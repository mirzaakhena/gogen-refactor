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

				if typeSpecName != targetTypeName {
					LogDebug(3, "Not Expected Target Type. it is %s as %T in %v", typeSpecName, typeSpec.Type, fset.File(astFile.Package).Name())
					return false
				}

				LogDebug(3, "Found Target Type %s in %v!!!", targetTypeName, fset.File(astFile.Package).Name())

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

//func traceGeneralType(packagePath string, targetTypeName string) (*GogenFieldType, error) {
//
//	var gft *GogenFieldType
//
//	fset := token.NewFileSet()
//	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, pkg := range pkgs {
//
//		for _, astFile := range pkg.Files {
//
//			if err != nil {
//				LogDebug(2, "ignore everything since we are done or has an err")
//				return nil, err
//			}
//
//			ast.Inspect(astFile, func(node ast.Node) bool {
//
//				if err != nil {
//					LogDebug(3, "ignore everything since we have err : %v", err.Error())
//					return false
//				}
//
//				typeSpec, ok := node.(*ast.TypeSpec)
//				if !ok {
//					return true
//				}
//
//				typeSpecName := typeSpec.Name.String()
//
//				if typeSpecName != targetTypeName {
//					LogDebug(3, "Not Expected Target Type. it is %s as %T in %v", typeSpecName, typeSpec.Type, fset.File(astFile.Package).Name())
//					return false
//				}
//
//				LogDebug(3, "Found Target Type %s in %v!!!", targetTypeName, fset.File(astFile.Package).Name())
//
//				gft = &GogenFieldType{
//					Name:         FieldType(targetTypeName),
//					Expr:         typeSpec.Type,
//					DefaultValue: "",
//					File:         astFile,
//				}
//
//				return true
//			})
//
//		}
//
//	}
//
//	if gft == nil {
//		return nil, fmt.Errorf("target type %v not found", targetTypeName)
//	}
//
//	return gft, nil
//}

