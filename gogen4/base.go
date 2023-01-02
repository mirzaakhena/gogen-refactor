package gogen4

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
	"strings"
)

type GogenBase struct {
	GoModProperties GoModProperties
}

func (g *GogenBase) handleDefaultValue(gf *GogenFieldType, expr ast.Expr, astFile *ast.File, collectedType map[FieldType]*TypeProperties) (string, error) {

	switch exprType := expr.(type) {
	case *ast.Ident:

		if exprType.Obj != nil {
			return g.handleDefaultValue(gf, exprType.Obj.Decl.(*ast.TypeSpec).Type, astFile, collectedType)
		}

		basicDefaultValue := ""

		for {

			if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
				basicDefaultValue = "0"
				break
			}

			if strings.HasPrefix(exprType.String(), "float") {
				basicDefaultValue = "0.0"
				break
			}

			if exprType.String() == "string" {
				basicDefaultValue = `""`
				break
			}

			if exprType.String() == "bool" {
				basicDefaultValue = `false`
				break
			}

			if exprType.String() == "any" {
				basicDefaultValue = `nil`
				break
			}

			if exprType.String() == "error" {
				basicDefaultValue = `nil`
				break
			}

			break
		}

		if basicDefaultValue != "" {
			if string(gf.Name) != exprType.String() {
				return fmt.Sprintf("%s(%s)", gf.Name, basicDefaultValue), nil
			}
			return basicDefaultValue, nil
		}

		typeProperties, exist := collectedType[gf.Name]
		if !exist {
			return "", fmt.Errorf("field type %v is not exist anywhere", gf.Name)
		}

		value, err := g.handleDefaultValue(gf, typeProperties.TypeSpec.Type, typeProperties.AstFile, collectedType)
		if err != nil {
			return "", err
		}

		return value, nil

	case *ast.StructType:
		return fmt.Sprintf("%v{}", gf.Name), nil

	case *ast.ArrayType:
		return fmt.Sprintf("%s{}", gf.Name), nil

	case *ast.SelectorExpr:

		importInFile, err := g.handleImport(astFile.Imports)
		if err != nil {
			return "", err
		}

		theX := Expression(exprType.X.(*ast.Ident).String())

		interfacePath := string(importInFile[theX].CompletePath)

		//interfacePath, err := getPathFromSelector(astFile, gomodProperties, exprType)
		//if err != nil {
		//	return "", err
		//}

		//// TODO hardcoded for context temporary
		//if importInFile[theX].Path == "context" {
		//	return string(gf.Name), nil
		//}

		var gft GogenFieldType

		err = traceGeneralType(interfacePath, exprType.Sel.String(), nil, func(tp TypeProperties) error {

			gft = GogenFieldType{
				Name:         FieldType(exprType.Sel.String()),
				Expr:         tp.TypeSpec.Type,
				DefaultValue: "",
				File:         astFile,
			}

			return nil

		})
		if err != nil {
			return "", err
		}

		result, err := g.handleDefaultValue(gf, gft.Expr, astFile, collectedType)
		if err != nil {
			return "", err
		}

		return result, nil

	case *ast.StarExpr:
		return "nil", nil

	case *ast.InterfaceType:
		return "nil", nil

	case *ast.MapType:
		return "nil", nil

	case *ast.ChanType:
		return "nil", nil

	case *ast.FuncType:
		return "nil", nil

	}

	return "", fmt.Errorf("unknown type")

}

