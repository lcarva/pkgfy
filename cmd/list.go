package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/lcarva/pkgfy/internal/persistance"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Long:  "",
	// Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		// TODO: Make DB path configurable
		dbPath := path.Join(homeDir, ".pkgfy.db")

		pkgs, err := persistance.List(dbPath)
		if err != nil {
			return
		}

		for _, pkg := range pkgs {
			if pkg.Alias == "" {
				fmt.Printf("%s\t%s\n", pkg.Name, pkg.URL)
			} else {
				fmt.Printf("%s (%s)\t%s\n", pkg.Name, pkg.Alias, pkg.URL)
			}
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
