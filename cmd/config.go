package cmd

type Config struct {
	Demo        EnvConfig
	Prod        EnvConfig
	DownloadDir string
}

type EnvConfig struct {
	RegistryEmail     string
	RegistryAPIKey    string
	AWSKey            string
	AWSSecret         string
	ReceivingBucket   string
	RestorationBucket string
}
