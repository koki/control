package cmd

import (
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	shortutil "github.com/koki/short/util"
)

func CreateFromShortFile(env *Env, filename string) error {
	data, err := parser.Parse([]string{filename}, false)
	if err != nil {
		return err
	}

	kubeObjs, err := converter.ConvertToKubeNative(data)
	if err != nil {
		return err
	}

	if kubeObjs, ok := kubeObjs.([]interface{}); ok {
		for _, kubeObj := range kubeObjs {
			err := CreateKubeObj(kubeObj, env)
			if err != nil {
				return err
			}
		}
	} else {
		return shortutil.PrettyTypeError(kubeObjs, "expected an array")
	}

	return nil
}