func (g *GogenInterfaceBuilder) traceInterfaceType(packagePath string, interfaceTargetName string) (*GogenInterface, error) {

	gogenInterfaceTarget := NewGogenInterface()

	unknownInterfaces := map[FieldType]*GogenInterface{}

	unknownFields := map[FieldMethodSignature]*GogenField{}

	collectedType := map[FieldType]*TypeProperties{}

	err := traceGeneralType(packagePath, interfaceTargetName, collectedType, func(tp TypeProperties) error {
		return g.handleGogenInterface(gogenInterfaceTarget, unknownInterfaces, unknownFields, &tp)
	})
	if err != nil {
		return nil, err
	}

	//LogDebug(0, "from path %v try to find interface with name %v", packagePath, interfaceTargetName)
	//fset := token.NewFileSet()
	//pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, pkg := range pkgs {
	//
	//	packageName := PackageName(pkg.Name)
	//
	//	LogDebug(1, "package %s", packageName)
	//
	//	for _, astFile := range pkg.Files {
	//
	//		if err != nil {
	//			LogDebug(2, "ignore everything since we are done or has an err")
	//			return nil, err
	//		}
	//
	//		ast.Inspect(astFile, func(node ast.Node) bool {
	//
	//			if err != nil {
	//				LogDebug(3, "ignore everything since we have err : %v", err.Error())
	//				return false
	//			}
	//
	//			typeSpec, ok := node.(*ast.TypeSpec)
	//			if !ok {
	//				return true
	//			}
	//
	//			typeSpecName := typeSpec.Name.String()
	//
	//			tp := TypeProperties{
	//				AstFile:  astFile,
	//				TypeSpec: typeSpec,
	//			}
	//
	//			collectedType[FieldType(typeSpecName)] = &tp
	//
	//			if typeSpecName != interfaceTargetName {
	//
	//				LogDebug(3, "Not Expected Interface Target Type. it is %s as %T in %v", typeSpecName, typeSpec.Type, fset.File(astFile.Package).Name())
	//
	//				return false
	//			}
	//
	//			LogDebug(3, "Found Interface Target Type %s in %v!!!", interfaceTargetName, fset.File(astFile.Package).Name())
	//
	//			err = g.handleGogenInterface(gogenInterfaceTarget, unknownInterfaces, unknownFields, &tp)
	//			if err != nil {
	//				return false
	//			}
	//
	//			return true
	//		})
	//
	//	}
	//
	//}

	// solve all unknown interfaces
	for fieldType, theInterface := range unknownInterfaces {

		typeProperties, exist := collectedType[fieldType]
		if !exist {
			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		err = g.handleGogenInterface(theInterface, unknownInterfaces, unknownFields, typeProperties)
		if err != nil {
			return nil, err
		}

		delete(unknownInterfaces, fieldType)

	}

	// solve all unknown fields
	for _, uf := range unknownFields {
		uf.DataType.DefaultValue, err = g.handleDefaultValue(uf, uf.DataType.Expr, uf.DataType.File)
		if err != nil {
			return nil, err
		}
	}

	return gogenInterfaceTarget, nil
}

func (g *GogenInterfaceBuilder) handleDefaultValue(gf *GogenField, expr ast.Expr, astFile *ast.File) (string, error) {

	switch exprType := expr.(type) {
	case *ast.Ident:

		if exprType.Obj != nil {
			return g.handleDefaultValue(gf, exprType.Obj.Decl.(*ast.TypeSpec).Type, astFile)
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
			if string(gf.DataType.Name) != exprType.String() {
				return fmt.Sprintf("%s(%s)", gf.DataType.Name, basicDefaultValue), nil
			}
			return basicDefaultValue, nil
		}

		return exprType.String(), nil

	case *ast.StructType:
		return fmt.Sprintf("%v{}", gf.DataType.Name), nil

	case *ast.ArrayType:
		return fmt.Sprintf("%s{}", gf.DataType.Name), nil

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

		// TODO hardcoded for context temporary
		if importInFile[theX].Path == "context" {
			return string(gf.Name), nil
		}

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

		result, err := g.handleDefaultValue(gf, gft.Expr, astFile)
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

func (g *GogenInterfaceBuilder) handleGogenInterface(gi *GogenInterface, unknownInterfaces map[FieldType]*GogenInterface, unknownFields map[FieldMethodSignature]*GogenField, typeProperties *TypeProperties) error {

	gi.InterfaceType = &GogenFieldType{
		Name:         FieldType(typeProperties.TypeSpec.Name.String()),
		Expr:         typeProperties.TypeSpec.Type,
		DefaultValue: "nil",
		File:         typeProperties.AstFile,
	}

	switch ts := typeProperties.TypeSpec.Type.(type) {
	case *ast.InterfaceType:

		LogDebug(3, "as interface with name %s", gi.InterfaceType.Name)

		for _, method := range ts.Methods.List {
			switch methodType := method.Type.(type) {

			case *ast.FuncType:
				err := g.handleDirectInlineMethod(gi, unknownFields, method)
				if err != nil {
					return err
				}

			case *ast.Ident:
				err := g.handleIdent(gi, unknownInterfaces, unknownFields, typeProperties, methodType)
				if err != nil {
					return err
				}

			case *ast.SelectorExpr:
				err := g.handleSelector(gi, typeProperties, methodType)
				if err != nil {
					return err
				}

			default:
				// TODO what about type alias?
				err := fmt.Errorf("un-handled method type %T", methodType)
				return err
			}
		}

	case *ast.Ident:

		err := g.handleIdent(gi, unknownInterfaces, unknownFields, typeProperties, ts)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("this is not an interface but %T", ts)

	}

	return nil
}

func (g *GogenInterfaceBuilder) handleSelector(gi *GogenInterface, typeProperties *TypeProperties, methodType *ast.SelectorExpr) error {

	//interfacePath, err := getPathFromSelector(typeProperties.AstFile, gomodProperties, methodType)
	//if err != nil {
	//	return err
	//}

	importInFile, err := g.handleImport(typeProperties.AstFile.Imports)
	if err != nil {
		return err
	}

	theX := Expression(methodType.X.(*ast.Ident).String())

	interfacePath := string(importInFile[theX].CompletePath)

	internalGi, err := g.traceInterfaceType(interfacePath, methodType.Sel.String())
	if err != nil {
		return err
	}

	gi.Interfaces = append(gi.Interfaces, internalGi)

	return nil
}

func (g *GogenInterfaceBuilder) handleIdent(gi *GogenInterface, unknownInterfaces map[FieldType]*GogenInterface, unknownFields map[FieldMethodSignature]*GogenField, typeProperties *TypeProperties, methodType *ast.Ident) error {
	internalGi := new(GogenInterface)
	gi.Interfaces = append(gi.Interfaces, internalGi)

	if methodType.Obj == nil {

		name := FieldType(methodType.String())

		internalGi.InterfaceType = &GogenFieldType{
			Name:         name,
			Expr:         methodType,
			DefaultValue: "nil",
			File:         typeProperties.AstFile,
		}

		unknownInterfaces[name] = internalGi

		LogDebug(3, "unknown %s", methodType.String())

		return nil
	}

	newTypeSpec, ok := methodType.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not type", methodType.String())
	}

	newTp := TypeProperties{
		AstFile:  typeProperties.AstFile,
		TypeSpec: newTypeSpec,
	}

	err := g.handleGogenInterface(internalGi, unknownInterfaces, unknownFields, &newTp)
	if err != nil {
		return err
	}

	return nil
}

func (g *GogenInterfaceBuilder) handleImport(importSpecs []*ast.ImportSpec) (map[Expression]*GogenImport, error) {

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

			path := importPath[strings.LastIndex(string(importPath), "/")+1:]

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

func (g *GogenInterfaceBuilder) handleGoMod(goModFilePath string) error {

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
			LogDebug(1, "path %v, version %v", ip, v)
			continue
		}
	}

	//return absolutePathProject, moduleName, requirePath, nil

	return nil
}

func (g *GogenInterfaceBuilder) handleDirectInlineMethod(gi *GogenInterface, unknownFields map[FieldMethodSignature]*GogenField, method *ast.Field) error {

	if method.Names == nil && len(method.Names) > 0 {
		err := fmt.Errorf("method must have name")
		return err
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil
	}

	LogDebug(4, "as function found method %s", methodName)

	gm := newGogenMethod(methodName)

	gi.Methods = append(gi.Methods, gm)

	methodType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("somehow cannot convert to FuncType")
	}

	g.handleFuncParamResultType(methodType, unknownFields, gm, gi.InterfaceType.File)

	return nil
}

