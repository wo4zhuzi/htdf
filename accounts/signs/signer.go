package signs

import (
	"errors"
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
		Msg:           stdTx.GetMsg(),
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	sig := stdTx.GetSignature()

	if sig.IsEmpty() {
		sig = stdSignature
	}
	//fmt.Print(stdTx.GetMsgs()[0].GetSignBytes())
	// fmt.Println("SignTx:stdTx", stdTx)
	signedStdTx = auth.NewStdTx(stdTx.GetMsg(), sig, stdTx.GetMemo())
	return
}

//
func BuildAndSign(txbuilder authtxb.TxBuilder, privKey tmcrypto.PrivKey, msg sdk.Msg) ([]byte, error) {
	signed, err := BuildSignMsg(txbuilder, msg)
	if err != nil {
		return nil, err
	}

	return Sign(txbuilder, privKey, signed)
}

//
func BuildSignMsg(txbuilder authtxb.TxBuilder, msg sdk.Msg) (authtxb.StdSignMsg, error) {
	fmt.Println("BuildSignMsg:txbuilder.Gas()", txbuilder.GasWanted())
	chainID := txbuilder.ChainID()
	if chainID == "" {
		return authtxb.StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}
	// junying-todo, 2019-11-08
	// converted from fee based to gas*gasprice based
	if txbuilder.GasPrices().IsZero() {
		return authtxb.StdSignMsg{}, errors.New("gasprices can't not be zero")
	}
	if txbuilder.GasWanted() <= 0 {
		return authtxb.StdSignMsg{}, errors.New("gas must be greater than zero")
	}
	fmt.Println("BuildSignMsg:Fee", sdk.NewStdFee(txbuilder.GasWanted(), txbuilder.GasPrices().String()), txbuilder.GasWanted())
	return authtxb.StdSignMsg{
		ChainID:       txbuilder.ChainID(),
		AccountNumber: txbuilder.AccountNumber(),
		Sequence:      txbuilder.Sequence(),
		Memo:          txbuilder.Memo(),
		Msg:           msg,
	}, nil
}

//
func Sign(txbuilder authtxb.TxBuilder, privKey tmcrypto.PrivKey, msg authtxb.StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(privKey, msg)
	if err != nil {
		return nil, err
	}

	en := txbuilder.TxEncoder()
	return en(auth.NewStdTx(msg.Msg, sig, msg.Memo))
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
