package gogen5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase101(t *testing.T) {

	actGi, err := Build("./data_testing/project/interface001/p1", "./data_testing/project/go.mod", "MyInterface1")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "MyInterface1", actGi.Name.String())
	assert.Equal(t, 4, len(actGi.Methods))

	assert.Equal(t, "Method11", actGi.Methods[0].Name.String())
	assert.Equal(t, 2, len(actGi.Methods[0].Params))
	assert.Equal(t, 2, len(actGi.Methods[0].Results))

	assert.Equal(t, "x", actGi.Methods[0].Params[0].Name.String())
	assert.Equal(t, "int", actGi.Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, `0`, actGi.Methods[0].Params[0].DataType.DefaultValue)

	assert.Equal(t, "y", actGi.Methods[0].Params[1].Name.String())
	assert.Equal(t, "string", actGi.Methods[0].Params[1].DataType.Name.String())
	assert.Equal(t, `""`, actGi.Methods[0].Params[1].DataType.DefaultValue)

	assert.Equal(t, "bool", actGi.Methods[0].Results[0].Name.String())
	assert.Equal(t, "bool", actGi.Methods[0].Results[0].DataType.Name.String())
	assert.Equal(t, `false`, actGi.Methods[0].Results[0].DataType.DefaultValue)

	assert.Equal(t, "error", actGi.Methods[0].Results[1].Name.String())
	assert.Equal(t, "error", actGi.Methods[0].Results[1].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[0].Results[1].DataType.DefaultValue)

	assert.Equal(t, "Method12", actGi.Methods[1].Name.String())
	assert.Equal(t, 2, len(actGi.Methods[1].Params))
	assert.Equal(t, 2, len(actGi.Methods[1].Results))

	assert.Equal(t, "int", actGi.Methods[1].Params[0].Name.String())
	assert.Equal(t, "int", actGi.Methods[1].Params[0].DataType.Name.String())
	assert.Equal(t, `0`, actGi.Methods[1].Params[0].DataType.DefaultValue)

	assert.Equal(t, "string", actGi.Methods[1].Params[1].Name.String())
	assert.Equal(t, "string", actGi.Methods[1].Params[1].DataType.Name.String())
	assert.Equal(t, `""`, actGi.Methods[1].Params[1].DataType.DefaultValue)

	assert.Equal(t, "x", actGi.Methods[1].Results[0].Name.String())
	assert.Equal(t, "bool", actGi.Methods[1].Results[0].DataType.Name.String())
	assert.Equal(t, `false`, actGi.Methods[1].Results[0].DataType.DefaultValue)

	assert.Equal(t, "y", actGi.Methods[1].Results[1].Name.String())
	assert.Equal(t, "error", actGi.Methods[1].Results[1].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[1].Results[1].DataType.DefaultValue)

	assert.Equal(t, "Method13", actGi.Methods[2].Name.String())
	assert.Equal(t, 2, len(actGi.Methods[2].Params))
	assert.Equal(t, 0, len(actGi.Methods[2].Results))

	assert.Equal(t, "ctx3", actGi.Methods[2].Params[0].Name.String())
	assert.Equal(t, "context.Context", actGi.Methods[2].Params[0].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[2].Params[0].DataType.DefaultValue)

	assert.Equal(t, "handler", actGi.Methods[2].Params[1].Name.String())
	assert.Equal(t, "gin.HandlerFunc", actGi.Methods[2].Params[1].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[2].Params[1].DataType.DefaultValue)

	assert.Equal(t, "Method14", actGi.Methods[3].Name.String())
	assert.Equal(t, 1, len(actGi.Methods[3].Params))
	assert.Equal(t, 0, len(actGi.Methods[3].Results))

	assert.Equal(t, "a", actGi.Methods[3].Params[0].Name.String())
	assert.Equal(t, "*p4.MyStruct1", actGi.Methods[3].Params[0].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[3].Params[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.Fields))

	//

	assert.Equal(t, 4, len(actGi.Imports))

	assert.Equal(t, ``, actGi.Imports["p2"].Name)
	assert.Equal(t, "p2", actGi.Imports["p2"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/interface001/p2", actGi.Imports["p2"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.Imports["p2"].ImportType.String())

	assert.Equal(t, ``, actGi.Imports["context"].Name)
	assert.Equal(t, "context", actGi.Imports["context"].Expression.String())
	assert.Equal(t, "context", actGi.Imports["context"].Path.String())
	assert.Equal(t, "GO_SDK", actGi.Imports["context"].ImportType.String())

	assert.Equal(t, ``, actGi.Imports["gin"].Name)
	assert.Equal(t, "gin", actGi.Imports["gin"].Expression.String())
	assert.Equal(t, "github.com/gin-gonic/gin", actGi.Imports["gin"].Path.String())
	assert.Equal(t, "EXTERNAL_MODULE", actGi.Imports["gin"].ImportType.String())

	assert.Equal(t, ``, actGi.Imports["p4"].Name)
	assert.Equal(t, "p4", actGi.Imports["p4"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/interface001/p4", actGi.Imports["p4"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.Imports["p4"].ImportType.String())

	assert.Equal(t, 5, len(actGi.CompositionTypes))
	assert.Equal(t, "MyInterface2", actGi.CompositionTypes[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[0].Methods))

	assert.Equal(t, "Method21", actGi.CompositionTypes[0].Methods[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[0].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[0].Methods[0].Results))

	assert.Equal(t, "b", actGi.CompositionTypes[0].Methods[0].Params[0].Name.String())
	assert.Equal(t, "MyAliasInteger", actGi.CompositionTypes[0].Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, `MyAliasInteger(0)`, actGi.CompositionTypes[0].Methods[0].Params[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[0].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[0].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[0].CompositionTypes))
	assert.Equal(t, "MyInterface3", actGi.CompositionTypes[1].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[1].Methods))

	assert.Equal(t, "Method31", actGi.CompositionTypes[1].Methods[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[1].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[1].Methods[0].Results))

	assert.Equal(t, "c", actGi.CompositionTypes[1].Methods[0].Params[0].Name.String())
	assert.Equal(t, "[]MyAliasInteger", actGi.CompositionTypes[1].Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, `[]MyAliasInteger{}`, actGi.CompositionTypes[1].Methods[0].Params[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[1].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[1].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[1].CompositionTypes))
	assert.Equal(t, "MyInterface4", actGi.CompositionTypes[2].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[2].Methods))

	assert.Equal(t, "Method41", actGi.CompositionTypes[2].Methods[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Methods[0].Results))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].CompositionTypes))
	assert.Equal(t, "MyInterface5", actGi.CompositionTypes[3].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[3].Methods))

	assert.Equal(t, "Method51", actGi.CompositionTypes[3].Methods[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[3].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[3].Methods[0].Results))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[3].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[3].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[3].CompositionTypes))
	assert.Equal(t, "p2.MyInterface6", actGi.CompositionTypes[4].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[4].Methods))

	assert.Equal(t, "Method61", actGi.CompositionTypes[4].Methods[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].Methods[0].Results))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].Fields))

	//

	assert.Equal(t, 1, len(actGi.CompositionTypes[4].Imports))

	assert.Equal(t, ``, actGi.CompositionTypes[4].Imports["p3differentname"].Name)
	assert.Equal(t, "p3differentname", actGi.CompositionTypes[4].Imports["p3differentname"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/interface001/p3", actGi.CompositionTypes[4].Imports["p3differentname"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.CompositionTypes[4].Imports["p3differentname"].ImportType.String())

	assert.Equal(t, 2, len(actGi.CompositionTypes[4].CompositionTypes))
	assert.Equal(t, "MyInterface7", actGi.CompositionTypes[4].CompositionTypes[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[0].Methods))

	assert.Equal(t, "Method71", actGi.CompositionTypes[4].CompositionTypes[0].Methods[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[0].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[0].Methods[0].Results))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[0].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[0].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[0].CompositionTypes))
	assert.Equal(t, "p3differentname.MyInterface8", actGi.CompositionTypes[4].CompositionTypes[1].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[1].Methods))

	assert.Equal(t, "Method81", actGi.CompositionTypes[4].CompositionTypes[1].Methods[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].Methods[0].Results))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].Imports))

	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes))
	assert.Equal(t, "MyInterface9", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods))

	assert.Equal(t, "Method91", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Params))
	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Results))

	assert.Equal(t, "MyStruct1", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Params[0].Name.String())
	assert.Equal(t, "[]p4.MyStruct1", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, `[]p4.MyStruct1{}`, actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Methods[0].Params[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Fields))

	//

	assert.Equal(t, 1, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Imports))

	assert.Equal(t, ``, actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Imports["p4"].Name)
	assert.Equal(t, "p4", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Imports["p4"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/interface001/p4", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Imports["p4"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].Imports["p4"].ImportType.String())

	assert.Equal(t, 0, len(actGi.CompositionTypes[4].CompositionTypes[1].CompositionTypes[0].CompositionTypes))

}

func TestCase102(t *testing.T) {

	actGi, err := Build("./data_testing/project/struct001/p1", "./data_testing/project/go.mod", "MyStruct1")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "MyStruct1", actGi.Name.String())
	assert.Equal(t, 1, len(actGi.Methods))

	assert.Equal(t, "MethodStruct1", actGi.Methods[0].Name.String())
	assert.Equal(t, 1, len(actGi.Methods[0].Params))
	assert.Equal(t, 1, len(actGi.Methods[0].Results))

	assert.Equal(t, "x", actGi.Methods[0].Params[0].Name.String())
	assert.Equal(t, "int", actGi.Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, `0`, actGi.Methods[0].Params[0].DataType.DefaultValue)

	assert.Equal(t, "error", actGi.Methods[0].Results[0].Name.String())
	assert.Equal(t, "error", actGi.Methods[0].Results[0].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Methods[0].Results[0].DataType.DefaultValue)

	//

	assert.Equal(t, 10, len(actGi.Fields))

	assert.Equal(t, "Field1", actGi.Fields[0].Name.String())
	assert.Equal(t, "int", actGi.Fields[0].DataType.Name.String())
	assert.Equal(t, `0`, actGi.Fields[0].DataType.DefaultValue)

	assert.Equal(t, "Field2", actGi.Fields[1].Name.String())
	assert.Equal(t, "sync.WaitGroup", actGi.Fields[1].DataType.Name.String())
	assert.Equal(t, `sync.WaitGroup{}`, actGi.Fields[1].DataType.DefaultValue)

	assert.Equal(t, "Field3", actGi.Fields[2].Name.String())
	assert.Equal(t, "p2.MyStruct2", actGi.Fields[2].DataType.Name.String())
	assert.Equal(t, `p2.MyStruct2{}`, actGi.Fields[2].DataType.DefaultValue)

	assert.Equal(t, "Field4", actGi.Fields[3].Name.String())
	assert.Equal(t, "p2.MyAliasBool", actGi.Fields[3].DataType.Name.String())
	assert.Equal(t, `p2.MyAliasBool(false)`, actGi.Fields[3].DataType.DefaultValue)

	assert.Equal(t, "Field5", actGi.Fields[4].Name.String())
	assert.Equal(t, "[]string", actGi.Fields[4].DataType.Name.String())
	assert.Equal(t, `[]string{}`, actGi.Fields[4].DataType.DefaultValue)

	assert.Equal(t, "Field6", actGi.Fields[5].Name.String())
	assert.Equal(t, "[]*p2.MyStruct2", actGi.Fields[5].DataType.Name.String())
	assert.Equal(t, `[]*p2.MyStruct2{}`, actGi.Fields[5].DataType.DefaultValue)

	assert.Equal(t, "Field7", actGi.Fields[6].Name.String())
	assert.Equal(t, "p2.MyInterface1", actGi.Fields[6].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Fields[6].DataType.DefaultValue)

	assert.Equal(t, "Field8", actGi.Fields[7].Name.String())
	assert.Equal(t, "struct{x int; y string}", actGi.Fields[7].DataType.Name.String())
	assert.Equal(t, `struct{x int; y string}{}`, actGi.Fields[7].DataType.DefaultValue)

	assert.Equal(t, "Field9", actGi.Fields[8].Name.String())
	assert.Equal(t, "func(string, ) int", actGi.Fields[8].DataType.Name.String())
	assert.Equal(t, `nil`, actGi.Fields[8].DataType.DefaultValue)

	assert.Equal(t, "Field10", actGi.Fields[9].Name.String())
	assert.Equal(t, "map[string]int", actGi.Fields[9].DataType.Name.String())
	assert.Equal(t, `map[string]int{}`, actGi.Fields[9].DataType.DefaultValue)

	//

	assert.Equal(t, 3, len(actGi.Imports))

	assert.Equal(t, ``, actGi.Imports["sync"].Name)
	assert.Equal(t, "sync", actGi.Imports["sync"].Expression.String())
	assert.Equal(t, "sync", actGi.Imports["sync"].Path.String())
	assert.Equal(t, "GO_SDK", actGi.Imports["sync"].ImportType.String())

	assert.Equal(t, ``, actGi.Imports["p2"].Name)
	assert.Equal(t, "p2", actGi.Imports["p2"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/struct001/p2", actGi.Imports["p2"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.Imports["p2"].ImportType.String())

	assert.Equal(t, ``, actGi.Imports["p3"].Name)
	assert.Equal(t, "p3", actGi.Imports["p3"].Expression.String())
	assert.Equal(t, "mirza/gogen/refactor/struct001/p3", actGi.Imports["p3"].Path.String())
	assert.Equal(t, "INTERNAL_PROJECT", actGi.Imports["p3"].ImportType.String())

	assert.Equal(t, 3, len(actGi.CompositionTypes))
	assert.Equal(t, "MyStruct3", actGi.CompositionTypes[0].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[0].Methods))

	//

	assert.Equal(t, 1, len(actGi.CompositionTypes[0].Fields))

	assert.Equal(t, "Field1", actGi.CompositionTypes[0].Fields[0].Name.String())
	assert.Equal(t, "bool", actGi.CompositionTypes[0].Fields[0].DataType.Name.String())
	assert.Equal(t, `false`, actGi.CompositionTypes[0].Fields[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[0].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[0].CompositionTypes))
	assert.Equal(t, "p3.MyStruct4", actGi.CompositionTypes[1].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[1].Methods))

	//

	assert.Equal(t, 1, len(actGi.CompositionTypes[1].Fields))

	assert.Equal(t, "Field1", actGi.CompositionTypes[1].Fields[0].Name.String())
	assert.Equal(t, "float64", actGi.CompositionTypes[1].Fields[0].DataType.Name.String())
	assert.Equal(t, `0.0`, actGi.CompositionTypes[1].Fields[0].DataType.DefaultValue)

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[1].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[1].CompositionTypes))
	assert.Equal(t, "*p3.MyStruct5", actGi.CompositionTypes[2].Name.String())
	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Methods))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Fields))

	//

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].Imports))

	assert.Equal(t, 0, len(actGi.CompositionTypes[2].CompositionTypes))

}
