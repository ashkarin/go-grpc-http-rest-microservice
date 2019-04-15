package commands

import "github.com/spf13/cobra"

func Execute(args []string) (*cobra.Command, error) {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdServer)
	return rootCmd, rootCmd.Execute()
}
