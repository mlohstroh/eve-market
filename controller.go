package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getJSONForItem(t *eveType) map[string]interface{} {
	return map[string]interface{}{
		"id":     t.typeID,
		"name":   t.typeName,
		"volume": t.volume,
		"price":  t.price,
	}
}

func getItem(context *gin.Context) {
	sid := context.Params.ByName("id")
	id, err := strconv.ParseInt(sid, 10, 64)

	if err != nil {
		context.String(http.StatusBadRequest, "ID must be an int")
		return
	}

	t, err := getTypeFromID(id)
	if err != nil {
		context.String(http.StatusNotFound, "typeID not found")
		return
	}

	resp := getJSONForItem(t)

	context.JSON(http.StatusOK, resp)
}

func getItems(context *gin.Context) {
	items := make([]map[string]interface{}, 0)
	for _, t := range typeMap {
		items = append(items, getJSONForItem(t))
	}

	context.JSON(http.StatusOK, items)
}

func (server *Server) getOrders(context *gin.Context) {
	character := context.MustGet("character").(*Character)

	client := server.oauth.Client(server.ctx, character.Token)
	resp, err := client.Get(fmt.Sprintf("%v/v1/characters/%v/orders/", "https://esi.evetech.net", character.CharacterID))
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	context.String(http.StatusOK, string(body))
}
