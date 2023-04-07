package cmd

import "fmt"

type Config struct {
	RegistryURL        string
	RegistryAPIVersion string
	RegistryEmail      string
	RegistryAPIKey     string
	AWSKey             string
	AWSSecret          string
	S3Host             string
	UploadBucket       string
	DownloadBucket     string
	LocalDownloadDir   string
	ConfigSource       string
}

func (config *Config) ValidateRegistryConfig() error {
	errMsg := ""
	if config.RegistryURL == "" {
		errMsg = "Registry URL is missing from config. "
	}
	if config.RegistryAPIVersion == "" {
		errMsg += "Registry API version is missing from config. "
	}
	if config.RegistryEmail == "" {
		errMsg += "Registry email is missing from config. "
	}
	if config.RegistryAPIKey == "" {
		errMsg += "Registry API key is missing from config. "
	}
	if errMsg != "" {
		errMsg += fmt.Sprintf("(Config source: %s)", config.ConfigSource)
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (config *Config) ValidateAWSCredentials() error {
	errMsg := ""
	if config.AWSKey == "" {
		errMsg = "AWS Key is missing from config. "
	}
	if config.AWSSecret == "" {
		errMsg += "AWS Secret is missing from config. "
	}
	if config.S3Host == "" {
		errMsg += "S3 host is missing from config. "
	}
	if errMsg != "" {
		errMsg += fmt.Sprintf("(Config source: %s)", config.ConfigSource)
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (config *Config) String() string {
	regAPIKey := "[redacted]"
	if config.RegistryAPIKey == "" {
		regAPIKey = "MISSING!"
	}
	awsKey := "[redacted]"
	if config.AWSKey == "" {
		awsKey = "MISSING!"
	}
	awsSecret := "[redacted]"
	if config.AWSSecret == "" {
		awsSecret = "MISSING!"
	}
	return fmt.Sprintf(`Configuration:
	RegistryURL:             %s
	RegistryAPIVersion:      %s
	RegistryEmail:           %s
	RegistryAPIKey:          %s
	AWSKey:                  %s
	AWSSecret:               %s
	S3Host:                  %s
	UploadBucket:            %s
	DownloadBucket:          %s
	LocalDownloadDir:        %s
	ConfigSource:            %s`,
		config.RegistryURL,
		config.RegistryAPIVersion,
		config.RegistryEmail,
		regAPIKey,
		awsKey,
		awsSecret,
		config.S3Host,
		config.UploadBucket,
		config.DownloadBucket,
		config.LocalDownloadDir,
		config.ConfigSource)
}
