package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/lcarva/pkgfy/internal/core"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [url to fetch] [package name]",
	Short: "Install a file from a URL",
	Long:  "The file is downloaded into ~/bin, ready to be executed.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		url := args[0]
		name := args[1]
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		config := core.PkgfyConfig{
			// TODO: Make this configurable
			InstallDir: path.Join(homeDir, "bin"),
		}
		client := core.HTTPClient(&http.Client{})

		alias, err := cmd.Flags().GetString("alias")
		if err != nil {
			return
		}
		pkg := core.Package{Alias: alias, Name: name, URL: url}
		p := core.Pkgfy{Config: config, Client: &client}

		fmt.Printf("Installing %s from %s...\n", name, url)
		err = p.Install(pkg)
		return
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("alias", "a", "", "Alias name for the local file")
}
