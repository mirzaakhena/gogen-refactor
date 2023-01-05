package gogen5

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

func GetDefaultValue(gf *GogenFieldType, expr ast.Expr, collectedType map[GogenFieldTypeName]TypeProperties, handleSelectorDefaultValue func(selectorExpr *ast.SelectorExpr) (ast.Expr, error)) (string, error) {

	switch exprType := expr.(type) {
	case *ast.Ident:

		if exprType.Obj != nil {
			return GetDefaultValue(gf, exprType.Obj.Decl.(*ast.TypeSpec).Type, collectedType, handleSelectorDefaultValue)
		}

		basicDefaultValue := handlePrimitiveType(exprType)

		if basicDefaultValue != "" {
			if string(gf.Name) != exprType.String() {
				return fmt.Sprintf("%s(%s)", gf.Name, basicDefaultValue), nil
			}
			return basicDefaultValue, nil
		}

		LogDebug(1, ">>>1 %v %v", gf.Name, "dv")

		tp, exist := collectedType[gf.Name]
		if !exist {
			return "", fmt.Errorf("field type %v is not exist anywhere", gf.Name)
		}

		value, err := GetDefaultValue(gf, tp.TypeSpec.Type, collectedType, nil)
		if err != nil {
			return "", err
		}

		return value, nil

	case *ast.StructType:
		return fmt.Sprintf("%v{}", gf.Name), nil

	case *ast.ArrayType:
		return fmt.Sprintf("%s{}", gf.Name), nil

	case *ast.SelectorExpr:

		if handleSelectorDefaultValue == nil {
			return "", fmt.Errorf("handleSelectorDefaultValue must not nil")
		}

		newExpr, err := handleSelectorDefaultValue(exprType)
		if err != nil {
			return "", err
		}

		result, err := GetDefaultValue(gf, newExpr, collectedType, nil)
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

func handlePrimitiveType(exprType *ast.Ident) string {
	basicDefaultValue := ""

	if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
		basicDefaultValue = "0"
	} else if strings.HasPrefix(exprType.String(), "float") {
		basicDefaultValue = "0.0"
	} else if exprType.String() == "string" {
		basicDefaultValue = `""`
	} else if exprType.String() == "bool" {
		basicDefaultValue = `false`
	} else if exprType.String() == "any" {
		basicDefaultValue = `nil`
	} else if exprType.String() == "error" {
		basicDefaultValue = `nil`
	}
	return basicDefaultValue
}

func TraceNode(packagePath string, afterFound func(astPackage *ast.Package, astFile *ast.File, node ast.Node) (bool, error)) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var continueFindOtherNode bool

	for _, pkg := range pkgs {

		for _, astFile := range pkg.Files {

			ast.Inspect(astFile, func(node ast.Node) bool {

				continueFindOtherNode, err = afterFound(pkg, astFile, node)
				if err != nil || !continueFindOtherNode {
					return false
				}

				return true
			})

			if err != nil {
				return err
			}

		}

	}

	return nil
}

func GetImportExpression(importSpecs []*ast.ImportSpec, gmp GoModProperties) (map[Expression]*GogenImport, error) {

	importInFile := map[Expression]*GogenImport{}

	for _, importSpec := range importSpecs {

		importPath := ImportPath(strings.Trim(importSpec.Path.Value, `"`))

		version, exist := gmp.RequirePath[importPath]
		if exist {

			completePath := fmt.Sprintf("%v/pkg/mod/%v@%v", build.Default.GOPATH, importPath, version)

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

			}

			continue

		}

		if strings.HasPrefix(string(importPath), gmp.ModuleName) {

			//path := importPath[strings.LastIndex(string(importPath), "/")+1:]
			path := importPath[len(gmp.ModuleName):]

			completePathStr := filepath.Join(gmp.AbsolutePathProject, string(path))

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
					CompletePath: completePathStr,
				}

				importInFile[expr] = &gi

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

			completePath := fmt.Sprintf("%s/src/%s", build.Default.GOROOT, expr)

			gi := GogenImport{
				Name:         name,
				Path:         importPath,
				Expression:   expr,
				ImportType:   ImportTypeGoSDK,
				CompletePath: completePath,
			}

			importInFile[expr] = &gi

		}

	}

	return importInFile, nil
}

