package mempool1559

import sdk "github.com/cosmos/cosmos-sdk/types"

// Sections to this right now:
// - Maintain something thats gets parsed from chain tx execution
// update eipState according to that.
// - Every time blockheight % 1000 = 0, reset eipState to default. (B/c this isn't committed on-chain, still gives us some consistency guarantees)
// - Evaluate CheckTx/RecheckTx against this.
//
// 1000 blocks = almost 2 hours, maybe we need a smaller time for resets?
//
// PROBLEMS: Currently, a node will throw out any tx that gets under its gas bound here.
// :OOO We can just do this on checkTx not recheck
//
// Variables we can control for:
// - fees paid per unit gas
// - gas wanted per block (Ethereum)
// - gas used and gas wanted difference
// TODO: Change this percentage update time to be faster

// TODO: Read this from config, can even make default 0, so this is only turned on by nodes who change it!
// ALt: do that with an enable/disable flag.
var DefaultBaseFee = sdk.MustNewDecFromStr("0.0025")
var MinBaseFee = sdk.MustNewDecFromStr("0.0025")
var TargetGas = int64(90_000_000)
var MaxBlockChangeRate = sdk.NewDec(1).Quo(sdk.NewDec(16))
var ResetInterval = int64(1000)

type eipState struct {
	// Signal when we are starting a new block
	// TODO: Or just use begin block
	lastBlockHeight         int64
	totalGasWantedThisBlock int64

	CurBaseFee sdk.Dec
}

var CurEipState = eipState{
	lastBlockHeight:         0,
	totalGasWantedThisBlock: 0,
	CurBaseFee:              DefaultBaseFee,
}

func (e *eipState) startBlock(height int64) {
	e.lastBlockHeight = height
	e.totalGasWantedThisBlock = 0

	if height%ResetInterval == 0 {
		e.CurBaseFee = DefaultBaseFee
	}
}

func (e *eipState) deliverTxCode(_ sdk.Context, tx sdk.FeeTx) {
	e.totalGasWantedThisBlock += int64(tx.GetGas())
}

// Equation is:
// baseFeeMultiplier = 1 + (gasUsed - targetGas) / targetGas * maxChangeRate
// newBaseFee = baseFee * baseFeeMultiplier
func (e *eipState) updateBaseFee(height int64) {
	gasUsed := e.totalGasWantedThisBlock
	// obvi fix
	e.lastBlockHeight = height
	gasDiff := gasUsed - TargetGas
	//  (gasUsed - targetGas) / targetGas * maxChangeRate
	baseFeeIncrement := sdk.NewDec(gasDiff).Quo(sdk.NewDec(TargetGas)).Mul(MaxBlockChangeRate)
	baseFeeMultiplier := sdk.NewDec(1).Add(baseFeeIncrement)
	e.CurBaseFee.MulMut(baseFeeMultiplier)

	// Make a min base fee
	if e.CurBaseFee.LT(MinBaseFee) {
		e.CurBaseFee = MinBaseFee
	}
}

func (e *eipState) GetCurBaseFee() sdk.Dec {
	return e.CurBaseFee
}
