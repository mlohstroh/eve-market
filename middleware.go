package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (server *Server) requireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		cID, err := c.Cookie("character_id")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id, err := strconv.Atoi(cID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		filter := bson.M{
			"CharacterID": id,
		}

		character := &Character{}
		err = server.db.Collection("users").FindOne(server.ctx, filter).Decode(&character)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("character", character)
		c.Next()
	}
}
