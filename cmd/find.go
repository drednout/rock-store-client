package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

func rock_find(cmd *cobra.Command, args []string) {
	params := url.Values{}
	params.Add("q", args[0])
	encodedParams := params.Encode()
    rock_find_url := fmt.Sprintf("%s/v2/rocks/find?%s",
		storeUrl, encodedParams)
    fmt.Printf("SCAFFOLD: rock_find_url is %s\n", rock_find_url)
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
	fmt.Printf("%s", body)
}
