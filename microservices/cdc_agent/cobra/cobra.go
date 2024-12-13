package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
)

// CusterNodeId Declare global flags
var CusterNodeId int

// Root command
var RootCmd = &cobra.Command{
	Use:   "Idra Cdc",
	Short: "Idra Cdc agent",
	Long:  `Idra is a Change Data Capture platform that uses Golang features to make simple data capture.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show user information",
	Run: func(cmd *cobra.Command, args []string) {
		CusterNodeId, _ := cmd.Flags().GetInt("cluster_node_id")
		fmt.Printf("Started node with id %d", CusterNodeId)
	},
}

func init() {
	// Add flags for global variables
	RootCmd.PersistentFlags().IntVarP(&CusterNodeId, "cluster_node_id", "c", 0, "Set id for this node if part of a cluster")

	// Add subcommands
	RootCmd.AddCommand(infoCmd)
}
