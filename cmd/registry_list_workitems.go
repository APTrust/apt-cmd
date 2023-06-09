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

Examples:

List recent ingests:

  apt-cmd registry list workitems action='Ingest' sort='date_processed__desc'

List all work items since April 6, 2023:

  apt-cmd registry list workitems date_processed__gteq='2023-04-06' sort='date_processed__desc'

List failed work items:

  apt-cmd registry list workitems status='Failed' sort='date_processed__desc'

List work items pertaining to a tar file you uploaded:

  apt-cmd registry list workitems name='bag-of-photos.tar'

List work items pertaining to a bag with a specific etag:

  apt-cmd registry list workitems etag='987654321-100'

List work items pertaining to a specific intellectual object:

  apt-cmd registry list workitems object_identifier='test.edu/TestBag'

List restorations or deletions of a specific file:

  apt-cmd registry list workitems generic_file_identifier='test.edu/TestBag/data/photo1.jpg'


Quick Reports:

List all items from the past 30 days that are still in process:

  apt-cmd registry list workitems --report=inprocess

List all items from the past 30 days that failed or were cancelled:

  apt-cmd registry list workitems --report=problems

List all restorations from the past 30 days:

  apt-cmd registry list workitems --report=restorations

When running quick reports, this tool ignores all other query params.

Full online documentation:

  https://aptrust.github.io/userguide/partner_tools/

	`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(config, args)
		EnsureDefaultListParams(urlValues)

		report := cmd.Flags().Lookup("report").Value.String()
		if report != "" {
			var err error
			logger.Debugf("Running WorkItem report %s", report)
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
