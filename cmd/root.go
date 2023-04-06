package cmd

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"path"

	"github.com/op/go-logging"
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
var logger *logging.Logger

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	initLogger()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aptrust)")
}

func initConfig() {
	useConfigFile := false
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		useConfigFile = true
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		if _, err := os.Stat(path.Join(home, ".aptrust")); !os.IsNotExist(err) {
			viper.AddConfigPath(home)
			viper.SetConfigType("env")
			viper.SetConfigName(".aptrust")
			useConfigFile = true
		}
	}
	viper.AutomaticEnv()

	if useConfigFile {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading config file:", err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
	}

	configSource := "Environment Variables"
	if useConfigFile {
		configSource = viper.ConfigFileUsed()
	}
	config = &Config{
		RegistryEmail:      viper.GetString("APTRUST_REGISTRY_EMAIL"),
		RegistryAPIKey:     viper.GetString("APTRUST_REGISTRY_API_KEY"),
		RegistryURL:        viper.GetString("APTRUST_REGISTRY_URL"),
		RegistryAPIVersion: viper.GetString("APTRUST_REGISTRY_API_VERSION"),
		AWSKey:             viper.GetString("APTRUST_AWS_KEY"),
		AWSSecret:          viper.GetString("APTRUST_AWS_SECRET"),
		ReceivingBucket:    viper.GetString("APTRUST_RECEIVING_BUCKET"),
		RestorationBucket:  viper.GetString("APTRUST_RESTORATION_BUCKET"),
		DownloadDir:        viper.GetString("APTRUST_DOWNLOAD_DIR"),
		ConfigSource:       configSource,
	}
}

func initLogger() {
	// We do this because the registry client requires
	// a logger, but we don't really want to log its output.
	logger = logging.MustGetLogger("aptrust")
	logging.SetLevel(logging.INFO, "aptrust")
	logBackend := logging.NewLogBackend(io.Discard, "", stdlog.LstdFlags)
	logging.SetBackend(logBackend)
}
