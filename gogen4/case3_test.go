package gogen4

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase201(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t201/p1/usecase/runordercreate", "./data_testing/project/go.mod", "Outport")

	assert.Nil(t, err)

	assert.Equal(t, "Outport", actGi.GogenFieldType.Name.String())
	assert.Equal(t, "SaveOrderRepo", actGi.CompositionTypes[0].GogenFieldType.Name.String())
}

func TestCase202(t *testing.T) {

	gi := NewGogenAnyTypeBuilder()
	actGi, err := gi.Build("./data_testing/project/t202/p1", "./data_testing/project/go.mod", "MyStruct1")

	assert.Nil(t, err)

	assert.Equal(t, "MyStruct1", actGi.GogenFieldType.Name.String())
}
