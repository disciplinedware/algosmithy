package ema_crossover

import (
	"fmt"
	"github.com/shopspring/decimal"
	"trading/libs/7_common/types"
	"trading/libs/7_common/utils"
)

type EMA struct {
	period int
	k      decimal.Decimal
	value  decimal.Decimal
	count  int
}

func NewEMA(period int) *EMA {
	k := decimal.NewFromFloat(2.0 / (float64(period) + 1.0))
	return &EMA{period: period, k: k}
}

func (e *EMA) Update(price decimal.Decimal) decimal.Decimal {
	if e.count == 0 {
		e.value = price
	} else {
		e.value = price.Mul(e.k).Add(e.value.Mul(decimal.NewFromFloat(1).Sub(e.k)))
	}
	e.count++
	return e.value
}

func (e *EMA) Ready() bool {
	return e.count >= e.period
}

type EmaCrossoverParams struct {
	ShortWindow  int        `mapstructure:"short_window"`
	LongWindow   int        `mapstructure:"long_window"`
	SignalWindow int        `mapstructure:"signal_window"`
	Pair         types.Pair `mapstructure:"pair"`
}

var _ types.IStrategy = (*EmaCrossover)(nil)

type EmaCrossover struct {
	stman  types.IStrategyManager
	params *EmaCrossoverParams

	emaShort  *EMA
	emaLong   *EMA
	emaSignal *EMA

	lastMACD decimal.Decimal
}

func NewEmaCrossover() types.IStrategy {
	return &EmaCrossover{}
}

func (st *EmaCrossover) Init(stman types.IStrategyManager, raw map[string]any) error {
	st.stman = stman

	cfg, err := utils.DecodeParams[EmaCrossoverParams](raw)
	if err != nil {
		return fmt.Errorf("failed to decode params: %w", err)
	}
	st.params = cfg

	st.emaShort = NewEMA(cfg.ShortWindow)
	st.emaLong = NewEMA(cfg.LongWindow)
	st.emaSignal = NewEMA(cfg.SignalWindow)

	return nil
}

func (st *EmaCrossover) OnStart(sctx types.ISmartContext) error {
	sctx.Infof("EMA crossover strategy started for %s", st.params.Pair)
	return nil
}

func (st *EmaCrossover) OnStopped(sctx types.ISmartContext) error {
	sctx.Infof("EMA crossover strategy stopped")
	return nil
}

func (st *EmaCrossover) OnTrade(sctx types.ISmartContext, trade *types.Trade) error {
	return nil
}

func (st *EmaCrossover) OnCandle(sctx types.ISmartContext, candle *types.Candle) error {
	price := candle.Close

	short := st.emaShort.Update(price)
	long := st.emaLong.Update(price)
	if !st.emaShort.Ready() || !st.emaLong.Ready() {
		return nil
	}

	macd := short.Sub(long)
	signal := st.emaSignal.Update(macd)
	if !st.emaSignal.Ready() {
		return nil
	}

	prev := st.lastMACD
	st.lastMACD = macd

	// --- Buy Signal ---
	if prev.LessThan(signal) && macd.GreaterThan(signal) {
		pos, err := st.stman.GetAvailableSizeToSell(sctx, st.params.Pair)
		if err != nil {
			sctx.Errorf("GetAvailableSizeToSell failed: %v", err)
			return err
		}
		if pos.LessThanOrEqual(decimal.Zero) {
			size, err := st.stman.GetAvailableSizeToBuy(sctx, st.params.Pair, price)
			if err != nil {
				sctx.Errorf("GetAvailableSizeToBuy failed: %v", err)
				return err
			}
			if size.GreaterThan(decimal.Zero) {
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
	}

	// --- Sell Signal ---
	if prev.GreaterThan(signal) && macd.LessThan(signal) {
		pos, err := st.stman.GetAvailableSizeToSell(sctx, st.params.Pair)
		if err != nil {
			sctx.Errorf("GetAvailableSizeToSell failed: %v", err)
			return err
		}
		if pos.GreaterThan(decimal.Zero) {
			return st.stman.RegisterOrder(sctx, &types.Order{
				Type:      types.Market,
				Side:      types.Sell,
				Pair:      st.params.Pair,
				Size:      pos,
				Price:     price,
				Timestamp: candle.Time,
			})
		}
	}

	return nil
}
