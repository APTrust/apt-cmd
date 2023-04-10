package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/APTrust/preservation-services/network"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/pflag"
)

const (
	// EXIT_OK means program completed successfully.
	EXIT_OK = 0

	// EXIT_RUNTIME_ERR means program did not complete
	// successfully due to an error. The error may have
	// occurred outside the program, such as a network
	// error or an error on a remote server.
	EXIT_RUNTIME_ERR = 1

	// EXIT_BAG_INVALID is used primarily for apt_validate.
	// It means the program completed its run and found that
	// the bag is not valid.
	EXIT_BAG_INVALID = 2

	// EXIT_USER_ERR means the user did not supply some
	// required option or argument, or the user supplied
	// invalid options/arguments.
	EXIT_USER_ERR = 3

	// EXIT_REQUEST_ERROR means the remote server responded
	// with a 4xx or 5xx HTTP status code.
	EXIT_REQUEST_ERROR = 4

	// EXIT_NO_OP means the user requested help message or
	// version info. The program printed the info, and no other
	// operations were performed.
	EXIT_NO_OP = 100
)

var ErrImbalancedArgPair = errors.New("odd number of filter args")

type ArgPair struct {
	Name  string
	Value string
}

func ParseArgPairs(args []string) []ArgPair {
	pairs := make([]ArgPair, 0)
	for _, arg := range args {
		if !strings.Contains(arg, "=") {
			continue
		}
		parts := strings.SplitN(arg, "=", 2)
		pairs = append(pairs, ArgPair{Name: parts[0], Value: parts[1]})
	}
	return pairs
}

func GetUrlValues(args []string) url.Values {
	pairs := ParseArgPairs(args)
	v := url.Values{}
	for _, pair := range pairs {
		v.Add(pair.Name, pair.Value)
	}
	return v
}

func NewRegistryClient(config *Config) (*network.RegistryClient, error) {
	err := config.ValidateRegistryConfig()
	if err != nil {
		return nil, err
	}
	client, err := network.NewRegistryClient(
		config.RegistryURL,
		config.RegistryAPIVersion,
		config.RegistryEmail,
		config.RegistryAPIKey,
		logger,
	)
	if client != nil {
		client.UseMemberAPI()
	}
	return client, err
}

func InitRegistryRequest(args []string) (*network.RegistryClient, url.Values) {
	urlValues := GetUrlValues(args)
	client, err := NewRegistryClient(config)
	if err != nil {
		// Only error here is that user didn't supply valid
		// Registry config or Registry crendentials
		fmt.Fprintln(os.Stderr, "Error getting Registry client:", err)
		os.Exit(EXIT_USER_ERR)
	}
	return client, urlValues
}

func EnsureDefaultListParams(values url.Values) {
	if values.Get("sort") == "" {
		values.Set("sort", "id")
	}
	if values.Get("per_page") == "" {
		values.Set("per_page", "25")
	}
}

func PrettyPrintJSON(jsonBytes []byte) {
	pretty := new(bytes.Buffer)
	err := json.Indent(pretty, jsonBytes, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting JSON:", err)
		os.Exit(EXIT_RUNTIME_ERR)
	}
	fmt.Println(pretty.String())
}

func GetS3Client(s3Host string) *minio.Client {
	err := config.ValidateAWSCredentials()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Missing S3 connection info:", err)
		os.Exit(EXIT_USER_ERR)
	}
	client, err := minio.New(
		s3Host,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.AWSKey, config.AWSSecret, ""),
			Secure: true,
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating S3 client:", err)
		os.Exit(EXIT_RUNTIME_ERR)
	}
	return client
}

func DisallowPreservationBucket(bucket string) {
	b := strings.ToLower(bucket)
	if strings.Contains(b, "aptrust") && strings.Contains(b, "preservation") {
		fmt.Fprintln(os.Stderr, "Illegal bucket:", bucket, "MC Hammer says you can't touch this!")
		os.Exit(EXIT_USER_ERR)
	}
}

func GetParam(flags *pflag.FlagSet, paramName, errMsg string) string {
	paramValue := flags.Lookup(paramName).Value.String()
	if paramValue == "" && errMsg != "" {
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(EXIT_USER_ERR)
	}
	return paramValue
}
