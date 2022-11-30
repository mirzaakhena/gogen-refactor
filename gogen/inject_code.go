package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"
)

func InjectCodeAtTheEndOfFile(filename, templateCode string) ([]byte, error) {

	// reopen the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// write the template in the end of file
	buffer.WriteString(templateCode)
	buffer.WriteString("\n")

	return buffer.Bytes(), nil

}

func InjectToMiddleCode(existingFile, injectedCode string) ([]byte, error) {

	file, err := os.Open(existingFile)
	if err != nil {
		return nil, err
	}

	needToInject := false

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if strings.TrimSpace(row) == getInjectedCodeLocation() {

			needToInject = true

			// inject code
			buffer.WriteString(injectedCode)
			buffer.WriteString("\n")

			continue
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// if no injected marker found, then abort the next step
	if !needToInject {
		return nil, nil
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	// rewrite the file
	if err := os.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func InjectToSpecificCodeLine(filename string, injectedLine int, templateWithData string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == injectedLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}

	return buffer.Bytes(), nil

}

func injectToOutport(domainName, usecaseName, repoMethodName string) error {

	fset := token.NewFileSet()

	var iFace *ast.InterfaceType
	var astFile *ast.File

	repoTypeName := getRepoTypeName(repoMethodName)

	rootFolderName := getUsecaseFolder(domainName, usecaseName)

	pkgs, err := parser.ParseDir(fset, rootFolderName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	isAlreadyInjectedBefore := false

	// in every package
	for _, pkg := range pkgs {

		for _, file := range pkg.Files {

			ast.Inspect(file, func(node ast.Node) bool {

				ts, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if ts.Name.String() == getOutportName() {

					interfaceType, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						return true
					}

					for _, meth := range interfaceType.Methods.List {

						selType, ok := meth.Type.(*ast.SelectorExpr)
						if !ok {
							return true
						}

						if selType.Sel.String() == getRepoTypeName(repoMethodName) {
							isAlreadyInjectedBefore = true
							astFile = file
							iFace = interfaceType
							return false
						}

					}

				}

				return true
			})

		}

	}

	// if it is already injected, then nothing to do
	if isAlreadyInjectedBefore {
		return nil
	}

	if iFace == nil {
		return fmt.Errorf("outport struct not found")
	}

	// we want to inject it now
	// add new repository to outport interface
	iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: GetPackageName(getRepositoryFolder(domainName)),
			},
			Sel: &ast.Ident{
				Name: repoTypeName,
			},
		},
	})

	fileReadPath := fset.Position(astFile.Pos()).Filename

	// rewrite the outport
	f, err := os.Create(fileReadPath)
	if err != nil {
		return err
	}

	if err := printer.Fprint(f, fset, astFile); err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	err = Reformat(fileReadPath, nil)
	if err != nil {
		return err
	}

	return nil

}

func InjectToErrorEnum(fset *token.FileSet, filepath string, errorName, separator string) {

	astFile, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	// in every declaration like type, func, const
	for _, decl := range astFile.Decls {

		genDecl := decl.(*ast.GenDecl)

		if genDecl.Tok != token.CONST {
			continue
		}

		var errorNumber = 1000
		for _, spec := range genDecl.Specs {

			valueSpec := spec.(*ast.ValueSpec)
			if len(valueSpec.Values) == 0 {
				break
			}

			if len(valueSpec.Names) > 0 && strings.ToLower(valueSpec.Names[0].String()) == strings.ToLower(errorName) {
				fmt.Printf("error code %s already exist\n", errorName)
				return
			}

			basicList := valueSpec.Values[0].(*ast.BasicLit)
			errorCodeWithMessage := strings.Split(basicList.Value, " ")
			if len(errorCodeWithMessage) == 0 {
				continue
			}

			errorCodeOnly := strings.Split(errorCodeWithMessage[0], separator)
			if len(errorCodeOnly) < 2 || errorCodeOnly[1] == "" {
				continue
			}

			n, err := strconv.Atoi(errorCodeOnly[1])
			if err != nil {
				continue
			}
			errorNumber = n
		}

		errorValue := fmt.Sprintf("\"%s%04d %s\"", separator, errorNumber+1, SpaceCase(errorName))

		genDecl.Specs = append(genDecl.Specs, &ast.ValueSpec{
			Names:  []*ast.Ident{{Name: PascalCase(errorName)}},
			Type:   &ast.SelectorExpr{X: &ast.Ident{Name: "apperror"}, Sel: &ast.Ident{Name: "ErrorType"}},
			Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: errorValue}},
		})

		//ast.Print(fset, decl)

	}

	{
		f, err := os.Create(filepath)
		if err != nil {
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				os.Exit(1)
			}
		}(f)
		err = printer.Fprint(f, fset, astFile)
		if err != nil {
			return
		}
	}

	//{
	//	err := printer.Fprint(os.Stdout, fset, astFile)
	//	if err != nil {
	//		return
	//	}
	//}
}

