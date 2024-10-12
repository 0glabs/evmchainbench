package option

import (
	"github.com/spf13/cobra"
)

func OptionsForGeneration(cmd *cobra.Command) {
	cmd.Flags().StringP("faucet-private-key", "f", "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306", "Private key of a faucet account")
	cmd.Flags().IntP("sender-count", "s", 4, "The number of senders of generated transactions")
	cmd.Flags().IntP("tx-count", "t", 100000, "The number of tx count each sender will broadcast")
}

func OptionsForTxStore(cmd *cobra.Command) {
	cmd.Flags().StringP("tx-store-dir", "d", "/tmp/0g-benchmark-dir", "The directory of storing generated transactions")
}
