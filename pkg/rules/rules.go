package rules

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"http-everything/httpe/pkg/share/logger"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

//go:embed schema.json
var schemaJSON string

func New(logger *logger.Logger) (rules *Rules) {
	rules = &Rules{
		logger: logger,
	}
	return rules
}

func (r *Rules) LoadAndValidate(rulesFile string) (rules *[]Rule, err error) {
	r.Rules = &[]Rule{}
	r.rulesFile = rulesFile
	rulesYaml, err := os.ReadFile(rulesFile)
	if err != nil {
		return r.Rules, err
	}

	err = yaml.Unmarshal(rulesYaml, r.Rules)
	if err != nil {
		return r.Rules, err
	}

	err = r.validate()
	if err != nil {
		return r.Rules, err
	}

	return r.Rules, nil
}

func (r *Rules) validate() (err error) {
	rulesJSON, err := json.Marshal(r.Rules)
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	documentLoader := gojsonschema.NewStringLoader(string(rulesJSON))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("json schema valdation failed: %w", err)
	}

	if result.Valid() {
		r.logger.Debugf("'%s' successfully validated against schema", r.rulesFile)
	} else {
		r.logger.Errorf("schema validation against %s failed", SchemaURL)
		for _, desc := range result.Errors() {
			r.logger.Errorf("%s\n", desc)
		}
		return fmt.Errorf("invalid rules file '%s'", r.rulesFile)
	}
	return nil
}
