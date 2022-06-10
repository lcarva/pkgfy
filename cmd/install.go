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
	Use:   "install <url>",
	Short: "Install a file from a URL",
	Long:  "The file is downloaded into ~/bin, ready to be executed.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		url := args[0]
		fmt.Println("installing ...")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		config := core.PkgfyConfig{
			// TODO: Make this configurable
			InstallDir: path.Join(homeDir, "bin"),
		}
		client := core.HTTPClient(&http.Client{})
		p := core.Pkgfy{Config: config, Client: &client}
		err = p.Install(url)
		return
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
