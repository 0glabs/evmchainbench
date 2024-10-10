/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().StringP("rpc-url", "r", "http://127.0.0.1:8545", "RPC url of the chain")
	runCmd.Flags().StringP("faucet-private-key", "f", "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306", "Private key of a faucet account")
	runCmd.Flags().IntP("sender-count", "s", 4, "The number of senders of generated transactions")
	runCmd.Flags().IntP("tx-count", "t", 100000, "The number of tx count each sender will broadcast")
}
