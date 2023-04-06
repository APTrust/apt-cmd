package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/APTrust/preservation-services/models/registry"
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

func EnsureDefaultListParams(values url.Values) {
	if values.Get("sort") == "" {
		values.Set("sort", "id")
	}
	if values.Get("per_page") == "" {
		values.Set("per_page", "25")
	}
}

func PrettyPrint(resp *network.RegistryResponse, objType network.RegistryObjectType) {
	fmt.Println(objType, ":", resp.Count, "results")
	switch objType {
	case network.RegistryIntellectualObject:
		prettyPrintObjectList(resp.IntellectualObjects())
	case network.RegistryGenericFile:
		prettyPrintFileList(resp.GenericFiles())
	case network.RegistryWorkItem:
		prettyPrintWorkItemList(resp.WorkItems())
	}
	if resp.HasNextPage() {
		fmt.Println("Next Page:", resp.Next)
	}
}

func prettyPrintObjectList(list []*registry.IntellectualObject) {

}

func prettyPrintFileList(list []*registry.GenericFile) {

}

func prettyPrintWorkItemList(list []*registry.WorkItem) {

}

func PrettyPrintObject(obj *registry.IntellectualObject) {

	// // Identifier Info
	// obj.ID
	// obj.Identifier
	// obj.BagName
	// obj.ETag
	// obj.InstitutionID
	// obj.InstitutionIdentifier
	// obj.Title
	// obj.Access
	// obj.State // use English: Active or Deleted
	// obj.StorageOption
	// obj.BagItProfileIdentifier // limit to short string: APTrust or BTR

	// // Stats
	// obj.Size
	// obj.FileCount
	// obj.PayloadFileCount
	// obj.PayloadSize
	// obj.CreatedAt
	// obj.UpdatedAt

	// // Inst internal info
	// obj.AltIdentifier
	// obj.BagGroupIdentifier
	// obj.InternalSenderDescription
	// obj.InternalSenderIdentifier
	// obj.SourceOrganization

	// // Show description last, as it can be lengthy
	// obj.Description

}

func PrettyPrintFile(gf *registry.GenericFile) {

	// File attrs, checksums, events, storage records

}

func PrettyPrintChecksum(cs *registry.Checksum) {
	fmt.Println("    Date:      ", cs.DateTime.Format(time.RFC3339))
	fmt.Println("    Algorithm: ", cs.Algorithm)
	fmt.Println("    Digest:    ", cs.Digest)
}

func PrettyPrintEvent(event *registry.PremisEvent) {
	fmt.Println("    ID:             ", event.ID)
	fmt.Println("    UUID:           ", event.Identifier)
	fmt.Println("    Date:           ", event.DateTime)
	fmt.Println("    Type:           ", event.EventType)
	fmt.Println("    Outcome:        ", event.Outcome)

	fmt.Println("    Agent:          ", event.Agent)
	fmt.Println("    Detail:         ", event.Detail)
	fmt.Println("    Object:         ", event.Object)
	fmt.Println("    Outcome Detail: ", event.OutcomeDetail)
	fmt.Println("    Info:           ", event.OutcomeInformation)

}

func PrettyPrintStorageRecord(sr *registry.StorageRecord) {
	fmt.Println("    ID:  ", sr.ID)
	fmt.Println("    URL: ", sr.URL)
}

func PrettyPrintWorkItem(item *registry.WorkItem) {

	// // General: common to all items
	// item.ID
	// item.Name
	// item.Size
	// item.Action
	// item.Stage
	// item.Status
	// item.Note
	// item.NeedsAdminReview
	// item.UpdatedAt

	// // Obj Info
	// item.ETag
	// item.BagDate
	// item.IntellectualObjectID
	// item.ObjectIdentifier
	// item.StorageOption

	// // File Info
	// item.GenericFileID
	// item.GenericFileIdentifier

	// // Restorations & deletions
	// item.User
	// item.InstApprover

	// fmt.Printf(``)
}