func GetGoModProperties(goModFilePath string) (GoModProperties, error) {

	gmp := NewGoModProperties()

	cleanPath := filepath.Clean(goModFilePath)

	absoluteGomodPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return gmp, err
	}

	gmp.AbsolutePathProject = absoluteGomodPath[:strings.LastIndex(absoluteGomodPath, filepath.Base(goModFilePath))]

	dataInBytes, err := os.ReadFile(goModFilePath)
	if err != nil {
		return gmp, err
	}

	parsedGoMod, err := modfile.Parse(goModFilePath, dataInBytes, nil)
	if err != nil {
		return gmp, err
	}

	//var moduleName string
	if len(parsedGoMod.Module.Syntax.Token) >= 2 {
		gmp.ModuleName = parsedGoMod.Module.Syntax.Token[1]
	}

	for _, req := range parsedGoMod.Require {
		if len(req.Syntax.Token) == 1 {
			gmp.RequirePath[ImportPath(req.Syntax.Token[0])] = ""
			continue
		}
		if len(req.Syntax.Token) >= 2 {
			ip := ImportPath(req.Syntax.Token[0])
			v := Version(req.Syntax.Token[1])
			gmp.RequirePath[ip] = v
			continue
		}
	}

	//return absolutePathProject, moduleName, requirePath, nil

	return gmp, nil
}

func GetBasicType(expr ast.Expr) string {

	switch fieldType := expr.(type) {

	case *ast.Ident:
		return fieldType.String()

	case *ast.SelectorExpr:
		return GetBasicType(fieldType.Sel)

	case *ast.StarExpr:
		return GetBasicType(fieldType.X)

	case *ast.ArrayType:
		return GetBasicType(fieldType.Elt)

	}

	return ""
}

func GetSelectorExpr(expr ast.Expr) *ast.SelectorExpr {

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:
		return fieldType

	case *ast.StarExpr:
		return GetSelectorExpr(fieldType.X)

	case *ast.ArrayType:
		return GetSelectorExpr(fieldType.Elt)

	}

	return nil
}

func PrintGogenAnyType(level int, gft *GogenAnyType) {

	LogDebug(level, "===<<<<=============================================================")

	printGogenAnyTypeStruct{}.printGogenAnyTypeLoop(level, gft)

	LogDebug(level, "===>>>>=============================================================")

}

type printGogenAnyTypeStruct struct{}

func (x printGogenAnyTypeStruct) printGogenAnyTypeLoop(level int, gft *GogenAnyType) {
	//if gft.GogenFieldType == nil {
	//	return
	//}

	//LogDebug(level, "GogenType %s %v", gft.GogenFieldType.Name, gft.GogenFieldType.DefaultValue)

	LogDebug(level, "GogenType %s", gft.Name)

	for _, v := range gft.CompositionTypes {
		x.printGogenAnyTypeLoop(level+1, v)
	}
	for _, p := range gft.Fields {
		LogDebug(level+1, "Field %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
	}
	for _, v := range gft.Methods {
		LogDebug(level+1, "Method %s", v.Name)

		for _, p := range v.Params {
			LogDebug(level+2, "Param  %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
		}

		for _, p := range v.Results {
			LogDebug(level+2, "Result %s %s %s", p.Name, p.DataType.Name, p.DataType.DefaultValue)
		}
	}

}

func GetGogenImport(selectorExpr *ast.SelectorExpr, imports []*ast.ImportSpec, goMod GoModProperties) (*GogenImport, error) {
	importInFile, err := GetImportExpression(imports, goMod)
	if err != nil {
		return nil, err
	}
	return importInFile[Expression(selectorExpr.X.(*ast.Ident).String())], nil
}

func LogDebug(identationLevel int, format string, args ...any) {
	x := fmt.Sprintf("%%%ds", identationLevel*2)
	y := fmt.Sprintf(x, "")
	fmt.Printf(y+format+"\n", args...)
}
