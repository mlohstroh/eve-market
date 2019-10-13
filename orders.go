package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ESIOrder is the db model and ESI representation for
type ESIOrder struct {
	Duration     int       `json:"duration"`
	IsBuyOrder   bool      `json:"is_buy_order"`
	Issued       time.Time `json:"issued"`
	Location     int       `json:"location_id"`
	MinVolume    int       `json:"min_volume"`
	OrderID      int64     `json:"order_id" bson:"_id"`
	Price        float64   `json:"price"`
	Range        string    `json:"range"`
	SystemID     int       `json:"system_id"`
	TypeID       int64     `json:"type_id"`
	VolumeRemain int       `json:"volume_remain"`
	VolumeTotal  int       `json:"volume_total"`

	// Custom Types
	TypeInfo *eveType `json:"type_info"`
}

func (server *Server) saveOrders(orders []*ESIOrder) error {
	for _, o := range orders {
		t, err := getTypeFromID(o.TypeID)
		o.TypeInfo = t
		if err != nil {
			continue
		}
		_, err = server.db.Collection("orders").ReplaceOne(server.ctx, bson.M{"_id": o.OrderID}, o, options.Replace().SetUpsert(true))
		if err != nil {
			continue
		}
	}

	return nil
}
