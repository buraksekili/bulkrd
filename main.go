package main

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
)

func main() {
	var cfg Config
	setConfiguration(&cfg)

	// TODO: it needs a custom decoder for Resources field of Config.
	err := envconfig.Process("BULKRD", &cfg)
	if err != nil {
		panic(err)
	}

	for _, resource := range cfg.Resources {
		panicIfError(validResource(&resource))

		runMultipleKubectl(
			getCount(resource.Count),
			func() error {
				return kubectl("apply", "-f", resource.TemplatePath, "-n", resource.Namespace, "--dry-run=client")
			},
		)
	}
}

func setConfiguration(cfg *Config) {
	panicIfError(envconfig.Process("BULKRD", cfg))
	panicIfError(readYAMLConfig(cfg.ConfigPath, cfg))
	panicIfError(envconfig.Process("BULKRD", cfg))
}

var (
	ErrInvalidExactCount = errors.New("failed to validate Resource Count: count cannot be negative")
	ErrInvalidMinCount   = errors.New("failed to validate Resource Count: minCount cannot be negative")
	ErrInvalidMaxCount   = errors.New("failed to validate Resource Count: maxCount cannot be negative")
	ErrInvalidBetween    = errors.New("failed to validate Resource Count: maxCount cannot be smaller than minCount")
)

func validResource(resource *Resource) error {
	if resource.Count.Exact < 0 {
		return ErrInvalidExactCount
	}

	if resource.Count.Exact > 0 &&
		(resource.Count.Between.MinCount > 0 || resource.Count.Between.MaxCount > 0) {
		log.Println("Both Count type are specified, the count specified in Exact field will be used.")
	}

	if resource.Count.Between.MinCount < 0 {
		return ErrInvalidMinCount
	}

	if resource.Count.Between.MaxCount < 0 {
		return ErrInvalidMaxCount
	}

	if resource.Count.Between.MaxCount < resource.Count.Between.MinCount {
		return ErrInvalidBetween
	}

	return nil
}

func getCount(count Count) int {
	if count.Exact > 0 {
		return count.Exact
	}

	return rand.Intn(count.Between.MaxCount-count.Between.MinCount) + count.Between.MinCount
}

func runMultipleKubectl(n int, kf func() error) {
	log.Printf("running %v times\n", n)
	for i := 0; i < n; i++ {
		panicIfError(kf())
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

func readYAMLConfig(configPath string, cfg *Config) error {
	f, err := os.Open(configPath)
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
