package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getStructureOrders(c *gin.Context) {
	character := c.MustGet("character").(*Character)
	structureID := c.Params.ByName("structure_id")

	client := server.oauth.Client(server.ctx, character.Token)
	resp, err := client.Get(fmt.Sprintf("%v/v1/markets/structures/%v/", "https://esi.evetech.net", structureID))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, string(body))
}
