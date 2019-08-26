package cli

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

// GetListAccCmd lists Accounts
func ExportPrivateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export all private key list",
		Long:  "export private key from .hscli/keystores",
		Args:  cobra.ExactArgs(1),
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
					//fmt.Printf("%s\n", account.Address)
					privkey, err := keystore.GetPrivKey(account, args[0], rootDir)
					if err != nil {
						fmt.Printf("GetPrivKey error|err=%s\n", err)
						return err
					}
					strPrivKey := hex.EncodeToString(privkey.Bytes())
					priv := strPrivKey[10:]
					fmt.Print("htdf	", account.Address, " ", priv, "\n")
					index++
				}
			}
			return nil
		},
	}
}
