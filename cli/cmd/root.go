package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd root cobra command.
var RootCmd = &cobra.Command{
	Use:   "cli <subcommand>",
	Short: "use the koki cli to do koki things",
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "[used for testing]: deploy a koki controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := EnvFromFlags()
		if err != nil {
			return err
		}

		_, err = CreateControllerIfNeeded(env)
		return err
	},
}

var createCmd = &cobra.Command{
	Use:   "create <resource.yaml>",
	Short: "create a plain kubernetes resource from koki shorthand",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := EnvFromFlags()
		if err != nil {
			return err
		}

		return CreateFromShortFile(env, args[0])
	},
}

var createAppCmd = &cobra.Command{
	Use:   "create-app <pod.yaml>",
	Short: "create a koki app from a pod",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := EnvFromFlags()
		if err != nil {
			return err
		}

		pod, err := loadTemplate(args[0], args[1:]...)
		if err != nil {
			return err
		}

		return CreateAppForPod(env, pod)
	},
}

var deleteAppCmd = &cobra.Command{
	Use:   "delete-app <pod.yaml>",
	Short: "delete the koki app for a pod (if it exists)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := EnvFromFlags()
		if err != nil {
			return err
		}

		pod, err := loadTemplate(args[0], args[1:]...)
		if err != nil {
			return err
		}

		return DeleteAppForPodIfExists(env.Namespace, env.Client, pod.Name)
	},
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "purge the koki controller and all koki apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := EnvFromFlags()
		if err != nil {
			return err
		}

		return PurgeAppsAndController(env.Namespace, env.Client)
	},
}

func init() {
	RootCmd.AddCommand(controllerCmd, createAppCmd, deleteAppCmd, purgeCmd, createCmd)
}
