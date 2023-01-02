package gogen4

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase001(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t001/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "p1", actGi.GogenFieldType.File.Name.String())
	assert.Equal(t, "nil", actGi.GogenFieldType.DefaultValue)

}

func TestCase002(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t002/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "p1", actGi.GogenFieldType.File.Name.String())
	assert.Equal(t, "nil", actGi.GogenFieldType.DefaultValue)

	assert.Equal(t, 3, len(actGi.Methods))

	{
		assert.Equal(t, "Method1", actGi.Methods[0].Name.String())

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
		assert.Equal(t, "Method2", actGi.Methods[1].Name.String())

		assert.Equal(t, "ctx3", actGi.Methods[1].Params[0].Name.String())
		assert.Equal(t, "context.Context", actGi.Methods[1].Params[0].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[1].Params[0].DataType.DefaultValue)

		assert.Equal(t, "handler", actGi.Methods[1].Params[1].Name.String())
		assert.Equal(t, "gin.HandlerFunc", actGi.Methods[1].Params[1].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[1].Params[1].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Method3", actGi.Methods[2].Name.String())

		assert.Equal(t, "int", actGi.Methods[2].Params[0].Name.String())
		assert.Equal(t, "int", actGi.Methods[2].Params[0].DataType.Name.String())
		assert.Equal(t, "0", actGi.Methods[2].Params[0].DataType.DefaultValue)

		assert.Equal(t, "string", actGi.Methods[2].Params[1].Name.String())
		assert.Equal(t, "string", actGi.Methods[2].Params[1].DataType.Name.String())
		assert.Equal(t, `""`, actGi.Methods[2].Params[1].DataType.DefaultValue)

		assert.Equal(t, "x", actGi.Methods[2].Results[0].Name.String())
		assert.Equal(t, "bool", actGi.Methods[2].Results[0].DataType.Name.String())
		assert.Equal(t, "false", actGi.Methods[2].Results[0].DataType.DefaultValue)

		assert.Equal(t, "y", actGi.Methods[2].Results[1].Name.String())
		assert.Equal(t, "error", actGi.Methods[2].Results[1].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[2].Results[1].DataType.DefaultValue)

	}

}

func TestCase012(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t012/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "p1", actGi.GogenFieldType.File.Name.String())
	assert.Equal(t, "nil", actGi.GogenFieldType.DefaultValue)

	assert.Equal(t, 1, len(actGi.Methods))

	{
		assert.Equal(t, "Method1", actGi.Methods[0].Name.String())

		assert.Equal(t, "handler", actGi.Methods[0].Params[0].Name.String())
		assert.Equal(t, "gin.HandlerFunc", actGi.Methods[0].Params[0].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Methods[0].Params[0].DataType.DefaultValue)

	}

}

func TestCase003(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t003/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "Method1", actGi.Methods[0].Name.String())

	assert.Equal(t, 4, len(actGi.CompositionTypes))

	assert.Equal(t, "MyInterface2", actGi.CompositionTypes[0].GogenFieldType.Name.String())
	assert.Equal(t, "Method2", actGi.CompositionTypes[0].Methods[0].Name.String())

	assert.Equal(t, "MyInterface3", actGi.CompositionTypes[1].GogenFieldType.Name.String())
	assert.Equal(t, "Method3", actGi.CompositionTypes[1].Methods[0].Name.String())

	assert.Equal(t, "MyInterface4", actGi.CompositionTypes[2].GogenFieldType.Name.String())
	assert.Equal(t, "Method4", actGi.CompositionTypes[2].Methods[0].Name.String())

	assert.Equal(t, "MyInterface5", actGi.CompositionTypes[3].GogenFieldType.Name.String())
	assert.Equal(t, "Method5", actGi.CompositionTypes[3].Methods[0].Name.String())

}

func TestCase004(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t004/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "MyInterface2", actGi.CompositionTypes[0].GogenFieldType.Name.String())
	assert.Equal(t, "Method2", actGi.CompositionTypes[0].Methods[0].Name.String())

}

func TestCase005(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t005/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "Binder", actGi.CompositionTypes[0].GogenFieldType.Name.String())
	assert.Equal(t, "Bind", actGi.CompositionTypes[0].Methods[0].Name.String())

	assert.Equal(t, "i", actGi.CompositionTypes[0].Methods[0].Params[0].Name.String())
	assert.Equal(t, "any", actGi.CompositionTypes[0].Methods[0].Params[0].DataType.Name.String())
	assert.Equal(t, "nil", actGi.CompositionTypes[0].Methods[0].Params[0].DataType.DefaultValue)

	assert.Equal(t, "c", actGi.CompositionTypes[0].Methods[0].Params[1].Name.String())
	assert.Equal(t, "Context", actGi.CompositionTypes[0].Methods[0].Params[1].DataType.Name.String())
	assert.Equal(t, "nil", actGi.CompositionTypes[0].Methods[0].Params[1].DataType.DefaultValue)

	assert.Equal(t, "error", actGi.CompositionTypes[0].Methods[0].Results[0].Name.String())
	assert.Equal(t, "error", actGi.CompositionTypes[0].Methods[0].Results[0].DataType.Name.String())
	assert.Equal(t, "nil", actGi.CompositionTypes[0].Methods[0].Results[0].DataType.DefaultValue)

}

func TestCase006(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/t006/p1", "./data_testing/project/go.mod", "MyInterface1")
	assert.Nil(t, err)

	assert.Equal(t, "MyInterface1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "TypeAlias", actGi.CompositionTypes[0].GogenFieldType.Name.String())
	assert.Equal(t, "MyInterface2", actGi.CompositionTypes[0].CompositionTypes[0].GogenFieldType.Name.String())

}
