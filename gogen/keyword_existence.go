package gogen

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

func IsRepoExist(repoPath, repoMethodName string) (bool, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, repoPath, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}

	found := false
	errMessage := ""

	// in every package
	for _, pkg := range pkgs {

		ast.Inspect(pkg, func(node ast.Node) bool {

			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if ts.Name.String() == getRepoTypeName(repoMethodName) {

				// repo is not an interface
				if _, ok := ts.Type.(*ast.InterfaceType); !ok {
					errMessage = "repo found but not an interface"
					return false
				}

				found = true
				return false
			}

			return true

		})

	}

	if len(errMessage) > 0 {
		return false, errors.New(errMessage)
	}

	return found, nil
}

func IsServiceExist(servicePath, serviceMethodName string) (bool, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, servicePath, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}

	found := false
	errMessage := ""

	// in every package
	for _, pkg := range pkgs {

		ast.Inspect(pkg, func(node ast.Node) bool {

			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if ts.Name.String() == getServiceTypeName(serviceMethodName) {

				_, isInterface := ts.Type.(*ast.InterfaceType)
				_, isStruct := ts.Type.(*ast.InterfaceType)

				// repo is not an interface or a service
				if isStruct || isInterface {
					errMessage = "service found but not interface or struct"
					return false
				}

				found = true
				return false
			}

			return true

		})

	}

	if len(errMessage) > 0 {
		return false, errors.New(errMessage)
	}

	return found, nil
}
