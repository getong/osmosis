package route_test

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/domain"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/domain/mocks"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/router/usecase/pools"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/router/usecase/route"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/router/usecase/routertesting"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v20/x/poolmanager/types"
)

type RouterTestSuite struct {
	routertesting.RouterTestHelper
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

var (
	// Concentrated liquidity constants
	ETH    = routertesting.ETH
	USDC   = routertesting.USDC
	USDT   = routertesting.USDT
	Denom0 = ETH
	Denom1 = USDC

	DefaultCurrentTick = routertesting.DefaultCurrentTick

	DefaultAmt0 = routertesting.DefaultAmt0
	DefaultAmt1 = routertesting.DefaultAmt1

	DefaultCoin0 = routertesting.DefaultCoin0
	DefaultCoin1 = routertesting.DefaultCoin1

	DefaultLiquidityAmt = routertesting.DefaultLiquidityAmt

	// router specific variables
	defaultTickModel = routertesting.DefaultTickModel

	noTakerFee = routertesting.NoTakerFee

	emptyRoute = routertesting.EmptyRoute
)

var (
	DefaultTakerFee     = routertesting.DefaultTakerFee
	DefaultPoolBalances = routertesting.DefaultPoolBalances
	DefaultSpreadFactor = routertesting.DefaultSpreadFactor

	DefaultPool = routertesting.DefaultPool
	EmptyRoute  = routertesting.EmptyRoute

	// Test denoms
	DenomOne   = routertesting.DenomOne
	DenomTwo   = routertesting.DenomTwo
	DenomThree = routertesting.DenomThree
	DenomFour  = routertesting.DenomFour
	DenomFive  = routertesting.DenomFive
	DenomSix   = routertesting.DenomSix
)

// This test validates that the pools in the route are converted into a new serializable
// type for clients with the following list of fields that are returned in each pool:
// - ID
// - Type
// - Balances
// - Spread Factor
// - Token Out Denom
// - Taker Fee
func (s *RouterTestSuite) TestPrepareResultPools() {
	s.Setup()

	balancerPoolID := s.PrepareBalancerPool()

	balancerPool, err := s.App.PoolManagerKeeper.GetPool(s.Ctx, balancerPoolID)
	s.Require().NoError(err)

	testcases := map[string]struct {
		route domain.Route

		expectedPools []domain.RoutablePool
	}{
		"empty route": {
			route: emptyRoute.DeepCopy(),

			expectedPools: []domain.RoutablePool{},
		},
		"single balancer pool in route": {
			route: WithRoutePools(
				emptyRoute,
				[]domain.RoutablePool{
					mocks.WithChainPoolModel(mocks.WithTokenOutDenom(DefaultPool, DenomOne), balancerPool),
				},
			),

			expectedPools: []domain.RoutablePool{
				pools.NewRoutableResultPool(
					balancerPoolID,
					poolmanagertypes.Balancer,
					DefaultPoolBalances,
					DefaultSpreadFactor,
					DenomOne,
					DefaultTakerFee,
				),
			},
		},

		// TODO:
		// add tests with more pool types as well as multiple pools in the route
		// https://app.clickup.com/t/86a1cfwag
	}

	for name, tc := range testcases {
		tc := tc
		s.Run(name, func() {

			resultPools := tc.route.PrepareResultPools()

			s.ValidateRoutePools(tc.expectedPools, resultPools)
			s.ValidateRoutePools(tc.expectedPools, tc.route.GetPools())
		})
	}
}

// This test validates that serialization and deserialization of the route.
// It only exists because I have no idea how JSON marshaling and unmarshaling works
// and Paddy didn't want to teach me the way.
func (s *RouterTestSuite) TestRouteMarshalUnmarshal() {
	allPools := s.PrepareAllSupportedPools()

	concentratedPool, err := s.App.ConcentratedLiquidityKeeper.GetConcentratedPoolById(s.Ctx, allPools.ConcentratedPoolID)
	s.Require().NoError(err)

	concentratedBalance := sdk.NewCoins(sdk.NewCoin(concentratedPool.GetToken0(), DefaultAmt0), sdk.NewCoin(concentratedPool.GetToken1(), DefaultAmt1))

	routablePools := []domain.RoutablePool{
		pools.NewRoutablePool(domain.NewPool(concentratedPool, routertesting.DefaultSpreadFactor, concentratedBalance), concentratedPool.GetToken1(), DefaultTakerFee),
	}

	poolBytes, err := json.Marshal(routablePools)
	s.Require().NoError(err)

	var poolData []pools.SerializedPoolByType
	err = json.Unmarshal(poolBytes, &poolData)
	s.Require().NoError(err)

	routes := []domain.Route{
		&route.RouteImpl{
			Pools: routablePools,
		},
	}

	bytes, err := json.Marshal(routes)
	s.Require().NoError(err)

	var data []route.RouteImpl
	err = json.Unmarshal(bytes, &data)
	s.Require().NoError(err)
}

func WithRoutePools(r domain.Route, pools []domain.RoutablePool) domain.Route {
	return routertesting.WithRoutePools(r, pools)
}