type Usecase struct {
	Name                 string
	InportRequestFields  []*StructField
	InportResponseFields []*StructField
}

type StructField struct {
	Name string
	Type string
}

func injectUsecaseInportFields(domainName string, usecaseName string, usecases []*Usecase) []*Usecase {

	inportRequestFields := make([]*StructField, 0)
	inportResponseFields := make([]*StructField, 0)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, getUsecaseFolder(domainName, usecaseName), nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	// in every package
	for _, pkg := range pkgs {

		// in every files
		for _, file := range pkg.Files {

			// in every declaration like type, func, const
			for _, decl := range file.Decls {

				// focus only to type
				gen, ok := decl.(*ast.GenDecl)

				if ok && gen.Tok == token.TYPE {

					for _, specs := range gen.Specs {

						ts, ok := specs.(*ast.TypeSpec)
						if !ok {
							continue
						}

						structObj, isStruct := ts.Type.(*ast.StructType)

						if isStruct {

							structName := strings.ToLower(ts.Name.String())

							if structName == "inportrequest" {

								for _, f := range structObj.Fields.List {
									//fieldType := typeHandler{PrefixExpression: strings.ToLower(usecaseName)}.Start(f.Type)
									fieldType := typeHandler{}.Start(f.Type)
									for _, name := range f.Names {
										inportRequestFields = append(inportRequestFields, &StructField{
											Name: name.String(),
											Type: fieldType,
										})
									}
								}
							}

							if structName == "inportresponse" {

								for _, f := range structObj.Fields.List {
									//fieldType := typeHandler{PrefixExpression: strings.ToLower(usecaseName)}.Start(f.Type)
									fieldType := typeHandler{}.Start(f.Type)
									for _, name := range f.Names {
										inportResponseFields = append(inportResponseFields, &StructField{
											Name: name.String(),
											Type: fieldType,
										})
									}
								}
							}

							if structName == fmt.Sprintf("%sinteractor", file.Name) {
								usecaseNameWithInteractor := ts.Name.String()
								usecaseNameOnly := usecaseNameWithInteractor[:strings.LastIndex(usecaseNameWithInteractor, "Interactor")]
								usecases = append(usecases, &Usecase{
									Name:                 usecaseNameOnly,
									InportRequestFields:  inportRequestFields,
									InportResponseFields: inportResponseFields,
								})

							}
						}

					}

				}
			}
		}
	}

	return usecases
}

func InjectRegisterUsecaseInApplication(usecaseName, appFile string, injectedLine int) ([]byte, error) {

	templateWithData := fmt.Sprintf("%s.NewUsecase(datasource),", usecaseName)

	// u.AddUsecase(runordercreate.NewUsecase(datasource))
	//templateWithData := fmt.Sprintf("u.AddUsecase(%s.NewUsecase(datasource))", usecaseName)

	//beforeLine, err := getInjectedLineInApplication(appFile)
	//if err != nil {
	//	return nil, err
	//}

	file, err := os.Open(appFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == injectedLine+1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}
	return buffer.Bytes(), nil
}
