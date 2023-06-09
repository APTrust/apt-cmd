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
	Long: `APTrust partner tools.

    * Create and validate bags.
    * Upload to and download from S3.
    * Report on WorkItems, objects and files in the registry.

	Source: https://github.com/APTrust/apt-cmd
	Docs: https://aptrust.github.io/userguide/partner_tools/

`,
}

var config *Config
var debug bool
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aptrust)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug output to stderr")
}

func initConfig() {
	initLogger()
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
		ConfigSource:       configSource,
	}
	logger.Debug(config.String())
}

func initLogger() {
	outStream := io.Discard
	if debug {
		outStream = os.Stderr
		fmt.Println("DEBUG = ", debug)
	}
	logger = logging.MustGetLogger("aptrust")
	logging.SetLevel(logging.DEBUG, "aptrust")
	logBackend := logging.NewLogBackend(outStream, "[debug] ", stdlog.Lmsgprefix)
	logging.SetBackend(logBackend)
}
