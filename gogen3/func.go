package gogen3

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

func handleImport(importSpecs []*ast.ImportSpec, gomodProperties *GoModProperties) (map[Expression]*GogenImport, error) {

	LogDebug(1, ">>>>>>>>>>>>>>>>>>>>>>>>>>>> %v", importSpecs)

	importInFile := map[Expression]*GogenImport{}

	for _, importSpec := range importSpecs {

		importPath := ImportPath(strings.Trim(importSpec.Path.Value, `"`))

		LogDebug(4, "read package with path %s", importPath)

		version, exist := gomodProperties.RequirePath[importPath]
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

		if strings.HasPrefix(string(importPath), gomodProperties.ModuleName) {

			path := importPath[strings.LastIndex(string(importPath), "/")+1:]

			completePathStr := filepath.Join(gomodProperties.AbsolutePathProject, string(path))

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

	return importInFile, nil
}

func handleGoMod(goModFilePath string) (*GoModProperties, error) {

	gm := GoModProperties{
		RequirePath: map[ImportPath]Version{},
	}

	cleanPath := filepath.Clean(goModFilePath)

	absoluteGomodPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, err
	}

	gm.AbsolutePathProject = absoluteGomodPath[:strings.LastIndex(absoluteGomodPath, filepath.Base(goModFilePath))]

	LogDebug(0, "absolute path project path : %s", gm.AbsolutePathProject)

	dataInBytes, err := os.ReadFile(goModFilePath)
	if err != nil {
		return nil, err
	}

	parsedGoMod, err := modfile.Parse(goModFilePath, dataInBytes, nil)
	if err != nil {
		return nil, err
	}

	//var moduleName string
	if len(parsedGoMod.Module.Syntax.Token) >= 2 {
		gm.ModuleName = parsedGoMod.Module.Syntax.Token[1]
	}
	LogDebug(0, "we get the module name : %v", gm.ModuleName)

	LogDebug(0, "read the require path in go.mod :")
	requirePath := map[ImportPath]Version{}
	for _, req := range parsedGoMod.Require {
		if len(req.Syntax.Token) == 1 {
			requirePath[ImportPath(req.Syntax.Token[0])] = ""
			continue
		}
		if len(req.Syntax.Token) >= 2 {
			ip := ImportPath(req.Syntax.Token[0])
			v := Version(req.Syntax.Token[1])
			requirePath[ip] = v
			LogDebug(1, "path %v, version %v", ip, v)
			continue
		}
	}

	//return absolutePathProject, moduleName, requirePath, nil

	return &gm, nil
}

func LogDebug(identationLevel int, format string, args ...any) {
	x := fmt.Sprintf("%%%ds", identationLevel*2)
	y := fmt.Sprintf(x, "")
	fmt.Printf(y+format+"\n", args...)
}

func handleDirectInlineMethod(method *ast.Field) (*GogenMethod, error) {

	if method.Names == nil && len(method.Names) > 0 {
		err := fmt.Errorf("method must have name")
		return nil, err
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil, nil
	}

	LogDebug(4, "as function found method %s", methodName)

	gm := NewGogenMethod(methodName)

	// TODO handle param later

	return gm, nil
}

func NewGogenMethod(methodName string) *GogenMethod {
	return &GogenMethod{
		Name:    GogenMethodName(methodName),
		Params:  make([]*GogenField, 0),
		Results: make([]*GogenField, 0),
	}
}

func PrintGogenInteface(level int, gft *GogenInterface) {
	if gft.InterfaceType != nil {

		LogDebug(level, "GogenType %s", gft.InterfaceType.Name)

		for _, v := range gft.Interfaces {
			PrintGogenInteface(level+1, v)
		}
		for _, v := range gft.Methods {
			LogDebug(level+1, "Method %s", v.Name)
		}
	}
}

func handleGogenInterface(gi *GogenInterface, unknownInterface map[FieldType]*GogenInterface, typeSpec *ast.TypeSpec, astFile *ast.File) error {

	switch ts := typeSpec.Type.(type) {
	case *ast.InterfaceType:

		gi.InterfaceType = &GogenFieldType{
			Name:         FieldType(typeSpec.Name.String()),
			Expr:         ts,
			DefaultValue: "nil",
			File:         astFile,
		}

		//gi.CurrentPackage = packageName
		gi.Interfaces = make([]*GogenInterface, 0)
		gi.Methods = make([]*GogenMethod, 0)

		LogDebug(3, "as interface with name %s", gi.InterfaceType.Name)

		for _, method := range ts.Methods.List {
			switch methodType := method.Type.(type) {
			case *ast.FuncType:

				gm, err := handleDirectInlineMethod(method)
				if err != nil {
					return err
				}

				if gm == nil {
					continue
				}

				gi.Methods = append(gi.Methods, gm)

				fields := make([]*GogenField, 0)

				handleFuncParamResultType(methodType, gm, fields)

			case *ast.Ident:

				newGi := new(GogenInterface)
				gi.Interfaces = append(gi.Interfaces, newGi)

				if methodType.Obj == nil {

					name := FieldType(methodType.String())

					newGi.InterfaceType = &GogenFieldType{
						Name:         name,
						Expr:         methodType,
						DefaultValue: "nil",
						File:         astFile,
					}

					unknownInterface[name] = newGi

					LogDebug(3, "unknown %s", methodType.String())
					continue
				}

				newTypeSpec, ok := methodType.Obj.Decl.(*ast.TypeSpec)
				if !ok {
					return fmt.Errorf("%s is not type", methodType.String())
				}

				err := handleGogenInterface(newGi, unknownInterface, newTypeSpec, astFile)
				if err != nil {
					return err
				}

			case *ast.SelectorExpr:

				newGi := new(GogenInterface)
				gi.Interfaces = append(gi.Interfaces, newGi)

				x := methodType.X.(*ast.Ident).String()
				sel := methodType.Sel.String()

				name := FieldType(fmt.Sprintf("%v.%v", x, sel))

				newGi.InterfaceType = &GogenFieldType{
					Name:         name,
					Expr:         methodType,
					DefaultValue: "nil",
					File:         astFile,
				}

				unknownInterface[name] = newGi

			default:
				// TODO what about type alias?
				err := fmt.Errorf("un-handled method type %T", methodType)
				return err
			}
		}

	default:
		return fmt.Errorf("this is not an interface but %T", ts)

	}

	return nil
}

func handleFuncParamResultType(methodType *ast.FuncType, gm *GogenMethod, fields []*GogenField) {

	// TODO later we need to calculate all the default value based on fields

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			if param.Names != nil {

				for _, n := range param.Names {
					gf := NewGogenField(n.String(), param.Type)
					gm.Params = append(gm.Params, gf)
					fields = append(fields, gf)
				}

			} else {

				gf := NewGogenField(getSel(param.Type), param.Type)
				gm.Params = append(gm.Params, gf)
				fields = append(fields, gf)

			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			if result.Names != nil {

				for _, n := range result.Names {
					gf := NewGogenField(n.String(), result.Type)
					gm.Params = append(gm.Params, gf)
					fields = append(fields, gf)
				}

			} else {

				gf := NewGogenField(getSel(result.Type), result.Type)
				gm.Params = append(gm.Params, gf)
				fields = append(fields, gf)

			}

		}
	}
}

