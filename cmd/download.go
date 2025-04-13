package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type RockInfoDownloadResponse struct {
	Name       string           `json:"name"`
	PackageId  string           `json:"package_id"`
	ChannelMap []ChannelMapItem `json:"channel-map"`
}

func download_oci_archive(packageName string, channelName string, downloadUrl string) {
	rock_registry_url := fmt.Sprintf("docker://%s", downloadUrl)
	target_path := fmt.Sprintf("%s_%s.rock", packageName, strings.ReplaceAll(channelName, "/", "-"))
	cmd := exec.Command(
		skopeoBinaryName, "copy",
		rock_registry_url,
		fmt.Sprintf("oci-archive:%s", target_path),
	)

	// Attach stdout and stderr to see command output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	fmt.Printf("Running %s...\n", skopeoBinaryName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Command execution failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nDone. Please inspect %s\n", target_path)
}

func rock_download(cmd *cobra.Command, args []string) {
	packageName := args[0]
	channelName := args[1]

	params := url.Values{}
	fields := []string{
		"channel-map",
		"revision",
	}
	params.Add("fields", strings.Join(fields, ","))
	encoded_params := params.Encode()
	rock_info_url := fmt.Sprintf("%s/v2/rocks/info/%s?%s",
		storeUrl, packageName, encoded_params)
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
	var resp RockInfoDownloadResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	for _, item := range resp.ChannelMap {
		if item.Channel.Name == channelName {
			downloadUrl := item.Revision.Download.URL
			download_oci_archive(packageName, channelName, downloadUrl)
		}
	}
}
