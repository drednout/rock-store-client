package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type RockInfoResponse struct {
	Name       string           `json:"name"`
	PackageId  string           `json:"package_id"`
	Metadata   Metadata         `json:"metadata"`
	ChannelMap []ChannelMapItem `json:"channel-map"`
}

type Metadata struct {
	Contact     string    `json:"contact"`
	Description string    `json:"description"`
	License     string    `json:"license"`
	Publisher   Publisher `json:"publisher"`
	Summary     string    `json:"summary"`
}

type Publisher struct {
	DisplayName string `json:"display-name"`
	ID          string `json:"id"`
	Username    string `json:"username"`
	Validation  string `json:"validation"`
}

type ChannelMapItem struct {
	Channel  Channel  `json:"channel"`
	Revision Revision `json:"revision"`
}

type Channel struct {
	Name       string   `json:"name"`
	Platform   Platform `json:"platform"`
	ReleasedAt string   `json:"released-at"`
	Risk       string   `json:"risk"`
	Track      string   `json:"track"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	Channel      string `json:"channel"`
	Name         string `json:"name"`
}

type Revision struct {
	CreatedAt string     `json:"created-at"`
	Download  Download   `json:"download"`
	Platforms []Platform `json:"platforms"`
	Revision  int        `json:"revision"`
	Version   string     `json:"version"`
}

type Download struct {
	SHA256 string `json:"sha-256"`
	URL    string `json:"url"`
}

func rock_info(cmd *cobra.Command, args []string) {
	params := url.Values{}
	fields := []string{
		"description",
		"summary",
		"publisher",
		"contact",
		"license",
		"channel-map",
		"revision",
	}
	params.Add("fields", strings.Join(fields, ","))
	encodedParams := params.Encode()
	package_name := args[0]
	rock_info_url := fmt.Sprintf("%s/v2/rocks/info/%s?%s",
		storeUrl, package_name, encodedParams)
	res, err := http.Get(rock_info_url)
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
	var resp RockInfoResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("name: %s\n", resp.Name)
	fmt.Printf("summary: %s\n", resp.Metadata.Summary)
	fmt.Printf("publisher: %s\n", resp.Metadata.Publisher.Username)
	fmt.Printf("contact: %s\n", resp.Metadata.Contact)
	fmt.Printf("license: %s\n", resp.Metadata.License)
	fmt.Printf("description: %s\n", resp.Metadata.Description)
	fmt.Printf("channels:\n")

	for _, item := range resp.ChannelMap {
		RevisionCreatedAt, err := time.Parse(time.RFC3339Nano, item.Revision.CreatedAt)
		if err != nil {
			fmt.Println("Unable to unparse revision creation date:", item.Revision.CreatedAt)
			return
		}
		fmt.Printf("    %s: %s %s (%d)\n", item.Channel.Name, item.Revision.Version, RevisionCreatedAt.Format("2006-01-02"), item.Revision.Revision)
	}
}
