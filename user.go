package main

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

// Character is the DB representation of an Eve Character
type Character struct {
	ID          string
	CharacterID int
	Name        string
	Token       *oauth2.Token
}

// CreateOrUpdateUserFromESI makes sure that user is logged in
func (server *Server) CreateOrUpdateUserFromESI(esiData []byte, token *oauth2.Token) (*Character, error) {
	esiChar := make(map[string]interface{})
	err := json.Unmarshal(esiData, &esiChar)
	if err != nil {
		return nil, err
	}

	characterID := int(esiChar["CharacterID"].(float64))
	characterName := esiChar["CharacterName"].(string)

	collection := server.db.Collection("users")
	filter := bson.M{"CharacterID": characterID}

	existing := &Character{}
	err = collection.FindOne(server.ctx, filter).Decode(&existing)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err != nil {
		existing = &Character{}
	}

	existing.CharacterID = characterID
	existing.Name = characterName
	existing.Token = token

	insertFilter := bson.M{
		"CharacterID": characterID,
	}

	update := bson.M{
		"$set": bson.M{
			"CharacterID": characterID,
			"Name":        characterName,
			"Token":       token,
		},
	}

	_, err = collection.UpdateOne(server.ctx, insertFilter, update, options.Update().SetUpsert(true))

	return existing, err
}
