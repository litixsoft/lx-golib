package lxSchema_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/litixsoft/lx-golib/schema"
	"bytes"
	"github.com/litixsoft/lx-golib/test-helper"
	"github.com/labstack/echo"
	"github.com/litixsoft/lx-golib/helper"
	"encoding/json"
)

const (
	SCHEMAROOTPATH            = "./../data"
	SCHEMA_EXISTS_FILENAME    = "schema_001.json"
	SCHEMA_NOTEXISTS_FILENAME = "schema_002.json"
)

func TestJSONSchemaStruct_SetSchemaRootDirectory(t *testing.T) {
	t.Run("Set root schema directory", func(t *testing.T) {

		lxSchema.InitJsonSchemaLoader()
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)
	})
}

func TestJSONSchemaStruct_LoadSchema(t *testing.T) {
	t.Run("Load a json schema", func(t *testing.T) {
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		json, err := lxSchema.Loader.LoadSchema(SCHEMA_EXISTS_FILENAME)
		assert.NotNil(t, json, "LoadSchema returns nil")
		assert.NoError(t, err)
	})

	t.Run("Return error if no json schema loader", func(t *testing.T) {
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		json, err := lxSchema.Loader.LoadSchema(SCHEMA_NOTEXISTS_FILENAME)
		assert.Nil(t, json, "LoadSchema return not nil")
		assert.Error(t, err)
	})
}

func TestJSONSchemaStruct_HasSchema(t *testing.T) {
	t.Run("Check if loaded schema exists", func(t *testing.T) {
		res := lxSchema.Loader.HasSchema(SCHEMA_EXISTS_FILENAME)
		assert.Equal(t, res, true)
	})

	t.Run("Check if not loaded schema not exists", func(t *testing.T) {
		res := lxSchema.Loader.HasSchema(SCHEMA_NOTEXISTS_FILENAME)
		assert.Equal(t, res, false)
	})
}

func TestJSONSchemaStruct_ValidateBind(t *testing.T) {
	t.Run("Error if no schema given", func(t *testing.T) {
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := lxSchema.Loader.ValidateBind("", nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
	})

	t.Run("Error if schema not exists", func(t *testing.T) {
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := lxSchema.Loader.ValidateBind(SCHEMA_NOTEXISTS_FILENAME, nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
	})

	t.Run("Error if no context given", func(t *testing.T) {
		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		res, err := lxSchema.Loader.ValidateBind(SCHEMA_EXISTS_FILENAME, nil, nil)

		assert.Nil(t, res, "Schema was loaded")
		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("context is required"))
	})
	t.Run("Error if no context given", func(t *testing.T) {
		lxSchema.InitJsonSchemaLoader()

		err := lxSchema.Loader.SetSchemaRootDirectory(SCHEMAROOTPATH)
		assert.NoError(t, err)

		// Convert login data to json
		user := lxHelper.M{"name": "Otto", "login_name": "otto", "email":"otto@otto.com"}

		jsonData, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		// Set request
		_, c := lxTestHelper.SetEchoRequest(echo.POST, "/schema", bytes.NewBuffer(jsonData))


		res, err := lxSchema.Loader.ValidateBind(SCHEMA_EXISTS_FILENAME, c, nil)

		t.Log(res.Valid())

		//assert.Nil(t, res, "Schema was loaded")
		//assert.Error(t, err)
		//assert.Equal(t, err, fmt.Errorf("context is required"))
	})
}


