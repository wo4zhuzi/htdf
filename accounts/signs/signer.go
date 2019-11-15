package signs

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	tmcrypto "github.com/orientwalt/tendermint/crypto"
)

func sign(privKey tmcrypto.PrivKey, msg []byte) (sig []byte, pub tmcrypto.PubKey, err error) {
	sig, err = privKey.Sign(msg)
	if err != nil {
		return nil, nil, err
	}
	pub = privKey.PubKey()
	return sig, pub, nil
}

// //SignTxWithLocalFile return sign byte msg and pubkey
// func SignTxWithLocalFile(filename string, auth string, msg []byte) (sig []byte, pub tmcrypto.PubKey, err error) {
// 	priv, err := accounts.GetPrivKey(filename, auth)
// 	sig, pub, err = sign(priv, msg)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return sig, pub, err
// }

//SignTx return stdTx from unsigned to sign
func SignTx(txbuilder authtxb.TxBuilder, stdTx auth.StdTx, privKey tmcrypto.PrivKey) (signedStdTx auth.StdTx, err error) {
	stdSignature, err := MakeSignature(privKey, authtxb.StdSignMsg{
		ChainID:       txbuilder.ChainID(),
		AccountNumber: txbuilder.AccountNumber(),
		Sequence:      txbuilder.Sequence(),
		Fee:           stdTx.Fee,
		Msgs:          stdTx.GetMsgs(),
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	sigs := stdTx.GetSignatures()

	if len(sigs) == 0 {
		sigs = []auth.StdSignature{stdSignature}
	} else {
		sigs = append(sigs, stdSignature)
	}
	//fmt.Print(stdTx.GetMsgs()[0].GetSignBytes())
	// fmt.Println("SignTx:stdTx", stdTx)
	signedStdTx = auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	return
}

//
func BuildAndSign(txbuilder authtxb.TxBuilder, privKey tmcrypto.PrivKey, msgs []sdk.Msg) ([]byte, error) {
	msg, err := BuildSignMsg(txbuilder, msgs)
	if err != nil {
		return nil, err
	}

	return Sign(txbuilder, privKey, msg)
}

//
func BuildSignMsg(txbuilder authtxb.TxBuilder, msgs []sdk.Msg) (authtxb.StdSignMsg, error) {
	fmt.Println("BuildSignMsg:txbuilder.GasWanted()", txbuilder.GasWanted())
	chainID := txbuilder.ChainID()
	if chainID == "" {
		return authtxb.StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}
	// junying-todo, 2019-11-08
	// converted from fee based to gas*gasprice based
	// if txbuilder.GasPrices().IsZero() {
	// 	return authtxb.StdSignMsg{}, errors.New("gasprices can't not be zero")
	// }
	// if txbuilder.GasWanted() <= 0 {
	// 	return authtxb.StdSignMsg{}, errors.New("gasWanted must be provided")
	// }
	fmt.Println("BuildSignMsg:Fee", auth.NewStdFee(txbuilder.GasWanted(), txbuilder.GasPrices()), txbuilder.GasWanted())
	return authtxb.StdSignMsg{
		ChainID:       txbuilder.ChainID(),
		AccountNumber: txbuilder.AccountNumber(),
		Sequence:      txbuilder.Sequence(),
		Memo:          txbuilder.Memo(),
		Msgs:          msgs,
		Fee:           auth.NewStdFee(txbuilder.GasWanted(), txbuilder.GasPrices()), // auth.NewStdFee(txbuilder.GasWanted(), fees),
	}, nil
}

//
func Sign(txbuilder authtxb.TxBuilder, privKey tmcrypto.PrivKey, msg authtxb.StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(privKey, msg)
	if err != nil {
		return nil, err
	}

	en := txbuilder.TxEncoder()
	fmt.Println("4--------------------", msg.Msgs, msg.Fee)
	return en(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
}

// MakeSignature builds a StdSignature given keybase, key name, passphrase, and a StdSignMsg.
func MakeSignature(privKey tmcrypto.PrivKey, msg authtxb.StdSignMsg) (sig auth.StdSignature, err error) {
	sigBytes, pubkey, err := sign(privKey, msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}
