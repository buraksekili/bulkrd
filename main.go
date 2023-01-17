package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/rand"
	"log"
	"os/exec"
)

func main() {
	cb, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("failed to read file, err: %v", err)
	}

	var c Config

	err = yaml.Unmarshal(cb, &c)
	if err != nil {
		log.Fatalf("failed to unmarshal config to struct, err: %v", err)
	}

	//var defaultConfigFlags = genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	//matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(defaultConfigFlags)
	//
	//f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	//namespace, enforceNamespace, err := f.ToRawKubeConfigLoader().Namespace()
	//if err != nil {
	//	panic(err)
	//}
	//
	//builder := f.NewBuilder()
	//b := builder.
	//	Unstructured().
	//	LocalParam(true).
	//	ContinueOnError().
	//	NamespaceParam(namespace).DefaultNamespace().
	//	FilenameParam(enforceNamespace, &resource.FilenameOptions{Filenames: []string{"./sample.yaml"}}).
	//	Flatten()
	//
	//r := b.Do()
	//if err := r.Err(); err != nil {
	//	panic(err)
	//}

	for _, resource := range c.Resources {
		if resource.Count.Exact > 0 && (resource.Count.Between.MinCount > 0 || resource.Count.Between.MaxCount > 0) {
			panic("failed to validate count")
		}

		if resource.Namespace == "" {
			resource.Namespace = "default"
		}

		rand.IntnRange(int(resource.Count.Between.MinCount), int(resource.Count.Between.MaxCount))

		err = kubectl("apply", "-f", resource.TemplatePath, "-n", resource.Namespace)
		if err != nil {
			panic(err)
		}
	}

}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	return cmd.Run()
}
