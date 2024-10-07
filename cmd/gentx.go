package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gentxCmd = &cobra.Command{
	Use:   "gentx",
	Short: "To generate transactions and store them onto disk",
	Long: "To generate transactions and store them onto disk",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gentx called")
	},
}

func init() {
	rootCmd.AddCommand(gentxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gentxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gentxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