func (g *GogenBase) handleImport(importSpecs []*ast.ImportSpec) (map[Expression]*GogenImport, error) {

	importInFile := map[Expression]*GogenImport{}

	for _, importSpec := range importSpecs {

		importPath := ImportPath(strings.Trim(importSpec.Path.Value, `"`))

		LogDebug(4, "read package with path %s", importPath)

		version, exist := g.GoModProperties.RequirePath[importPath]
		if exist {

			completePath := AbsolutePath(fmt.Sprintf("%v/pkg/mod/%v@%v", build.Default.GOPATH, importPath, version))

			pkgs, err := parser.ParseDir(token.NewFileSet(), string(completePath), nil, parser.PackageClauseOnly)
			if err != nil {
				return nil, err
			}

			for _, pkg := range pkgs {

				name := ""
				expr := Expression(pkg.Name)
				if importSpec.Name != nil {
					name = importSpec.Name.String()
					expr = Expression(name)
				}

				gi := GogenImport{
					Name:         name,
					Path:         importPath,
					Expression:   expr,
					ImportType:   ImportTypeExternalModule,
					CompletePath: completePath,
				}

				importInFile[expr] = &gi

				LogDebug(5, "found %v in path %v as %s", expr, completePath, gi.ImportType)

			}

			continue

		}

		if strings.HasPrefix(string(importPath), g.GoModProperties.ModuleName) {

			//path := importPath[strings.LastIndex(string(importPath), "/")+1:]
			path := importPath[len(g.GoModProperties.ModuleName):]

			completePathStr := filepath.Join(g.GoModProperties.AbsolutePathProject, string(path))

			pkgs, err := parser.ParseDir(token.NewFileSet(), completePathStr, nil, parser.PackageClauseOnly)
			if err != nil {
				return nil, err
			}

			for _, pkg := range pkgs {

				name := ""
				expr := Expression(pkg.Name)
				if importSpec.Name != nil {
					name = importSpec.Name.String()
					expr = Expression(name)
				}

				gi := GogenImport{
					Name:         name,
					Path:         importPath,
					Expression:   expr,
					ImportType:   ImportTypeInternalProject,
					CompletePath: AbsolutePath(completePathStr),
				}

				importInFile[expr] = &gi

				LogDebug(5, "found %v in path %v as %s", expr, completePathStr, gi.ImportType)

			}

			continue
		}

		// others
		{
			name := ""
			expr := Expression(importPath[strings.LastIndex(string(importPath), "/")+1:])
			if importSpec.Name != nil {
				name = importSpec.Name.String()
				expr = Expression(name)
			}

			completePath := AbsolutePath(fmt.Sprintf("%s/src/%s", build.Default.GOROOT, expr))

			gi := GogenImport{
				Name:         name,
				Path:         importPath,
				Expression:   expr,
				ImportType:   ImportTypeGoSDK,
				CompletePath: completePath,
			}

			importInFile[expr] = &gi

			LogDebug(5, "found %v in path %v as %s", expr, completePath, gi.ImportType)
		}

	}

	return importInFile, nil
}

func (g *GogenBase) handleGoMod(goModFilePath string) error {

	//gm := GoModProperties{
	//	RequirePath: map[ImportPath]Version{},
	//}

	cleanPath := filepath.Clean(goModFilePath)

	absoluteGomodPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return err
	}

	g.GoModProperties.AbsolutePathProject = absoluteGomodPath[:strings.LastIndex(absoluteGomodPath, filepath.Base(goModFilePath))]

	LogDebug(0, "absolute path project path : %s", g.GoModProperties.AbsolutePathProject)

	dataInBytes, err := os.ReadFile(goModFilePath)
	if err != nil {
		return err
	}

	parsedGoMod, err := modfile.Parse(goModFilePath, dataInBytes, nil)
	if err != nil {
		return err
	}

	//var moduleName string
	if len(parsedGoMod.Module.Syntax.Token) >= 2 {
		g.GoModProperties.ModuleName = parsedGoMod.Module.Syntax.Token[1]
	}
	LogDebug(0, "we get the module name : %v", g.GoModProperties.ModuleName)

	LogDebug(0, "read the require path in go.mod :")

	for _, req := range parsedGoMod.Require {
		if len(req.Syntax.Token) == 1 {
			g.GoModProperties.RequirePath[ImportPath(req.Syntax.Token[0])] = ""
			continue
		}
		if len(req.Syntax.Token) >= 2 {
			ip := ImportPath(req.Syntax.Token[0])
			v := Version(req.Syntax.Token[1])
			g.GoModProperties.RequirePath[ip] = v
			//LogDebug(1, "path %v, version %v", ip, v)
			continue
		}
	}

	//return absolutePathProject, moduleName, requirePath, nil

	return nil
}

