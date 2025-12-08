package bollinger_bands_breakout_strategy

import (
	"fmt"
	"math"
	"github.com/shopspring/decimal"
	"trading/libs/7_common/types"
	"trading/libs/7_common/utils"
)

type BollingerBands struct {
	period  int
	stddev  decimal.Decimal
	prices  []decimal.Decimal
	current decimal.Decimal
}

func NewBollingerBands(period int, stddev float64) *BollingerBands {
	return &BollingerBands{
		period: period,
		stddev: decimal.NewFromFloat(stddev),
		prices: make([]decimal.Decimal, 0),
	}
}

func (b *BollingerBands) Update(price decimal.Decimal) (upper, lower, middle decimal.Decimal, ready bool) {
	b.prices = append(b.prices, price)
	if len(b.prices) > b.period {
		b.prices = b.prices[1:]
	}
	if len(b.prices) < b.period {
		return decimal.Zero, decimal.Zero, decimal.Zero, false
	}
	var sum decimal.Decimal
	for _, p := range b.prices {
		sum = sum.Add(p)
	}
	middle = sum.Div(decimal.NewFromFloat(float64(b.period)))
	var varianceSum decimal.Decimal
	for _, p := range b.prices {
		varianceSum = varianceSum.Add(p.Sub(middle).Pow(decimal.NewFromFloat(2)))
	}
	variance := varianceSum.Div(decimal.NewFromFloat(float64(b.period)))
	stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
	b.current = price
	bandWidth := stdDev.Mul(b.stddev)
	upper = middle.Add(bandWidth)
	lower = middle.Sub(bandWidth)
	return upper, lower, middle, true
}

type BollingerBandsBreakoutParams struct {
	WindowSize int         `mapstructure:"window_size"`
	StdDev     float64     `mapstructure:"stddev"`
	Pair       types.Pair  `mapstructure:"pair"`
}

var _ types.IStrategy = (*BollingerBandsBreakoutStrategy)(nil)

type BollingerBandsBreakoutStrategy struct {
	stman     types.IStrategyManager
	params    *BollingerBandsBreakoutParams
	bbands    *BollingerBands
	signalDelivered bool
}

func NewBollingerBandsBreakoutStrategy() types.IStrategy {
	return &BollingerBandsBreakoutStrategy{}
}

func (st *BollingerBandsBreakoutStrategy) Init(stman types.IStrategyManager, raw map[string]any) error {
	st.stman = stman
	cfg, err := utils.DecodeParams[BollingerBandsBreakoutParams](raw)
	if err != nil {
		return fmt.Errorf("failed to decode params: %w", err)
	}
	st.params = cfg
	st.bbands = NewBollingerBands(cfg.WindowSize, cfg.StdDev)
	return nil
}

func (st *BollingerBandsBreakoutStrategy) OnStart(sctx types.ISmartContext) error {
	sctx.Infof("Bollinger Bands breakout strategy started for %s", st.params.Pair)
	return nil
}

func (st *BollingerBandsBreakoutStrategy) OnStopped(sctx types.ISmartContext) error {
	sctx.Infof("Bollinger Bands breakout strategy stopped")
	return nil
}

func (st *BollingerBandsBreakoutStrategy) OnTrade(sctx types.ISmartContext, trade *types.Trade) error {
	return nil
}

func (st *BollingerBandsBreakoutStrategy) OnCandle(sctx types.ISmartContext, candle *types.Candle) error {
	upper, lower, _, ready := st.bbands.Update(candle.Close)
	if !ready {
		return nil
	}

	price := candle.Close
	// Breakout above the upper band
	if price.GreaterThan(upper) {
		if pos, _ := st.stman.GetAvailableSizeToSell(sctx, st.params.Pair); pos.LessThanOrEqual(decimal.Zero) {
			size, _ := st.stman.GetAvailableSizeToBuy(sctx, st.params.Pair, price)
			if size.GreaterThan(decimal.Zero){
				return st.stman.RegisterOrder(sctx, &types.Order{
					Type:      types.Market,
					Side:      types.Buy,
					Pair:      st.params.Pair,
					Size:      size,
					Price:     price,
					Timestamp: candle.Time,
				})
			}
		}
		st.signalDelivered = true
	}

	// Breakout below the lower band
	if price.LessThan(lower) {
		if pos, _ := st.stman.GetAvailableSizeToSell(sctx, st.params.Pair); pos.GreaterThanOrEqual(decimal.Zero) {
			return st.stman.RegisterOrder(sctx, &types.Order{
				Type:      types.Market,
				Side:      types.Sell,
				Pair:      st.params.Pair,
				Size:      pos,
				Price:     price,
				Timestamp: candle.Time,
			})
		}
		st.signalDelivered = true
	}
	return nil
}
