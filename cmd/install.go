package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		fmt.Println("installing ...")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		config := core.PkgfyConfig{
			// TODO: Make this configurable
			InstallDir: path.Join(homeDir, "bin"),
		}
		if err := core.Install(pkg, &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
