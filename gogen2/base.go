package gogen2

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"strings"
)

type GogenBuilder struct {
	path          string
	goModPath     string
	importMap     map[Expression]GogenImport // menampung semua import tapi tidak semuanya dipakai oleh struct yg kita cari
	usedImport    map[Expression]GogenImport // menampung import yg dipakai saja
	mapOfRequire  map[RequirePath]CompletePath
	expressionMap map[Expression][]string   // grouping FieldType
	typeMap       map[FieldType]ast.Expr    // menampung semua type yang mungkin akan dipakai oleh struct target
	unknownTypes  map[FieldName]*GogenField // menampung hanya type yang diperlukan oleh field struct target
}

func (gsb *GogenBuilder) handleGoMod() error {

	const goModFileName = "go.mod"

	dataInBytes, err := os.ReadFile(goModFileName)
	if err != nil {
		return err
	}

	parsedGoMod, err := modfile.Parse(goModFileName, dataInBytes, nil)
	if err != nil {
		return err
	}

	for _, r := range parsedGoMod.Require {

		if len(r.Syntax.Token) == 1 {
			gsb.mapOfRequire[RequirePath(r.Syntax.Token[0])] = CompletePath(fmt.Sprintf("%v/pkg/mod/%v", build.Default.GOPATH, r.Syntax.Token[0]))
			continue
		}

		gsb.mapOfRequire[RequirePath(r.Syntax.Token[0])] = CompletePath(fmt.Sprintf("%v/pkg/mod/%v@%v", build.Default.GOPATH, r.Syntax.Token[0], r.Syntax.Token[1]))
	}

	return nil
}

func (gsb *GogenBuilder) handleImport(genDecl *ast.GenDecl) {

	for _, spec := range genDecl.Specs {

		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		// kita ambil import dari file yg sedang dibaca
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		lenImportPath := len(importPath)

		// kita cek apakah dia ada ada di gomod, 3rdlib, atau internal apps
		cp, exist := gsb.mapOfRequire[RequirePath(importPath)]
		if exist {

			pathToLib := fmt.Sprintf("%v%v", cp, importPath[lenImportPath:])

			pkgs, err := parser.ParseDir(token.NewFileSet(), pathToLib, nil, parser.PackageClauseOnly)
			if err != nil {
				panic(err)
			}

			for _, pkg := range pkgs {

				name := ""
				expr := pkg.Name
				if importSpec.Name != nil {
					name = importSpec.Name.String()
					expr = name
				}

				gi := GogenImport{
					Name:         name,
					Path:         ImportPath(importPath),
					Expression:   Expression(expr),
					ImportType:   ImportTypeExtModule,
					CompletePath: CompletePath(pathToLib),
				}

				gsb.importMap[gi.Expression] = gi
			}

			continue
		}

		name := ""
		expr := importPath[strings.LastIndex(importPath, "/")+1:]
		if importSpec.Name != nil {
			name = importSpec.Name.String()
			expr = name
		}
		importType := ImportTypeGoSDK
		completePath := CompletePath(fmt.Sprintf("%s/src/%s", build.Default.GOROOT, expr))

		if strings.HasPrefix(importPath, gsb.goModPath) {
			importType = ImportTypeProject
			completePath = CompletePath(importPath[len(gsb.goModPath)+1:])
		}

		gi := GogenImport{
			Name:         name,
			Path:         ImportPath(importPath),
			CompletePath: completePath,
			Expression:   Expression(expr),
			ImportType:   importType,
		}

		gsb.importMap[gi.Expression] = gi

	}

}

func (gsb *GogenBuilder) handleUsedImport(expr ast.Expr) []Expression {

	switch fieldType := expr.(type) {
	case *ast.StructType:

		str := make([]Expression, 0)
		for _, f := range fieldType.Fields.List {
			str = append(str, gsb.handleUsedImport(f.Type)...)
		}
		return str

	case *ast.SelectorExpr:
		x := Expression(fieldType.X.(*ast.Ident).String())
		sel := fieldType.Sel.String()

		gsb.expressionMap[x] = append(gsb.expressionMap[x], sel)

		return []Expression{Expression(fieldType.X.(*ast.Ident).String())}

	case *ast.StarExpr:
		return gsb.handleUsedImport(fieldType.X)

	case *ast.MapType:
		str := make([]Expression, 0)
		key := gsb.handleUsedImport(fieldType.Key)
		if key != nil {
			str = append(str, key...)
		}
		value := gsb.handleUsedImport(fieldType.Value)
		if value != nil {
			str = append(str, value...)
		}
		return str

	case *ast.ArrayType:
		return gsb.handleUsedImport(fieldType.Elt)

	case *ast.ChanType:
		return gsb.handleUsedImport(fieldType.Value)

	case *ast.FuncType:
		str := make([]Expression, 0)

		if fieldType.Params.NumFields() > 0 {
			for _, x := range fieldType.Params.List {
				str = append(str, gsb.handleUsedImport(x.Type)...)
			}
		}

		if fieldType.Results.NumFields() > 0 {
			for _, x := range fieldType.Results.List {
				str = append(str, gsb.handleUsedImport(x.Type)...)
			}
		}

		return str

	}

	return []Expression{}
}

func (gsb *GogenBuilder) checkDefaultValue(gf *GogenField) {

	if string(gf.DataType.Type) == gf.DataType.DefaultValue {

		logDebug("karena defaultValue utk field %v dengan type %v belum final, kita cek ke map", gf.Name, gf.DataType.Type)

		typeSpecFromMap, exist := gsb.typeMap[gf.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam map. pencarian default value utk var %v ditunda dan sudah didaftarkan dalam unknownTypes", gf.DataType.DefaultValue, gf.Name)
			gsb.unknownTypes[gf.Name] = gf

			logDebug("status unknownTypes : %+v", gsb.unknownTypes)
			return
		}

		oldDefaultValue := gf.DataType.DefaultValue
		gf.handleDefaultValue(typeSpecFromMap)

		logDebug("dataType %v ada di map. defaultValue %v sudah di replace dengan %v", gf.DataType.DefaultValue, oldDefaultValue, gf.DataType.DefaultValue)

		return
	}

	logDebug("default value utk field %v dengan dataType %v sudah final, yaitu : %v", gf.Name, gf.DataType.DefaultValue, gf.DataType.DefaultValue)

}
