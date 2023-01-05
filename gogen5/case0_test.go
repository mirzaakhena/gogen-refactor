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
