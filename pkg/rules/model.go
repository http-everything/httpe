package rules

import "http-everything/httpe/pkg/share/logger"

type Rules struct {
	Rules     *[]Rule
	logger    *logger.Logger
	rulesFile string
}

const SchemaURL = "https://github.com/http-everything/httpe/main/pkg/rules/schema.json"

type Rule struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	On   *On    `yaml:"on,omitempty" json:"on,omitempty"`
	Do   *Do    `yaml:"do,omitempty" json:"do,omitempty"`
}

type On struct {
	Path    string   `yaml:"path,omitempty" json:"path,omitempty"`
	Methods []string `yaml:"methods,omitempty" json:"methods,omitempty"`
}

type Do struct {
	Name    string `yaml:"name,omitempty" json:"name,omitempty"`
	Execute string `yaml:"execute,omitempty" json:"execute,omitempty"`
	Args    *Args  `yaml:"args,omitempty" json:"args,omitempty"`
}

type Args struct {
	Interpreter string  `yaml:"interpreter,omitempty" json:"interpreter,omitempty"`
	Timeout     float64 `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

type With struct {
	AuthBasic *AuthBasic `yaml:"auth_basic,omitempty" json:"auth_basic,omitempty"`
}

type AuthBasic struct {
	Users []User
}

type User struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Hash     string `yaml:"hash" json:"hash"`
}
