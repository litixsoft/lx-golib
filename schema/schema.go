package lxSchema

import (
	"github.com/xeipuuv/gojsonschema"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/labstack/echo"
	"path/filepath"
)

type IJSONSchema interface {
	SetSchemaRootDirectory(dirname string) error
	HasSchema(filename string) bool
	LoadSchema(filename string) (gojsonschema.JSONLoader, error)
	ValidateBind(schema string, c echo.Context, s interface{}) (*gojsonschema.Result, error)
}

type JSONSchema struct {
	root    string
	schemas map[string]gojsonschema.JSONLoader
}

var Loader IJSONSchema

func InitJsonSchemaLoader() {
	Loader = &JSONSchema{
		root:    "",
		schemas: make(map[string]gojsonschema.JSONLoader),
	}
}

func (js *JSONSchema) SetSchemaRootDirectory(dirname string) error {
	var err error = nil

	if !filepath.IsAbs(dirname) {
		dirname, err = filepath.Abs(dirname)
	}

	if err == nil {
		js.root = "file:///" + filepath.ToSlash(dirname) + "/"
	}

	return err
}

func (js *JSONSchema) HasSchema(filename string) bool {
	return js.schemas[filename] != nil
}

func (js *JSONSchema) LoadSchema(filename string) (gojsonschema.JSONLoader, error) {
	if !js.HasSchema(filename) {
		var jsonURILoader = gojsonschema.NewReferenceLoader(js.root + filename)

		if _, err := jsonURILoader.LoadJSON(); err != nil {
			return nil, err
		}

		js.schemas[filename] = jsonURILoader
	}

	return js.schemas[filename], nil
}

func (js *JSONSchema) ValidateBind(schema string, c echo.Context, s interface{}) (*gojsonschema.Result, error) {
	schemaLoader, err := js.LoadSchema(schema)

	if err != nil {
		return nil, err
	}

	// No schema
	if schemaLoader == nil {
		return nil, fmt.Errorf("schema could not be loaded %s", schema)
	}

	if c == nil {
		return nil, fmt.Errorf("context is required")
	}

	var (
		jsonRAWDoc []byte
	)

	// Read []byte stream from Request Body
	if jsonRAWDoc, err = ioutil.ReadAll(c.Request().Body); err != nil {
		return nil, err
	}

	// Transfer []byte stream to gojsonschema.documentLoader and Validate
	documentLoader := gojsonschema.NewBytesLoader(jsonRAWDoc)
	res, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if s != nil && err == nil && res.Valid() {
		// if schema valid; and S given; translate
		if err = json.Unmarshal(documentLoader.JsonSource().([]byte), s); err != nil {
			return nil, err
		}
	}

	// Return result
	return res, err
}