package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/APTrust/dart-runner/bagit"
	"github.com/APTrust/preservation-services/network"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// ParseArgPairs converts command-line args that look like
// key-value pairs into ArgPair objects. It ignores flags
// and args that do not contain an equal sign.
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

// GetUrlValues converts command-line args that look like
// name-value pairs into URL values. It converts only those
// args that contain an equal sign. If an arg is specified
// more than once, the resulting url.Values will include
// all of that arg's values.
func GetUrlValues(args []string) url.Values {
	pairs := ParseArgPairs(args)
	v := url.Values{}
	for _, pair := range pairs {
		v.Add(pair.Name, pair.Value)
	}
	return v
}

// GetTagValues parses tag values from the supplied string slice,
// which should be collected from the --tags flags.
//
// Format is "tagfile.txt/Tag-Name=Value". If tag file name
// is missing from param, it's assumed to be bag-info.txt,
// which is the only customizable tag file in the BagIt standard.
//
// As with the LOC's BagIt-Python library, we convert the first
// letter of each word in tag names to upper-case. For example,
// "source-organization" will be converted here to
// "Source-Organization".
//
// While the BagIt spec does not officially say that tag names must
// use title case, the examples in the BagIt spec and common convention
// use title-cased tag names. Some parses may expect or demand
// title-cased names when validating bags, so we will stick to title
// case for now.
func GetTagValues(args []string) []*bagit.TagDefinition {
	pairs := ParseArgPairs(args)
	titleCase := cases.Title(language.English)
	tagDefs := make([]*bagit.TagDefinition, 0)
	for _, pair := range pairs {
		var tagDef *bagit.TagDefinition
		parts := strings.SplitN(pair.Name, "/", 2)
		if len(parts) == 1 {
			tagDef = &bagit.TagDefinition{
				TagFile:   "bag-info.txt",
				TagName:   titleCase.String(strings.ToLower(pair.Name)),
				UserValue: pair.Value,
			}
		} else {
			tagDef = &bagit.TagDefinition{
				TagFile:   parts[0],
				TagName:   titleCase.String(strings.ToLower(parts[1])),
				UserValue: pair.Value,
			}
		}
		tagDefs = append(tagDefs, tagDef)
	}
	return tagDefs
}

// NewRegistryClient returns a new client that can talk to
// the APTrust Registry. It will return an error if the
// config lacks essential Registry settings.
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

// InitRegistryRequest initializes a registry REST client
// and the params to be sent in a query string.
func InitRegistryRequest(config *Config, args []string) (*network.RegistryClient, url.Values) {
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

// EnsureDefaultListParams adds sort and per_page params to url.Values
// if they're not already there. This ensures that we don't get too many
// results per page, and that paging works correctly.
func EnsureDefaultListParams(values url.Values) {
	if values.Get("sort") == "" {
		values.Set("sort", "id")
	}
	if values.Get("per_page") == "" {
		values.Set("per_page", "25")
	}
}

// PrettyPrintJSON converts unformatted JSON, such as that returned by
// the Registry, to formatted JSON.
func PrettyPrintJSON(jsonBytes []byte) {
	pretty := new(bytes.Buffer)
	err := json.Indent(pretty, jsonBytes, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting JSON:", err)
		fmt.Fprintln(os.Stderr, "Response body:", string(jsonBytes))
		os.Exit(EXIT_RUNTIME_ERR)
	}
	fmt.Println(pretty.String())
}

// NewS3Client returns a client that can talk to an S3 endpoint.
// It will return an error if the config is lacking S3 authentication
// settings.
func NewS3Client(config *Config, s3Host string) *minio.Client {
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

// LooksLikePreservationBucket returns true if the bucket name
// looks like the name of an APTrust preservation bucket.
//
// The caller should exit if the user tries to access an APTrust
// preservation bucket. Though these buckets have strict access
// controls with only a single IAM user allowed to write or delete files,
// we want to protect against the case where some inattentive APTrust
// admin runs this tool with credentials belonging to that one IAM user.
// (E.g. An APTrust admin runs this on a production server, and
// the app pulls crendentials from the environment.)
// We call this only on S3 upload and delete.
func LooksLikePreservationBucket(bucket string) bool {
	looksLikePrez := false
	b := strings.ToLower(bucket)
	if strings.Contains(b, "aptrust.") && strings.Contains(b, ".preservation") {
		looksLikePrez = true
	}
	if strings.Contains(b, "aptrust-") && strings.Contains(b, "-wasabi-") {
		looksLikePrez = true
	}
	return looksLikePrez
}

// GetFlagValue returns the parsed value of flagName. If you specify
// a non-empty error message, this will cause the program to exit with
// code EXIT_USER_ERROR and will print the message.
//
// Call this with an empty error message for optional flags, and it
// won't exit if the flag is missing.
func GetFlagValue(flags *pflag.FlagSet, flagName, errMsg string) string {
	paramValue := flags.Lookup(flagName).Value.String()
	if paramValue == "" && errMsg != "" {
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(EXIT_USER_ERR)
	}
	return paramValue
}

// LoadProfile loads a BagIt profile.
func LoadProfile(name string) (*bagit.Profile, error) {
	profile := &bagit.Profile{}
	var data []byte
	var err error
	switch name {
	case "aptrust":
		data, err = profiles.ReadFile("profiles/aptrust-v2.2.json")
	case "btr":
		data, err = profiles.ReadFile("profiles/btr-v1.0.json")
	case "empty":
		data, err = profiles.ReadFile("profiles/empty_profile.json")
	default:
		err = fmt.Errorf("missing or invalid profile. Only 'aptrust', 'btr' and 'empty' are supported")
	}
	if err == nil && len(data) > 1 {
		err = json.Unmarshal(data, profile)
	}
	return profile, err
}

func PrintErrors(errors []string) {
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
}
