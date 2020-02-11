package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long: `Prints the specific version of this iteration of PVGO`,
	Run: func(cmd *cobra.Command, args []string) {
		println(VERSION)
	},
}


func init(){
	pvgoCmd.AddCommand(versionCmd)
}
