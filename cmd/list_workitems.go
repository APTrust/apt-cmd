package cmd

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// workitemsCmd represents the workitems command
var workitemsCmd = &cobra.Command{
	Use:   "workitems",
	Short: "List work item records from the APTrust registry.",
	Long: `List work items from the APTrust registry, or run a report.

--------------
Basic Examples
--------------

List recent ingests:

	aptrust list workitems action='Ingest' sort='date_processed__desc'

List all work items since April 6, 2023:

	aptrust list workitems date_processed__gteq='2023-04-06' sort='date_processed__desc'

List failed work items:

	aptrust list workitems status='Failed' sort='date_processed__desc'

List work items pertaining to a tar file you uploaded:

	aptrust list workitems name='bag-of-photos.tar'

List work items pertaining to a bag with a specific etag:

	aptrust list workitems etag='987654321-100'

List work items pertaining to a specific intellectual object:

	aptrust list workitems object_identifier='test.edu/TestBag'

List restorations or deletions of a specific file:

	aptrust list workitems generic_file_identifier='test.edu/TestBag/data/photo1.jpg'

-------------
Quick Reports
-------------

List all items from the past 30 days that are still in process:

	aptrust list workitems --report=inprocess

List all items from the past 30 days that failed or were cancelled:

	aptrust list workitems --report=problems

List all restorations from the past 30 days:

	aptrust list workitems --report=restorations

When you run a quick report, this tool ignores your other query params.

	`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		EnsureDefaultListParams(urlValues)

		report := cmd.Flags().Lookup("report").Value.String()
		if report != "" {
			var err error
			logger.Infof("Running WorkItem report %s", report)
			urlValues, err = valuesForWorkItemReport(report)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(EXIT_USER_ERR)
			}
		}

		resp := client.WorkItemList(urlValues)
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
	},
}

func init() {
	listCmd.AddCommand(workitemsCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	workitemsCmd.Flags().StringP("report", "r", "", "Run report: inprocess, problems, restorations")
}

func valuesForWorkItemReport(report string) (url.Values, error) {
	thirtyDaysAgo := time.Now().UTC().Add(-30 * 24 * time.Hour).Format(time.RFC3339)
	values := url.Values{}
	values.Add("date_processed__gteq", thirtyDaysAgo)
	values.Add("sort", "date_processed__desc")
	values.Add("per_page", "25")
	switch report {
	case "inprocess":
		values.Add("status__in", "Started")
		values.Add("status__in", "Pending")
	case "problems":
		values.Add("status__in", "Cancelled")
		values.Add("status__in", "Failed")
		values.Add("status__in", "Suspended")
	case "restorations":
		values.Add("action__in", "Glacier Restore")
		values.Add("action__in", "Restore File")
		values.Add("action__in", "Restore Object")
	default:
		return nil, fmt.Errorf("invalid report '%s' - try 'inprocess', 'problems', or 'restorations'", report)
	}
	return values, nil
}
