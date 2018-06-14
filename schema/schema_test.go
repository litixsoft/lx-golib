package lxSchema

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const (
	cSCHEMAROOTPATH            = "./../data"
	cSCHEMA_EXISTS_FILENAME    = "schema_001.json"
	cSCHEMA_NOTEXISTS_FILENAME = "schema_002.json"
)

func TestJSONSchemaStruct_SetSchemaRootDirectory(t *testing.T) {
	t.Run("Set root schema directory", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)
	})
}

func TestJSONSchemaStruct_LoadSchema(t *testing.T) {
	t.Run("Load a json schema", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)

		json, err := Loader.LoadSchema(cSCHEMA_EXISTS_FILENAME)
		assert.NotNil(t, json, "LoadSchema returns nil")
		assert.NoError(t, err)
	})

	t.Run("Return error if no json schema loader", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)

		json, err := Loader.LoadSchema(cSCHEMA_NOTEXISTS_FILENAME)
		assert.Nil(t, json, "LoadSchema return not nil")
		assert.Error(t, err)
	})
}

func TestJSONSchemaStruct_HasSchema(t *testing.T) {
	t.Run("Check if loaded schema exists", func(t *testing.T) {
		res := Loader.HasSchema(cSCHEMA_EXISTS_FILENAME)
		assert.Equal(t, res, true)
	})

	t.Run("Check if not loaded schema not exists", func(t *testing.T) {
		res := Loader.HasSchema(cSCHEMA_NOTEXISTS_FILENAME)
		assert.Equal(t, res, false)
	})
}

func TestJSONSchemaStruct_ValidateBind(t *testing.T) {
	t.Run("Error if no schema given", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := Loader.ValidateBind("", nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
	})

	t.Run("Error if schema not exists", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := Loader.ValidateBind(cSCHEMA_NOTEXISTS_FILENAME, nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
	})

	t.Run("Error if no context given", func(t *testing.T) {
		err := Loader.SetSchemaRootDirectory(cSCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := Loader.ValidateBind(cSCHEMA_EXISTS_FILENAME, nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("context is required"))
	})
}
