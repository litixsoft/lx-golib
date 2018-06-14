package lxSchema

import (
	"github.com/xeipuuv/gojsonschema"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/labstack/echo"
	"path/filepath"
)

type
	JSONSchemaStruct struct {
		root    string
		schemas map[string]gojsonschema.JSONLoader
	}

var (
	Loader = &JSONSchemaStruct{
		root:    "",
		schemas: make(map[string]gojsonschema.JSONLoader),
	}
)

func (self *JSONSchemaStruct) SetSchemaRootDirectory(dirname string) (error) {
	var err error = nil

	if !filepath.IsAbs(dirname) {
		dirname, err = filepath.Abs(dirname)
	}

	if err == nil {
		self.root = "file:///" + filepath.ToSlash(dirname) + "/"
	}

	return err
}

func (self *JSONSchemaStruct) HasSchema(filename string) (bool) {
	return self.schemas[filename] != nil
}

func (self *JSONSchemaStruct) LoadSchema(filename string) (gojsonschema.JSONLoader, error) {
	if !self.HasSchema(filename) {
		var jsonURILoader = gojsonschema.NewReferenceLoader(self.root + filename)

		if _, err := jsonURILoader.LoadJSON(); err != nil {
			return nil, err
		}

		self.schemas[filename] = jsonURILoader
	}

	return self.schemas[filename], nil
}

func (self *JSONSchemaStruct) ValidateBind(schema string, c echo.Context, s interface{}) (*gojsonschema.Result, error) {
	schemaLoader, err := self.LoadSchema(schema)

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
