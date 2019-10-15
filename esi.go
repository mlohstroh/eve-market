package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// DeserializeFunc is the function is a wrapper to help the getAllPages function
// properly deserialize things into an array
type DeserializeFunc func([]byte) (int, error)

// TypeID is just a quick typedef for an int
type TypeID int64

func urlFor(path string, args ...interface{}) string {
	return fmt.Sprintf("%v%v", "https://esi.evetech.net", fmt.Sprintf(path, args...))
}

func getAllRegionTypes(region int, client *http.Client) ([]TypeID, error) {
	url := urlFor("/v1/markets/%v/types/", region)

	allTypes := make([]TypeID, 0)
	err := getAllPages(url, nil, client, func(data []byte) (int, error) {
		var types []TypeID
		err := json.Unmarshal(data, &types)
		allTypes = append(allTypes, types...)
		return len(types), err
	})
	return allTypes, err
}

func getAllStructureOrders(structureID int, client *http.Client) ([]*ESIOrder, error) {
	url := urlFor("/v1/markets/structures/%v/", structureID)
	allOrders := make([]*ESIOrder, 0)
	err := getAllPages(url, nil, client, func(data []byte) (int, error) {
		var orders []*ESIOrder
		err := json.Unmarshal(data, &orders)
		if err != nil {
			log.Printf("Error %v", err)
		}
		allOrders = append(allOrders, orders...)
		return len(orders), err
	})
	return allOrders, err
}

func getAllRegionOrders(region int, client *http.Client) ([]*ESIOrder, error) {
	url := urlFor("/v1/markets/%v/orders/", region)
	allOrders := make([]*ESIOrder, 0)
	err := getAllPages(url, nil, client, func(data []byte) (int, error) {
		var orders []*ESIOrder
		err := json.Unmarshal(data, &orders)
		if err != nil {
			log.Printf("Error %v", err)
		}
		allOrders = append(allOrders, orders...)
		return len(orders), err
	})

	return allOrders, err
}

func getAllPages(path string, params []string, client *http.Client, f DeserializeFunc) error {
	currentPage := 1
	// will be set later
	maxPages := 0

	if params == nil {
		params = make([]string, 0)
	}

	for {
		copiedParams := make([]string, len(params))
		copy(copiedParams, params)
		copiedParams = append(copiedParams, fmt.Sprintf("page=%v", currentPage))
		if maxPages != 0 && currentPage > maxPages {
			break
		}

		q := strings.Join(copiedParams, "&")
		url := fmt.Sprintf("%v?%v", path, q)
		log.Printf("ESI - Calling %v", url)
		resp, err := client.Get(url)
		currentPage++
		if err != nil {
			return err
		}
		maxPages, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			maxPages = 0
		}
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		count, err := f(body)

		// If none were read
		if count == 0 {
			break
		}
	}

	return nil
}
