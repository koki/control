package cmd

import (
	"github.com/koki/control/cli/cmd/ctl"
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
	// TODO: create/delete/get work. replace/apply seem to work
	//       attach doesn't work. exec/run/stop ???
	RootCmd.AddCommand(controllerCmd, createAppCmd, deleteAppCmd, purgeCmd,
		ctl.ShortGenericCmd("create", true), ctl.ShortGenericCmd("delete", false),
		ctl.ShortGenericCmd("get", false), ctl.ShortGenericCmd("replace", true),
		ctl.ShortGenericCmd("apply", true), ctl.ShortGenericCmd("attach", false),
		ctl.ShortGenericCmd("exec", false), ctl.ShortGenericCmd("run", false),
		ctl.ShortGenericCmd("stop", false))
}
