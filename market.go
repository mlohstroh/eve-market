package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const esiRoot = "https://esi.evetech.net"

// typeID,groupID,typeName,description,mass,volume,capacity,portionSize,raceID,basePrice,published,marketGroupID,iconID,soundID,graphicID
type eveType struct {
	typeID      int64
	groupID     int64
	typeName    string
	description string
	mass        float64
	volume      float64
	// pricing data not in SDE
	price *evePrice
	// not parsed
	capacity      float64
	portionSize   int64
	raceID        string
	basePrice     string
	published     int
	marketGroupID string
	iconID        string
	soundID       string
	graphicID     string
}

type evePrice struct {
	AdjustedPrice float64 `json:"adjusted_price"`
	AveragePrice  float64 `json:"average_price"`
	TypeID        int64   `json:"type_id"`
}

var (
	typeMap map[int64]*eveType
)

func loadSDE() error {
	data, err := ioutil.ReadFile("data/invTypes.csv")

	if err != nil {
		return err
	}

	reader := csv.NewReader(bytes.NewReader(data))

	typeMap = make(map[int64]*eveType)

	var header bool
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// Skip the header, ugly
		if !header {
			header = true
			continue
		}

		// Yay for no struct marshalling without a library
		typeID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return err
		}
		groupID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}
		typeName := record[2]
		description := record[3]
		mass, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return err
		}
		volume, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return err
		}

		t := &eveType{
			typeID:      typeID,
			groupID:     groupID,
			typeName:    typeName,
			description: description,
			mass:        mass,
			volume:      volume,
		}

		// Add to global array
		typeMap[typeID] = t
	}

	return nil
}

func getTypeFromID(id int64) (*eveType, error) {
	t, ok := typeMap[id]

	if ok {
		return t, nil
	}

	return nil, errors.New("Type was not found")
}

func getPrices() error {
	resp, err := get("/v1/markets/prices/")

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data/prices.json", body, 0777)
	if err != nil {
		return err
	}

	var prices []*evePrice
	err = json.Unmarshal(body, &prices)

	if err != nil {
		return err
	}

	for _, p := range prices {
		t, err := getTypeFromID(p.TypeID)

		if err != nil {
			continue
		}

		t.price = p
	}

	return nil
}

func get(path string) (*http.Response, error) {
	return http.DefaultClient.Get(fmt.Sprintf("%v%v", esiRoot, path))
}