func traceType(packagePath string, interfaceTargetName string, collectedType map[FieldType]*TypeProperties, afterFound func(tp *TypeProperties) error) error {

	LogDebug(0, "from path %v try to find interface with name %v", packagePath, interfaceTargetName)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		packageName := PackageName(pkg.Name)

		LogDebug(1, "package %s", packageName)

		for _, file := range pkg.Files {

			if err != nil {
				LogDebug(2, "ignore everything since we are done or has an err")
				return err
			}

			ast.Inspect(file, func(node ast.Node) bool {

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
					File:     file,
					TypeSpec: typeSpec,
				}

				collectedType[FieldType(typeSpecName)] = &tp

				if typeSpecName != interfaceTargetName {
					LogDebug(3, "Not Expected Interface Target Type. it is %s as %T in %v", typeSpecName, typeSpec.Type, fset.File(file.Package).Name())

					return false
				}

				LogDebug(3, "Found Interface Target Type %s in %v!!!", interfaceTargetName, fset.File(file.Package).Name())

				err = afterFound(&tp)
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

func NewGogenField(name string, expr ast.Expr) *GogenField {

	return &GogenField{
		Name: GogenFieldName(name),
		DataType: &GogenFieldType{
			Name:         FieldType(getTypeAsString(expr)),
			Expr:         expr,
			DefaultValue: "",
			File:         nil,
		},
	}

}
