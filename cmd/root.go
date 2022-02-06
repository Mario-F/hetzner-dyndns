package cmd

import (
	"fmt"
	"os"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile     string
	token       string
	recordInput string
)

var Version = "development"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hetzner-dyndns",
	Short: "Updates Hetzner DNS with external IP.",
	Long:  `A tool for updating Hetzner DNS with your external IP.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&logger.DebugMode, "debug", false, "Turn on debug messages")
	rootCmd.PersistentFlags().BoolVar(&logger.QuietMode, "quiet", false, "Suppress logmessages")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hetzner-dyndns.yaml)")
}

// initLogging set up logging with logrus
func initLogging() {
	if logger.DebugMode {
		log.Info("Debug loglevel enabled")
		log.SetLevel(log.DebugLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".hetzner-dyndns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hetzner-dyndns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
