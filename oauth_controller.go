package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (server *Server) oauthBegin(c *gin.Context) {
	url := server.oauth.AuthCodeURL("nothing", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (server *Server) oauthCallback(c *gin.Context) {
	code := c.Query("code")

	if len(code) == 0 {
		c.String(http.StatusBadRequest, "No code found in request parameters")
		return
	}

	token, err := server.oauth.Exchange(server.ctx, code)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	client := server.oauth.Client(server.ctx, token)
	resp, err := client.Get("https://login.eveonline.com/oauth/verify")
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
