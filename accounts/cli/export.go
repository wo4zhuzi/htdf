package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/orientwalt/htdf/crypto/keys/mintkey"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
)

func GetExportPivKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export all private key list",
		Long:  "export private key from .hscli/keystores",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome)

			accounts, err := ksw.Accounts()
			if err != nil {
				return err
			}

			for _, account := range accounts {
				priv, err := getPrivateKey(ksw, account.Address, args[0])
				if err != nil {
					return err
				}
				fmt.Printf("%s	%s\n", account.Address, priv)
			}
			return nil
		},
	}
}

func getPrivateKey(ksw *keystore.KeyStoreWallet, addr string, passphrase string) (string, error) {
	privKeyArmor, err := ksw.FindPrivKey(addr)
	if err != nil {
		return "", err
	}
	privKey, err := mintkey.UnarmorDecryptPrivKey(privKeyArmor, passphrase)
	if err != nil {
		return "", err
	}
	strPrivKey := hex.EncodeToString(privKey.Bytes())
	priv := strPrivKey[10:]
	return priv, err
}
