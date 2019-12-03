package cli

import (
	_ "path/filepath"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	_ "github.com/tendermint/tmlibs/cli"
)

func GetDelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "del",
		Short: "delete account",
		Long:  "delete account from keystores",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome)
			address := args[0]
			err := ksw.Drop(address)
			if err != nil {
				return err
			}

			println("Delete success: ", address)
			return nil
		},
	}
}
