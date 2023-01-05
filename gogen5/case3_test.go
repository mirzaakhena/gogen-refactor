package gogen5

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase003(t *testing.T) {

	actGi, err := Build("./data_testing/project/alias001/p1", "./data_testing/project/go.mod", "MyAlias1")

	PrintGogenAnyType(0, actGi)

	jsonInBytes, err := json.Marshal(actGi)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(jsonInBytes))

	assert.Nil(t, err)

	assert.Equal(t, "MyAlias1", actGi.Name.String())
	assert.Equal(t, 0, len(actGi.Fields))
	assert.Equal(t, 0, len(actGi.Methods))
	assert.Equal(t, 1, len(actGi.CompositionTypes))
	assert.Equal(t, 0, len(actGi.Imports))

	assert.Equal(t, "MyStruct1", actGi.CompositionTypes[0].Name.String())
	assert.Equal(t, 1, len(actGi.CompositionTypes[0].Fields))
	assert.Equal(t, "X", actGi.CompositionTypes[0].Fields[0].Name.String())
	assert.Equal(t, "int", actGi.CompositionTypes[0].Fields[0].DataType.Name.String())
	assert.Equal(t, "0", actGi.CompositionTypes[0].Fields[0].DataType.DefaultValue)

}
