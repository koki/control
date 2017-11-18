package ctl

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/koki/short/client"
	"github.com/koki/short/parser"
	shortutil "github.com/koki/short/util"
)

func KubectlGeneric(subCommand string, args []string) error {
	var err error
	kubeArgs := []string{subCommand}
	kubeArgs = append(kubeArgs, args...)
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

func KubectlGenericDashO(subCommand string, args []string, outputFormat string, kubeOutput io.Writer) error {
	var err error
	kubeArgs := []string{subCommand}
	kubeArgs = append(kubeArgs, args...)
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

func KubectlGenericWithFiles(subCommand string, args []string, kubeFiles io.Reader) error {
	var err error
	kubeArgs := []string{subCommand}
	kubeArgs = append(kubeArgs, args...)
	kubeArgs = append(kubeArgs, "-f", "-")
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = kubeFiles
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func KubectlGenericWithFilesDashO(subCommand string, args []string, kubeFiles io.Reader, outputFormat string, kubeOutput io.Writer) error {
	var err error
	kubeArgs := []string{subCommand}
	kubeArgs = append(kubeArgs, args...)
	kubeArgs = append(kubeArgs, "-f", "-", "-o", outputFormat)
	cmd := exec.Command("kubectl", kubeArgs...)
	cmd.Stdin = kubeFiles
	cmd.Stdout = kubeOutput
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func ShortGenericCmd(subCmd string, acceptsFiles bool) *cobra.Command {
	var (
		// filenames holds the input files that are to be converted to shorthand or kuberenetes native syntax
		filenames []string
		// output denotes the format of the converted data
		outputFormat string
	)

	genericCmd := &cobra.Command{
		Use:   subCmd,
		Short: fmt.Sprintf("equivalent to 'kubectl %s'", subCmd),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			useStdin := false
			// TODO: support stdin again?
			/*
				if len(args) == 1 && args[0] == "-" {
					glog.V(3).Info("using stdin for input data")
					useStdin = true
				}
			*/

			outputFormat = strings.ToLower(outputFormat)
			if len(outputFormat) > 0 && outputFormat != "yaml" && outputFormat != "json" {
				return shortutil.UsageErrorf("unexpected value %s for -o --output", outputFormat)
			}

			kubeOutStream := &bytes.Buffer{}
			if acceptsFiles {
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

				if len(outputFormat) == 0 {
					return KubectlGenericWithFiles(subCmd, args, kubeStream)
				}

				err = KubectlGenericWithFilesDashO(subCmd, args, kubeStream, outputFormat, kubeOutStream)
				if err != nil {
					return err
				}
			} else {
				if len(outputFormat) == 0 {
					return KubectlGeneric(subCmd, args)
				}

				err = KubectlGenericDashO(subCmd, args, outputFormat, kubeOutStream)
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
		},
	}

	genericCmd.Flags().StringSliceVarP(&filenames, "filenames", "f", nil, "path or url to input files to read manifests")
	genericCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "output format (yaml|json)")

	// parse the go default flagset to get flags for glog and other packages in future
	genericCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	_ = flag.Set("logtostderr", "true")

	//suppress the incorrect prefix in glog output
	_ = flag.CommandLine.Parse([]string{})

	return genericCmd
}
