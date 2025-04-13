package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var storeUrl string
var skopeoBinaryName string
var todoMessage = color.YellowString("TODO: implement me")

func Execute() {
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Get program version and exit",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: " + version)
		},
	}

	var cmdFind = &cobra.Command{
		Use:     "find [text]",
		Short:   "Find rock package in the Store",
		Args:    cobra.MinimumNArgs(1),
		Example: "rock-store-client find postgresql",
		Run:     rock_find,
	}
	var cmdInfo = &cobra.Command{
		Use:   "info [rock-name]",
		Short: "Get extended info about the rock package from the Store",
		Args:  cobra.MinimumNArgs(1),
		Run:   rock_info,
	}
	var cmdDownload = &cobra.Command{
		Use:     "download [rock-name] [channel]",
		Short:   "Download rock in oci-archive format via skopeo",
		Args:    cobra.MinimumNArgs(2),
		Example: "rock-store-client download postgresql 14/stable",
		Run:     rock_download,
	}

	var rootCmd = &cobra.Command{Use: "rock-store-client"}
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdFind)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdDownload)
	rootCmd.PersistentFlags().StringVar(
		&storeUrl,
		"store-url",
		"https://api.staging.snapcraft.io",
		"Store URL, default api.snapcraft.io",
	)
	rootCmd.PersistentFlags().StringVar(
		&skopeoBinaryName,
		"skopeo-binary-name",
		"rockcraft.skopeo",
		"Skopeo utility name, default rockcraft.skopeo",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
