/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/0glabs/evmchainbench/cmd/option"
	"github.com/0glabs/evmchainbench/lib/run"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "To run the benchmark",
	Long:  "To run the benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpc-url")
		faucetPrivateKey, _ := cmd.Flags().GetString("faucet-private-key")
		senderCount, _ := cmd.Flags().GetInt("sender-count")
		txCount, _ := cmd.Flags().GetInt("tx-count")
		run.Run(rpcUrl, faucetPrivateKey, senderCount, txCount)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	option.OptionsForGeneration(runCmd)
}
