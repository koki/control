package ctl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/koki/short/client"
	"github.com/koki/short/parser"
)

func KubectlGeneric(kubeArgs []string) error {
	var err error
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func KubectlGenericDashO(kubeArgs []string, outputFormat string, kubeOutput io.Writer) error {
	var err error
	kubeArgs = append(kubeArgs, "-o", outputFormat)
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = kubeOutput
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func KubectlGenericWithFiles(kubeArgs []string, kubeFile string) error {
	var err error

	kubeArgs = append(kubeArgs, "-f", kubeFile)
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func KubectlGenericWithFilesDashO(kubeArgs []string, kubeFile string, outputFormat string, kubeOutput io.Writer) error {
	var err error
	kubeArgs = append(kubeArgs, "-f", kubeFile, "-o", outputFormat)
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = kubeOutput
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func ProcessArgs(args []string) (filenames []string, outputFormat string, kubeArgs []string) {
	filenames = []string{}
	kubeArgs = []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-f":
			if i < len(args)-1 {
				i++
				filenames = append(filenames, args[i])
				continue
			}
		case "-o":
			if i < len(args)-1 {
				o := strings.ToLower(args[i+1])

				if o == "yaml" || o == "json" {
					i++
					outputFormat = o
					continue
				}
			}
		}

		kubeArgs = append(kubeArgs, arg)
	}

	return
}

func ShortGenericCmd(args []string) error {
	var err error
	filenames, outputFormat, kubeArgs := ProcessArgs(args)

	useStdin := false
	if len(filenames) == 1 && filenames[0] == "-" {
		useStdin = true
	}

	kubeOutStream := &bytes.Buffer{}
	if len(filenames) > 0 {
		var streams []io.ReadCloser
		if useStdin {
			streams = []io.ReadCloser{
				os.Stdin,
			}
		} else {
			streams, err = parser.OpenStreamsFromFiles(filenames)
			if err != nil {
				return err
			}
		}

		kubeObjs, err := client.ConvertEitherStreamsToKube(streams)
		if err != nil {
			return nil
		}
		kubeStream := &bytes.Buffer{}
		err = client.WriteObjsToYamlStream(kubeObjs, kubeStream)
		if err != nil {
			return err
		}

		kubeFile, err := ioutil.TempFile("", "tmp-koki-manifest")
		if err != nil {
			return fmt.Errorf("couldn't create temporary file for manifest")
		}
		defer os.Remove(kubeFile.Name())
		if _, err := kubeFile.Write(kubeStream.Bytes()); err != nil {
			return fmt.Errorf("couldn't write manifest to temporary file")
		}
		if err := kubeFile.Close(); err != nil {
			return fmt.Errorf("couldn't close temporary manifest file")
		}

		if len(outputFormat) == 0 {
			return KubectlGenericWithFiles(kubeArgs, kubeFile.Name())
		}

		err = KubectlGenericWithFilesDashO(kubeArgs, kubeFile.Name(), outputFormat, kubeOutStream)
		if err != nil {
			return err
		}
	} else {
		if len(outputFormat) == 0 {
			return KubectlGeneric(kubeArgs)
		}

		err = KubectlGenericDashO(kubeArgs, outputFormat, kubeOutStream)
		if err != nil {
			return err
		}
	}

	// Copy the contents of kubeOutStream in case parsing fails and we need to backtrack.
	kokiObjs, err := client.ConvertKubeStreams([]io.ReadCloser{ioutil.NopCloser(bytes.NewBuffer(kubeOutStream.Bytes()))})
	if err != nil {
		// Couldn't parse output. Just pass it through.
		_, _ = os.Stdout.Write(kubeOutStream.Bytes())
		return nil
	}

	switch outputFormat {
	case "yaml":
		err = client.WriteObjsToYamlStream(kokiObjs, os.Stdout)
	case "json":
		err = client.WriteObjsToJSONStream(kokiObjs, os.Stdout)
	}

	return err

}
