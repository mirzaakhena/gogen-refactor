package gogen

import (
	"fmt"
	"strings"
)

func getRepoTypeName(repoMethodName string) string {
	return fmt.Sprintf("%sRepo", repoMethodName)
}

func getServiceTypeName(serviceMethodName string) string {
	return fmt.Sprintf("%sService", serviceMethodName)
}

func getUsecaseFolder(domainName, usecaseName string) string {
	return fmt.Sprintf("domain_%s/usecase/%s", strings.ToLower(domainName), strings.ToLower(usecaseName))
}

func getRepositoryFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/model/repository", strings.ToLower(domainName))
}

func getRepositoryFilename(domainName string) string {
	return fmt.Sprintf("%s/repository.go", getRepositoryFolder(domainName))
}

func getInteractorFilename(domainName, usecaseName string) string {
	return fmt.Sprintf("%s/interactor.go", getUsecaseFolder(domainName, usecaseName))
}

func getOutportName() string {
	return "Outport"
}

func getInjectedCodeLocation() string {
	return "//!"
}
