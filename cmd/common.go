package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/APTrust/preservation-services/network"
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
