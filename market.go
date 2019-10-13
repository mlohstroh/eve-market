package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const esiRoot = "https://esi.evetech.net"

// typeID,groupID,typeName,description,mass,volume,capacity,portionSize,raceID,basePrice,published,marketGroupID,iconID,soundID,graphicID
type eveType struct {
	TypeID      int64
	GroupID     int64
	TypeName    string
	description string
	Mass        float64
	Volume      float64
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
			TypeID:      typeID,
			GroupID:     groupID,
			TypeName:    typeName,
			description: description,
			Mass:        mass,
			Volume:      volume,
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

	return nil
}

var (
	// Jita
	theForge = 10000002
	// 1DQ Imperial Palace
	importantLocation = []int{60003760, 1030049082711}
	// The Forge
	importantRegions = []int{10000002}
)

func (server *Server) backgroundGetStructureOrders() error {
	character, err := server.GetAnyCharacter()
	if err != nil {
		return err
	}

	client := server.oauth.Client(server.ctx, character.Token)

	// Public structures
	for _, region := range importantRegions {
		orders, err := getAllRegionOrders(region, client)
		if err != nil {
			log.Printf("Unable to fetch orders from region %v. Skipping... Error: %v", region, err)
			continue
		}

		// Ugly filter down
		filteredOrders := make([]*ESIOrder, 0)
		for _, o := range orders {
			if ContainsI(o.Location, importantLocation) {
				filteredOrders = append(filteredOrders, o)
			}
		}

		log.Printf("Got %v orders from region %v", len(filteredOrders), region)

		profileFunction("Save Orders", func() {
			err = server.saveOrders(filteredOrders)
			if err != nil {
				log.Printf("Error saving all orders from region %v. %v", region, err)
			}
		})
	}

	// Private structures
	for _, location := range importantLocation {
		orders, err := getAllStructureOrders(location, client)
		if err != nil {
			log.Printf("Unable to fetch orders from location %v. Skipping... Error: %v", location, err)
			continue
		}
		log.Printf("Got %v orders for location %v", len(orders), location)

		err = server.saveOrders(orders)
		if err != nil {
			log.Printf("Error saving all orders from location %v. %v", location, err)
		}
	}

	return nil
}

func get(path string) (*http.Response, error) {
	return http.DefaultClient.Get(fmt.Sprintf("%v%v", esiRoot, path))
}
