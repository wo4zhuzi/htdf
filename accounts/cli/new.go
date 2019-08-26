package cli

import (
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/orientwalt/htdf/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

//
const FlagPublicKey string = "pubkey"

// GetNewAccountCmd creates A new account
func GetNewAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "create a new account.",
		Long:  "new_Account: If OFF to use standard client tool ,else use common stdin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//var acc accounts.Account
			var encryptPassword string
			var err error
			if args[0] == "OFF" {

				buf := client.BufferStdin()
				//interactive := viper.GetBool(flagInteractive)
				encryptPassword = "12345678"
				if viper.GetString(FlagPublicKey) == "" && !viper.GetBool(client.FlagUseLedger) {
					encryptPassword, err = client.GetCheckPassword(
						"Enter a passphrase to encrypt your key to disk:",
						"Repeat the passphrase:", buf)
					if err != nil {
						return err
					}
				}

			} else {
				encryptPassword = args[0]
			}
			rootDir := viper.GetString(cli.HomeFlag)
			defaultKeyStoreHome := filepath.Join(rootDir, "keystores")
			address, _, err := keystore.StoreKey(defaultKeyStoreHome, encryptPassword)
			if err != nil {
				return err
			}
			fmt.Print("", address, "\n")

			return err
		},
	}
}
