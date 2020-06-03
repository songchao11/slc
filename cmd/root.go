package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slc",
	Short: "show linux command cheat sheet",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(
		NewShowCommand(),
	)
}
