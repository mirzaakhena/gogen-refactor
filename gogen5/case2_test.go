package gogen5

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase002(t *testing.T) {

	actGi, err := Build("./data_testing/project/struct001/p1", "./data_testing/project/go.mod", "MyStruct1")

	PrintGogenAnyType(0, actGi)

	jsonInBytes, err := json.Marshal(actGi)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(jsonInBytes))

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.Name.String())
	assert.Equal(t, 10, len(actGi.Fields))
	assert.Equal(t, 0, len(actGi.Methods))
	assert.Equal(t, 0, len(actGi.CompositionTypes))
	assert.Equal(t, 2, len(actGi.Imports))

	{
		assert.Equal(t, "Field1", actGi.Fields[0].Name.String())
		assert.Equal(t, "int", actGi.Fields[0].DataType.Name.String())
		assert.Equal(t, "0", actGi.Fields[0].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field2", actGi.Fields[1].Name.String())
		assert.Equal(t, "sync.WaitGroup", actGi.Fields[1].DataType.Name.String())
		assert.Equal(t, "sync.WaitGroup{}", actGi.Fields[1].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field3", actGi.Fields[2].Name.String())
		assert.Equal(t, "p2.MyStruct2", actGi.Fields[2].DataType.Name.String())
		assert.Equal(t, "p2.MyStruct2{}", actGi.Fields[2].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field4", actGi.Fields[3].Name.String())
		assert.Equal(t, "p2.MyAliasBool", actGi.Fields[3].DataType.Name.String())
		assert.Equal(t, "p2.MyAliasBool(false)", actGi.Fields[3].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field5", actGi.Fields[4].Name.String())
		assert.Equal(t, "[]string", actGi.Fields[4].DataType.Name.String())
		assert.Equal(t, "[]string{}", actGi.Fields[4].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field6", actGi.Fields[5].Name.String())
		assert.Equal(t, "[]*p2.MyStruct2", actGi.Fields[5].DataType.Name.String())
		assert.Equal(t, "[]*p2.MyStruct2{}", actGi.Fields[5].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field7", actGi.Fields[6].Name.String())
		assert.Equal(t, "p2.MyInterface1", actGi.Fields[6].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Fields[6].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field8", actGi.Fields[7].Name.String())
		assert.Equal(t, "struct{x int; y string}", actGi.Fields[7].DataType.Name.String())
		assert.Equal(t, "struct{x int; y string}{}", actGi.Fields[7].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field9", actGi.Fields[8].Name.String())
		assert.Equal(t, "func(string, ) int", actGi.Fields[8].DataType.Name.String())
		assert.Equal(t, "nil", actGi.Fields[8].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "Field10", actGi.Fields[9].Name.String())
		assert.Equal(t, "map[string]int", actGi.Fields[9].DataType.Name.String())
		assert.Equal(t, "map[string]int{}", actGi.Fields[9].DataType.DefaultValue)
	}

	{
		assert.Equal(t, "", actGi.Imports["p2"].Name)
		assert.Equal(t, Expression("p2"), actGi.Imports["p2"].Expression)
		assert.Equal(t, ImportPath("mirza/gogen/refactor/struct001/p2"), actGi.Imports["p2"].Path)
		assert.Equal(t, ImportType("INTERNAL_PROJECT"), actGi.Imports["p2"].ImportType)
	}

	{
		assert.Equal(t, "", actGi.Imports["sync"].Name)
		assert.Equal(t, Expression("sync"), actGi.Imports["sync"].Expression)
		assert.Equal(t, ImportPath("sync"), actGi.Imports["sync"].Path)
		assert.Equal(t, ImportType("GO_SDK"), actGi.Imports["sync"].ImportType)
	}

}
