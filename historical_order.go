package main

import (
	"fmt"
	"time"
)

// HistoricalOrder is a DB model for recording an Eve type's volume,
// price, and whatever else should be recorded
type HistoricalOrder struct {
	HighSellPrice float64
	LowSellPrice  float64
	HighBuyPrice  float64
	LowBuyPrice   float64
	Orders        int64
	Volume        int64
	TypeID        int64
	Location      int64
	CreatedAt     time.Time
	TypeInfo      *eveType
}

const (
	// Inf is a large constant representing positive infininty
	Inf float64 = 999999999999999999999
	// NInf is a large constant represneting negative infinity
	NInf float64 = -999999999999999999999
)

func (server *Server) createHourlyReport(orders []*ESIOrder, location, typeID int64, createdAt time.Time) error {
	var (
		highsell    float64 = 0
		lowsell     float64 = Inf
		highbuy     float64 = 0
		lowbuy      float64 = Inf
		totalVolume int64   = 0
		orderCount  int64   = int64(len(orders))
	)

	if orderCount == 0 {
		return fmt.Errorf("No orders for this hour on type %v", typeID)
	}

	for _, o := range orders {
		totalVolume += int64(o.VolumeRemain)

		if o.IsBuyOrder {
			if o.Price > highbuy {
				highbuy = o.Price
			}
			if o.Price < lowbuy {
				lowbuy = o.Price
			}
		} else {
			if o.Price > highsell {
				highsell = o.Price
			}
			if o.Price < lowsell {
				lowsell = o.Price
			}
		}
	}

	if lowbuy > Inf-1 {
		lowbuy = 0
	}
	if highbuy > Inf-1 {
		highbuy = 0
	}

	info, err := getTypeFromID(typeID)

	if err != nil {
		return err
	}

	history := &HistoricalOrder{
		HighSellPrice: highsell,
		LowSellPrice:  lowsell,
		HighBuyPrice:  highbuy,
		LowBuyPrice:   lowbuy,
		Orders:        orderCount,
		Volume:        totalVolume,
		TypeID:        typeID,
		CreatedAt:     createdAt,
		TypeInfo:      info,
	}

	_, err = server.db.Collection("history").InsertOne(server.ctx, history)

	return err
}
