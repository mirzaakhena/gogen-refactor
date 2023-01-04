package gogen5

import (
	"gen/gogen5/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase001(t *testing.T) {

	actGi, err := Build("./data_testing/project/interface001/p1", "./data_testing/project/go.mod", "MyInterface1")

	util.PrintGogenAnyType(1, actGi)

	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.Name.String())
	assert.Nil(t, actGi.GogenFieldType)
	assert.Equal(t, 3, len(actGi.Methods))

	{
		assert.Equal(t, "Method11", actGi.Methods[0].Name.String())

		assert.Equal(t, 2, len(actGi.Methods[0].Params))
		assert.Equal(t, 2, len(actGi.Methods[0].Results))

		assert.Equal(t, "x", actGi.Methods[0].Params[0].Name.String())
		assert.Equal(t, "int", actGi.Methods[0].Params[0].DataType.Name.String())
		assert.Equal(t, "0", actGi.Methods[0].Params[0].DataType.DefaultValue)

		assert.Equal(t, "y", actGi.Methods[0].Params[1].Name.String())
		assert.Equal(t, "string", actGi.Methods[0].Params[1].DataType.Name.String())
		assert.Equal(t, `""`, actGi.Methods[0].Params[1].DataType.DefaultValue)

		assert.Equal(t, "bool", actGi.Methods[0].Results[0].Name.String())
		assert.Equal(t, "bool", actGi.Methods[0].Results[0].DataType.Name.String())
		assert.Equal(t, "false", actGi.Methods[0].Results[0].DataType.DefaultValue)

		assert.Equal(t, "error", actGi.Methods[0].Results[1].Name.String())
		assert.Equal(t, "error", actGi.Methods[0].Results[1].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[0].Results[1].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Method12", actGi.Methods[1].Name.String())

		assert.Equal(t, 2, len(actGi.Methods[1].Params))
		assert.Equal(t, 2, len(actGi.Methods[1].Results))

		assert.Equal(t, "int", actGi.Methods[1].Params[0].Name.String())
		assert.Equal(t, "int", actGi.Methods[1].Params[0].DataType.Name.String())
		assert.Equal(t, "0", actGi.Methods[1].Params[0].DataType.DefaultValue)

		assert.Equal(t, "string", actGi.Methods[1].Params[1].Name.String())
		assert.Equal(t, "string", actGi.Methods[1].Params[1].DataType.Name.String())
		assert.Equal(t, `""`, actGi.Methods[1].Params[1].DataType.DefaultValue)

		assert.Equal(t, "x", actGi.Methods[1].Results[0].Name.String())
		assert.Equal(t, "bool", actGi.Methods[1].Results[0].DataType.Name.String())
		assert.Equal(t, "false", actGi.Methods[1].Results[0].DataType.DefaultValue)

		assert.Equal(t, "y", actGi.Methods[1].Results[1].Name.String())
		assert.Equal(t, "error", actGi.Methods[1].Results[1].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[1].Results[1].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Method13", actGi.Methods[2].Name.String())

		assert.Equal(t, 2, len(actGi.Methods[2].Params))
		assert.Equal(t, 0, len(actGi.Methods[2].Results))

		assert.Equal(t, "ctx3", actGi.Methods[2].Params[0].Name.String())
		assert.Equal(t, "context.Context", actGi.Methods[2].Params[0].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[2].Params[0].DataType.DefaultValue)

		assert.Equal(t, "handler", actGi.Methods[2].Params[1].Name.String())
		assert.Equal(t, "gin.HandlerFunc", actGi.Methods[2].Params[1].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[2].Params[1].DataType.DefaultValue)
	}

	assert.Equal(t, 5, len(actGi.CompositionTypes))

	assert.Equal(t, "MyInterface2", actGi.CompositionTypes[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[0].Methods))
	assert.Equal(t, "Method21", actGi.CompositionTypes[0].Methods[0].Name.String())

	assert.Equal(t, "MyInterface3", actGi.CompositionTypes[1].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[1].Methods))
	assert.Equal(t, "Method31", actGi.CompositionTypes[1].Methods[0].Name.String())

	assert.Equal(t, "MyInterface4", actGi.CompositionTypes[2].Name.String())
	//assert.Equal(t, "Method41", actGi.CompositionTypes[2].Methods[0].Name.String())

	assert.Equal(t, "MyInterface5", actGi.CompositionTypes[3].Name.String())
	//assert.Equal(t, "Method51", actGi.CompositionTypes[3].Methods[0].Name.String())

	assert.Equal(t, "p2.MyInterface6", actGi.CompositionTypes[4].Name.String())
	//assert.Equal(t, "Method61", actGi.CompositionTypes[4].Methods[0].Name.String())

}