func (g *GogenInterfaceBuilder) handleFuncParamResultType(methodType *ast.FuncType, unknownFields map[FieldMethodSignature]*GogenField, gm *GogenMethod, astFile *ast.File) {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			fieldType := FieldType(getTypeAsString(param.Type))

			if param.Names != nil {

				for _, n := range param.Names {
					gf := NewGogenField(n.String(), param.Type, fieldType, astFile)
					gm.Params = append(gm.Params, gf)
					unknownFields[NewFieldMethodSignature(gm.Name, gf.Name)] = gf
				}

			} else {

				gf := NewGogenField(getSel(param.Type), param.Type, fieldType, astFile)
				gm.Params = append(gm.Params, gf)
				unknownFields[NewFieldMethodSignature(gm.Name, gf.Name)] = gf

			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			fieldType := FieldType(getTypeAsString(result.Type))

			if result.Names != nil {

				for _, n := range result.Names {
					gf := NewGogenField(n.String(), result.Type, fieldType, astFile)
					gm.Results = append(gm.Results, gf)
					unknownFields[NewFieldMethodSignature(gm.Name, gf.Name)] = gf
				}

			} else {

				gf := NewGogenField(getSel(result.Type), result.Type, fieldType, astFile)
				gm.Results = append(gm.Results, gf)
				unknownFields[NewFieldMethodSignature(gm.Name, gf.Name)] = gf

			}

		}
	}
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

func PrintGogenInterface(level int, gft *GogenInterface) {
	if gft.InterfaceType != nil {

		LogDebug(level, "GogenType %s", gft.InterfaceType.Name)

		for _, v := range gft.Interfaces {
			PrintGogenInterface(level+1, v)
		}
		for _, v := range gft.Methods {
			LogDebug(level+1, "Method %s", v.Name)
		}
	}
}

func LogDebug(identationLevel int, format string, args ...any) {
	x := fmt.Sprintf("%%%ds", identationLevel*2)
	y := fmt.Sprintf(x, "")
	fmt.Printf(y+format+"\n", args...)
}

func PrintAllMethod(gi *GogenInterface) []*GogenMethod {
	return traceMethod(gi, make([]*GogenMethod, 0, 5))
}

func traceMethod(gi *GogenInterface, gms []*GogenMethod) []*GogenMethod {

	for _, m := range gi.Methods {
		gms = appendGogenMethod(gms, m)
	}

	for _, f := range gi.Interfaces {
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
