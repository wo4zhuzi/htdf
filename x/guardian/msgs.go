package guardian

import (
	"github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
)

const MsgType = "guardian"

var _, _, _, _ sdk.Msg = MsgAddProfiler{}, MsgAddTrustee{}, MsgDeleteProfiler{}, MsgDeleteTrustee{}

//______________________________________________________________________
// MsgAddProfiler - struct for add a profiler
type MsgAddProfiler struct {
	AddGuardian
}

func NewMsgAddProfiler(description string, address, addedBy sdk.AccAddress) MsgAddProfiler {
	return MsgAddProfiler{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedBy,
			Fee:         sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
			// GasWanted:        uint64(10000),
			// GasPrice:         config.DefaultMinGasPrices,
		},
	}
}

//
func (msg MsgAddProfiler) Route() string { return MsgType }

//
func (msg MsgAddProfiler) Type() string { return "guardian add-profiler" }

//
func (msg MsgAddProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//
func (msg MsgAddProfiler) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

//
func (msg MsgAddProfiler) GetSigner() sdk.AccAddress {
	return msg.AddedBy
}

// junying -todo, 2019-11-14
//
func (msg MsgAddProfiler) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgAddProfiler) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgAddProfiler) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgAddProfiler) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgAddProfiler) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgAddProfiler) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

//______________________________________________________________________
// MsgDeleteProfiler - struct for delete a profiler
type MsgDeleteProfiler struct {
	DeleteGuardian
}

func NewMsgDeleteProfiler(address, deletedBy sdk.AccAddress) MsgDeleteProfiler {
	return MsgDeleteProfiler{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
			Fee:       sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
			// GasWanted:        uint64(10000),
			// GasPrice:         config.DefaultMinGasPrices,
		},
	}
}

//
func (msg MsgDeleteProfiler) Route() string { return MsgType }

//
func (msg MsgDeleteProfiler) Type() string { return "guardian delete-profiler" }

//
func (msg MsgDeleteProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//
func (msg MsgDeleteProfiler) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

//
func (msg MsgDeleteProfiler) GetSigner() sdk.AccAddress {
	return msg.DeletedBy
}

// junying -todo, 2019-11-14
//
func (msg MsgDeleteProfiler) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgDeleteProfiler) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgDeleteProfiler) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgDeleteProfiler) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgDeleteProfiler) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgDeleteProfiler) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

//______________________________________________________________________
// MsgAddTrustee - struct for add a trustee
type MsgAddTrustee struct {
	AddGuardian
}

func NewMsgAddTrustee(description string, address, addedAddress sdk.AccAddress) MsgAddTrustee {
	return MsgAddTrustee{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedAddress,
			Fee:         sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		},
	}
}

//
func (msg MsgAddTrustee) Route() string { return MsgType }

//
func (msg MsgAddTrustee) Type() string { return "guardian add-trustee" }

//
func (msg MsgAddTrustee) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//
func (msg MsgAddTrustee) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

//
func (msg MsgAddTrustee) GetSigner() sdk.AccAddress {
	return msg.AddedBy
}

// junying -todo, 2019-11-14
//
func (msg MsgAddTrustee) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgAddTrustee) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgAddTrustee) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgAddTrustee) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgAddTrustee) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgAddTrustee) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

//______________________________________________________________________
// MsgDeleteTrustee - struct for delete a trustee
type MsgDeleteTrustee struct {
	DeleteGuardian
}

func NewMsgDeleteTrustee(address, deletedBy sdk.AccAddress) MsgDeleteTrustee {
	return MsgDeleteTrustee{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
			Fee:       sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
			// GasWanted: uint64(10000),
			// GasPrice:  config.DefaultMinGasPrices,
		},
	}
}

//
func (msg MsgDeleteTrustee) Route() string { return MsgType }

//
func (msg MsgDeleteTrustee) Type() string { return "guardian delete-trustee" }

//
func (msg MsgDeleteTrustee) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//
func (msg MsgDeleteTrustee) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

//
func (msg MsgDeleteTrustee) GetSigner() sdk.AccAddress {
	return msg.DeletedBy
}

// junying -todo, 2019-11-14
//
func (msg MsgDeleteTrustee) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgDeleteTrustee) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgDeleteTrustee) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgDeleteTrustee) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgDeleteTrustee) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgDeleteTrustee) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

//______________________________________________________________________
type AddGuardian struct {
	Description string         `json:"description"`
	Address     sdk.AccAddress `json:"address"`  // address added
	AddedBy     sdk.AccAddress `json:"added_by"` // address that initiated the tx
	Fee         sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

//
type DeleteGuardian struct {
	Address   sdk.AccAddress `json:"address"`    // address deleted
	DeletedBy sdk.AccAddress `json:"deleted_by"` // address that initiated the tx
	Fee       sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

//
func (g AddGuardian) ValidateBasic() sdk.Error {
	if len(g.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace)
	}
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.AddedBy) == 0 {
		return sdk.ErrInvalidAddress(g.AddedBy.String())
	}
	if err := g.EnsureLength(); err != nil {
		return err
	}
	return nil
}

//
func (g DeleteGuardian) ValidateBasic() sdk.Error {
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.DeletedBy) == 0 {
		return sdk.ErrInvalidAddress(g.DeletedBy.String())
	}
	return nil
}

//
func (g AddGuardian) EnsureLength() sdk.Error {
	if len(g.Description) > 70 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidGuardian, "description", len(g.Description), 70)
	}
	return nil
}
