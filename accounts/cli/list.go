package cli

import (
	"fmt"
	_ "path/filepath"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	_ "github.com/tendermint/tmlibs/cli"
)

func GetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "show all account",
		Long:  "show all account in keystore",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

			accounts, err := ksw.Accounts()

			if err != nil {
				return err
			}

			for _, account := range accounts {
				fmt.Printf("%s\n", account.Address)
			}
			return nil
		},
	}
}
