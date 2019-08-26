package cli

import (
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

// GetListAccCmd lists Accounts
func GetListAccCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show local account list",
		Long:  "list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var am *accounts.Manager
			// Assemble the account manager and supported backends
			rootDir := viper.GetString(cli.HomeFlag)
			defaultKeyStoreHome := filepath.Join(rootDir, "keystores")
			backends := []accounts.Backend{
				keystore.NewKeyStore(defaultKeyStoreHome),
			}

			am = accounts.NewManager(backends...)
			if am == nil {
				fmt.Print("Get account list error !")
				return nil
			}
			var index int
			for _, wallet := range am.Wallets() {
				for _, account := range wallet.Accounts() {
					fmt.Printf("%s\n", account.Address)
					index++
				}
			}
			return nil
		},
	}
}
