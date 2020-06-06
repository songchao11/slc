package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewInitializeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "init all commands from remote.",
		Run: func(cmd *cobra.Command, args []string) {
			initCmd()
		},
	}
}

func initCmd() {
	l := len(commands)
	for i, v := range commands {
		schedule := fmt.Sprintf("[busy working]: download command:%s, schedule is (%d | %d)", v, i+1, l)
		downloadCmdFile(v)
		fmt.Println(schedule)
	}
}
