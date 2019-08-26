package cli

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

//recover a account from a private key
func GetRecoverAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recover [private-key] [password]",
		Short: "recover a account from a private key.",
		Long:  "recover a account from a private key. when success,you can find the recovered account keyfile in keystore",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			//var acc accounts.Account
			//var encryptPassword string
			var err error
			strPrivateKey := args[0]
			strPasswd := args[1]

			rootDir := viper.GetString(cli.HomeFlag)
			defaultKeyStoreHome := filepath.Join(rootDir, "keystores")

			privateKey, err := hex.DecodeString(strPrivateKey)
			if err != nil {
				fmt.Printf("decodeString error|err=%s\n", err)
				return err
			}

			ks := keystore.NewKeyStore(defaultKeyStoreHome)
			acc, err := ks.RecoverAccount(privateKey, strPasswd)
			if err != nil {
				fmt.Printf("RecoverAccount error|err=%s\n", err)
				return err
			}

			privkey, err := keystore.GetPrivKey(acc, strPasswd, rootDir)
			if err != nil {
				fmt.Printf("GetPrivKey error|err=%s\n", err)
				return err
			}

			strPubkey := hex.EncodeToString(privkey.PubKey().Bytes()[5:])
			fmt.Printf("strPubkey=%v\n", strPubkey)

			fmt.Printf("RecoverAccount success|address=%s\n", acc.Address)

			return err
		},
	}
}
