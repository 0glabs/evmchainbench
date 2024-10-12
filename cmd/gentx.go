package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/0glabs/evmchainbench/cmd/option"
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
	option.OptionsForGeneration(gentxCmd)
	option.OptionsForTxStore(gentxCmd)
}
