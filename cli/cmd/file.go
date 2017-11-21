package cmd

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"

	templates "github.com/koki/templates/yaml"
)

func loadTemplate(templatePath string, paramsPaths ...string) (*v1.Pod, error) {
	templateYaml, err := ioutil.ReadFile(templatePath)
	if err != nil {
		glog.Fatalf("couldn't load template file (%s): %v",
			templatePath, err)
	}

	for _, paramsPath := range paramsPaths {
		paramsYaml, err := ioutil.ReadFile(paramsPath)
		if err != nil {
			glog.Fatalf("couldn't load params file (%s): %v",
				paramsPath, err)
		}

		templateYaml, err = templates.Fill(templateYaml, paramsYaml)
		if err != nil {
			glog.Fatalf("couldn't fill template with params\n(%s)\n(%s)\n(%v)",
				templateYaml, paramsYaml, err)
		}
	}

	pod := v1.Pod{}

	err = yaml.Unmarshal(templateYaml, &pod)
	if err != nil {
		return nil, err
	}

	return &pod, nil
}
