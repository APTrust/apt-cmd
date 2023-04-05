package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aptrust",
	Short: "APTrust partner tools.",
	Long:  `APTrust partner tools to validate bags and query the Registry.`,
}

var config *Config
var cfgFile string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aptrust)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("env")
		viper.SetConfigName(".aptrust")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Error reading config file:", err.Error())
		os.Exit(EXIT_RUNTIME_ERR)
	}

	demoConfig := EnvConfig{
		RegistryEmail:     viper.GetString("DEMO_REGISTRY_EMAIL"),
		RegistryAPIKey:    viper.GetString("DEMO_REGISTRY_API_KEY"),
		AWSKey:            viper.GetString("DEMO_AWS_KEY"),
		AWSSecret:         viper.GetString("DEMO_AWS_SECRET"),
		ReceivingBucket:   viper.GetString("DEMO_RECEIVING_BUCKET"),
		RestorationBucket: viper.GetString("DEMO_RESTORATION_BUCKET"),
	}

	prodConfig := EnvConfig{
		RegistryEmail:     viper.GetString("PROD_REGISTRY_EMAIL"),
		RegistryAPIKey:    viper.GetString("PROD_REGISTRY_API_KEY"),
		AWSKey:            viper.GetString("PROD_AWS_KEY"),
		AWSSecret:         viper.GetString("PROD_AWS_SECRET"),
		ReceivingBucket:   viper.GetString("PROD_RECEIVING_BUCKET"),
		RestorationBucket: viper.GetString("PROD_RESTORATION_BUCKET"),
	}

	config = &Config{
		Demo:        demoConfig,
		Prod:        prodConfig,
		DownloadDir: viper.GetString("DOWNLOAD_DIR"),
	}

	fmt.Println("CONFIG:", config)
}
