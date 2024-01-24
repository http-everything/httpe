package rules

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"http-everything/httpe/pkg/share/logger"
	"os"
	"strings"

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

func (r *Rules) Load(rulesFile string) (rules *[]Rule, err error) {
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

	return r.Rules, nil
}

func (r *Rules) Validate() (err error) {
	rulesJSON, err := json.Marshal(r.Rules)
	if err != nil {
		return err
	}

	// Validate the rule action. Doing it here and before the schema validation creates a more user-friendly error
	// message. While it is possible to do the below validation with the 'oneOf' statement of JSON schema validation,
	// the error message is confusing because only the first element of the oneOf list is named as required.
	var hasErrors = false
	for i, rule := range *r.Rules {
		if rule.Action() == "" {
			r.logger.Errorf("rule %d '%s' is missing a valid action in the 'do' section. Use one of '%s'.",
				i,
				rule.Name, strings.Join(ValidActions, ", "))
			hasErrors = true
		}
	}
	if hasErrors {
		return fmt.Errorf("invalid rules file '%s'", r.rulesFile)
	}

	// Validate against schema
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	documentLoader := gojsonschema.NewStringLoader(string(rulesJSON))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("json schema validation failed: %w", err)
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

func (r *Rules) AsJSONString() string {
	rulesJSON, err := json.MarshalIndent(r.Rules, "", "  ")
	if err != nil {
		return ""
	}
	return string(rulesJSON)
}

func (rule *Rule) Action() (action string) {
	if rule.Do == nil {
		return ""
	}
	if rule.Do.RunScript != "" {
		return RunScript
	}
	if rule.Do.SendEmail != "" {
		return SendEmail
	}
	if rule.Do.AnswerContent != "" {
		return AnswerContent
	}
	if rule.Do.AnswerFile != "" {
		return AnswerFile
	}
	return ""
}
