package main

type BetweenCount struct {
	MinCount int64 `yaml:"minCount"`
	MaxCount int64 `yaml:"maxCount"`
}

type Count struct {
	Exact   int64        `yaml:"exact"`
	Between BetweenCount `yaml:"between"`
}

type RandomField struct {
	Field string `yaml:"field"`
	Type  string `yaml:"type"`
}

type Resource struct {
	TemplatePath string        `yaml:"templatePath"`
	Namespace    string        `yaml:"namespace"`
	Count        Count         `yaml:"count"`
	Randomize    []RandomField `yaml:"randomize"`
}

type Config struct {
	Resources []Resource `yaml:"resources"`
}