func traceGeneralType(packagePath string, targetTypeName string, collectedType map[FieldType]*TypeProperties, afterFound func(tp TypeProperties) error) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		for _, astFile := range pkg.Files {

			if err != nil {
				LogDebug(2, "ignore everything since we are done or has an err")
				return err
			}

			ast.Inspect(astFile, func(node ast.Node) bool {

				if err != nil {
					LogDebug(3, "ignore everything since we have err : %v", err.Error())
					return false
				}

				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				typeSpecName := typeSpec.Name.String()

				tp := TypeProperties{
					AstFile:  astFile,
					TypeSpec: typeSpec,
				}

				if collectedType != nil {
					collectedType[FieldType(typeSpecName)] = &tp
				}

				LogDebug(3, "%v == %v", typeSpecName, targetTypeName)

				if typeSpecName != targetTypeName {
					LogDebug(3, "👎🏿 Not Expected Target Type. it is %s as %T in %v", typeSpecName, typeSpec.Type, fset.File(astFile.Package).Name())
					return false
				}

				LogDebug(3, "🎉 Found Target Type %s in %v!!!", targetTypeName, fset.File(astFile.Package).Name())

				err = afterFound(tp)
				if err != nil {
					return false
				}

				return true
			})

		}

	}

	return nil
}

func getSel(expr ast.Expr) string {

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:
		return getSel(fieldType.Sel)

	case *ast.StarExpr:
		return getSel(fieldType.X)

	case *ast.Ident:
		return fieldType.String()
	}

	return ""
}

func LogDebug(identationLevel int, format string, args ...any) {
	x := fmt.Sprintf("%%%ds", identationLevel*2)
	y := fmt.Sprintf(x, "")
	fmt.Printf(y+format+"\n", args...)
}

func PrintAllMethod(gi *GogenAnyType) []*GogenMethod {
	return traceMethod(gi, make([]*GogenMethod, 0, 5))
}

func traceMethod(gi *GogenAnyType, gms []*GogenMethod) []*GogenMethod {

	for _, m := range gi.Methods {
		gms = appendGogenMethod(gms, m)
	}

	for _, f := range gi.CompositionTypes {
		gms = traceMethod(f, gms)
	}

	return gms
}

func appendGogenMethod(gms []*GogenMethod, gm *GogenMethod) []*GogenMethod {

	lenGMS := len(gms)

	if lenGMS == cap(gms) {
		newSlice := make([]*GogenMethod, lenGMS, 2*lenGMS+1)
		copy(newSlice, gms)
		return append(newSlice, gm)

	}

	return append(gms, gm)

}

func PrintGogenAnyType(level int, gft *GogenAnyType) {
	if gft.GogenFieldType == nil {
		return
	}

	LogDebug(level, "GogenType %s %v", gft.GogenFieldType.Name, gft.GogenFieldType.DefaultValue)

	for _, v := range gft.CompositionTypes {
		PrintGogenAnyType(level+1, v)
	}
	for _, p := range gft.Fields {
		LogDebug(level+1, "Field %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
	}
	for _, v := range gft.Methods {
		LogDebug(level+1, "Method %s", v.Name)

		for _, p := range v.Params {
			LogDebug(level+2, "Param %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
		}

		for _, p := range v.Results {
			LogDebug(level+2, "Result %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
		}
	}

}

//func getPathFromSelector(astFile *ast.File, gomodProperties *GoModProperties, methodType *ast.SelectorExpr) (string, error) {
//	importInFile, err := handleImport(astFile.Imports, gomodProperties)
//	if err != nil {
//		return "", err
//	}
//
//	theX := Expression(methodType.X.(*ast.Ident).String())
//
//	return string(importInFile[theX].CompletePath), nil
//}
