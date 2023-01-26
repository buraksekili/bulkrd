package main

type BetweenCount struct {
	MinCount int `yaml:"minCount" json:"minCount"`
	MaxCount int `yaml:"maxCount" json:"maxCount"`
}

type Count struct {
	Exact   int          `yaml:"exact" json:"exact,omitempty"`
	Between BetweenCount `yaml:"between" json:"between,omitempty"`
}

type RandomField struct {
	Field string `yaml:"field" json:"field"`
	Type  string `yaml:"type" json:"type"`
}

type Resource struct {
	TemplatePath string        `yaml:"templatePath" json:"templatePath"`
	Namespace    string        `yaml:"namespace" json:"namespace"`
	Count        Count         `yaml:"count" json:"count"`
	Randomize    []RandomField `yaml:"randomize" json:"randomize,omitempty"`
}

type Config struct {
	Debug     bool       `yaml:"debug" json:"debug,omitempty"`
	Resources []Resource `yaml:"resources" json:"resources,omitempty"`
}
