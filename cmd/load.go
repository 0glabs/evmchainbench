package cmd

import (
	"log"

	"github.com/0glabs/evmchainbench/cmd/option"
	"github.com/0glabs/evmchainbench/lib/load"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load previously generated transactions and run the benchmark",
	Long:  "Load previously generated transactions and run the benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		httpRpc, _ := cmd.Flags().GetString("http-rpc")
		txStoreDir, _ := cmd.Flags().GetString("tx-store-dir")
		loader := load.NewLoader(httpRpc, txStoreDir)
		err := loader.LoadAndRun()
		if err != nil {
			log.Fatalf("Failed to load and run: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	option.OptionsForTxStore(loadCmd)
}
