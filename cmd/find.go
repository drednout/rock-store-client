package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type RockFindResponse struct {
	Results []RockFindResult `json:"results"`
}
type RockFindResult struct {
	Name             string   `json:"name"`
	PackageId        string   `json:"package_id"`
	RockFindMetadata Metadata `json:"metadata"`
}
type RockFindMetadata struct {
	Description string    `json:"description"`
	Publisher   Publisher `json:"publisher"`
	Summary     string    `json:"summary"`
}

func rock_find(cmd *cobra.Command, args []string) {
	fields := []string{
		"summary",
		"description",
		"publisher",
	}
	params := url.Values{}
	params.Add("q", args[0])
	params.Add("fields", strings.Join(fields, ","))
	encodedParams := params.Encode()
	rock_find_url := fmt.Sprintf("%s/v2/rocks/find?%s",
		storeUrl, encodedParams)
	res, err := http.Get(rock_find_url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n",
			res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	var resp RockFindResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Name\tPublisher\tSummary\t")
	sep := strings.Repeat("-", 25)
	fmt.Fprintf(w, "%s\t%s\t%s\t\n", sep, sep, sep)

	for _, item := range resp.Results {
		fmt.Fprintf(w, "%s\t%s (%s)\t%s\t\n",
			item.Name,
			item.RockFindMetadata.Publisher.DisplayName,
			item.RockFindMetadata.Publisher.Username,
			item.RockFindMetadata.Summary,
		)
	}
	w.Flush()
}
