package main

import (
	"log"
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

func (server *Server) getAllTypes(structureID int) ([]int64, error) {
	query := bson.M{
		"location": structureID,
	}
	types, err := server.db.Collection("orders").Distinct(server.ctx, "typeid", query)
	if err != nil {
		return nil, err
	}

	casted := make([]int64, len(types))
	for i, t := range types {
		// How's that for a gross cast?
		casted[i] = t.(int64)
	}

	return casted, nil
}

func (server *Server) getOrdersByType(typeID int64) ([]*ESIOrder, error) {
	query := bson.M{
		"typeid": typeID,
	}
	cursor, err := server.db.Collection("orders").Find(server.ctx, query)
	if err != nil {
		return nil, err
	}

	orders := make([]*ESIOrder, 0)
	defer cursor.Close(server.ctx)
	for cursor.Next(server.ctx) {
		o := &ESIOrder{}
		err = cursor.Decode(&o)
		if err != nil {
			log.Printf("Unable to decode order due to %v", err)
			continue
		}
		orders = append(orders, o)
	}

	return orders, cursor.Err()
}
