package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Provides a guided setup",
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "TODO",
		}

		_, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
