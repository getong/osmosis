package poolmanager

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v25/x/poolmanager/types"
)

var IntMaxValue = intMaxValue

func (k Keeper) GetNextPoolIdAndIncrement(ctx sdk.Context) uint64 {
	return k.getNextPoolIdAndIncrement(ctx)
}

// SetPoolRoutesUnsafe sets the given routes to the poolmanager keeper
// to allow routing from a pool type to a certain swap module.
// For example, balancer -> gamm.
// This utility function is only exposed for testing and should not be moved
// outside of the _test.go files.
func (k *Keeper) SetPoolRoutesUnsafe(routes map[types.PoolType]types.PoolModuleI) {
	k.routes = routes
	k.cachedPoolModules = &sync.Map{}
}

// SetPoolModulesUnsafe sets the given modules to the poolmanager keeper.
// This utility function is only exposed for testing and should not be moved
// outside of the _test.go files.
func (k *Keeper) SetPoolModulesUnsafe(poolModules []types.PoolModuleI) {
	k.poolModules = poolModules
	k.cachedPoolModules = &sync.Map{}
}

func (k Keeper) GetAllPoolRoutes(ctx sdk.Context) []types.ModuleRoute {
	return k.getAllPoolRoutes(ctx)
}

func (k Keeper) ValidateCreatedPool(ctx sdk.Context, poolId uint64, pool types.PoolI) error {
	return k.validateCreatedPool(poolId, pool)
}

func (k Keeper) CreateMultihopExpectedSwapOuts(
	ctx sdk.Context,
	route []types.SwapAmountOutRoute,
	tokenOut sdk.Coin,
) ([]osmomath.Int, error) {
	return k.createMultihopExpectedSwapOuts(ctx, route, tokenOut)
}

func (k Keeper) TrackVolume(ctx sdk.Context, poolId uint64, volumeGenerated sdk.Coin) {
	k.trackVolume(ctx, poolId, volumeGenerated)
}

func (k Keeper) ChargeTakerFee(ctx sdk.Context, tokenIn sdk.Coin, tokenOutDenom string, sender sdk.AccAddress, exactIn bool) (sdk.Coin, sdk.Coin, error) {
	return k.chargeTakerFee(ctx, tokenIn, tokenOutDenom, sender, exactIn)
}

func (k Keeper) QueryAndCheckAlloyedDenom(ctx sdk.Context, contractAddr sdk.AccAddress) (string, error) {
	return k.queryAndCheckAlloyedDenom(ctx, contractAddr)
}

func (k Keeper) SnapshotTakerFeeShareAlloyComposition(ctx sdk.Context, contractAddr sdk.AccAddress) ([]types.TakerFeeShareAgreement, error) {
	return k.snapshotTakerFeeShareAlloyComposition(ctx, contractAddr)
}

func (k Keeper) RecalculateAndSetTakerFeeShareAlloyComposition(ctx sdk.Context, poolId uint64) error {
	return k.recalculateAndSetTakerFeeShareAlloyComposition(ctx, poolId)
}
