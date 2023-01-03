package gogen4

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase101(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t101/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())
}

func TestCase102(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t102/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())

	assert.Equal(t, 5, len(actGi.Fields))

	assert.Equal(t, "Name", actGi.Fields[0].Name.String())
	assert.Equal(t, "string", actGi.Fields[0].DataType.Name.String())
	assert.Equal(t, `""`, actGi.Fields[0].DataType.DefaultValue)

	assert.Equal(t, "Age", actGi.Fields[1].Name.String())
	assert.Equal(t, "int", actGi.Fields[1].DataType.Name.String())
	assert.Equal(t, "0", actGi.Fields[1].DataType.DefaultValue)

	assert.Equal(t, "Hobbies", actGi.Fields[2].Name.String())
	assert.Equal(t, "[]string", actGi.Fields[2].DataType.Name.String())
	assert.Equal(t, "[]string{}", actGi.Fields[2].DataType.DefaultValue)

	assert.Equal(t, "Wife", actGi.Fields[3].Name.String())
	assert.Equal(t, "*MyStruct1", actGi.Fields[3].DataType.Name.String())
	assert.Equal(t, "nil", actGi.Fields[3].DataType.DefaultValue)

	assert.Equal(t, "ThisStruct", actGi.Fields[4].Name.String())
	assert.Equal(t, "struct{x int; y float64}", actGi.Fields[4].DataType.Name.String())
	assert.Equal(t, "struct{x int; y float64}{}", actGi.Fields[4].DataType.DefaultValue)

}

func TestCase103(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t103/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())

	assert.Equal(t, 2, len(actGi.Fields))

	assert.Equal(t, "ctx8", actGi.Fields[0].Name.String())
	assert.Equal(t, "context.Context", actGi.Fields[0].DataType.Name.String())
	assert.Equal(t, "nil", actGi.Fields[0].DataType.DefaultValue)

	assert.Equal(t, "theHandler", actGi.Fields[1].Name.String())
	assert.Equal(t, "echo.HandlerFunc", actGi.Fields[1].DataType.Name.String())
	assert.Equal(t, "nil", actGi.Fields[1].DataType.DefaultValue)

}

func TestCase104(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t104/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())

	assert.Equal(t, 5, len(actGi.Fields))

	assert.Equal(t, "A01", actGi.Fields[0].Name.String())
	assert.Equal(t, "MyStruct2", actGi.Fields[0].DataType.Name.String())
	assert.Equal(t, "MyStruct2{}", actGi.Fields[0].DataType.DefaultValue)

	assert.Equal(t, "A02", actGi.Fields[1].Name.String())
	assert.Equal(t, "MyStruct3", actGi.Fields[1].DataType.Name.String())
	assert.Equal(t, "MyStruct3{}", actGi.Fields[1].DataType.DefaultValue)

	assert.Equal(t, "A03", actGi.Fields[2].Name.String())
	assert.Equal(t, "p1.MyStruct4", actGi.Fields[2].DataType.Name.String())
	assert.Equal(t, "p1.MyStruct4{}", actGi.Fields[2].DataType.DefaultValue)

	assert.Equal(t, "A04", actGi.Fields[3].Name.String())
	assert.Equal(t, "MyStruct5", actGi.Fields[3].DataType.Name.String())
	assert.Equal(t, "MyStruct5{}", actGi.Fields[3].DataType.DefaultValue)

	assert.Equal(t, "A05", actGi.Fields[4].Name.String())
	assert.Equal(t, "MyStruct6", actGi.Fields[4].DataType.Name.String())
	assert.Equal(t, "MyStruct6{}", actGi.Fields[4].DataType.DefaultValue)

}

func TestCase105(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t105/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "Age", actGi.Fields[0].Name.String())

	assert.Equal(t, 6, len(actGi.CompositionTypes))

	assert.Equal(t, "MyStruct2", actGi.CompositionTypes[0].GogenFieldType.Name.String())
	assert.Equal(t, "Name", actGi.CompositionTypes[0].Fields[0].Name.String())

	assert.Equal(t, "MyStruct3", actGi.CompositionTypes[1].GogenFieldType.Name.String())
	assert.Equal(t, "MyStruct4", actGi.CompositionTypes[2].GogenFieldType.Name.String())
	assert.Equal(t, "MyStruct5", actGi.CompositionTypes[3].GogenFieldType.Name.String())
	assert.Equal(t, "YourInterface", actGi.CompositionTypes[4].GogenFieldType.Name.String())
	assert.Equal(t, "MyStruct6", actGi.CompositionTypes[5].GogenFieldType.Name.String())

}
