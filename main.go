package main

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
)

func main() {
	var cfg Config
	panicIfError(readYAMLConfig(&cfg))

	// TODO: it needs a custom decoder for Resources field of Config.
	err := envconfig.Process("BULKRD", &cfg)
	if err != nil {
		panic(err)
	}

	for _, resource := range cfg.Resources {
		panicIfError(
			kubectl("apply", "-f", resource.TemplatePath, "-n", resource.Namespace, "--dry-run=client"),
		)
	}
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to run %v, output: \n%v", cmd.String(), string(output))
		return err
	}

	log.Println(string(output))

	return nil
}

func readYAMLConfig(cfg *Config) error {
	f, err := os.Open("config.yaml")
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
