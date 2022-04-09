package cmd

import (
	"fmt"
	"log"

	"github.com/Mario-F/hetzner-dyndns/internal/setup"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Provides a guided setup",
	Run: func(cmd *cobra.Command, args []string) {

		prompt := promptui.Select{
			Label: "Select a setup option",
			Items: []string{"Cronjob"},
		}

		_, option, err := prompt.Run()
		if err != nil {
			fmt.Printf("Select failed %v\n", err)
			return
		}

		switch option {
		case "Cronjob":
			err := setup.Cronjob()
			if err != nil {
				log.Fatalf("Cronjob setup failed %v\n", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
