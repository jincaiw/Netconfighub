package devicemodel

type DeviceModelConfig struct {
	Model         string       `yaml:"model" json:"model"`
	Vendor        string       `yaml:"vendor" json:"vendor"`
	Protocols     []string     `yaml:"protocols" json:"protocols"`
	DisablePaging []string     `yaml:"disable_paging" json:"disable_paging"`
	Commands      []string     `yaml:"commands" json:"commands"`
	Exit          []string     `yaml:"exit" json:"exit"`
	PromptRegex   string       `yaml:"prompt_regex" json:"prompt_regex"`
	EnableCmd     string       `yaml:"enable_cmd" json:"enable_cmd"`
	SecretRules   []SecretRule `yaml:"secret_rules" json:"secret_rules"`
}

type SecretRule struct {
	Pattern string `yaml:"pattern" json:"pattern"`
	Replace string `yaml:"replace" json:"replace"`
}
