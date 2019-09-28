package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

type apiResponse struct {
	CookieName  string `json:"cookieName"`
	CookieValue string `json:"cookieValue"`
	StatusCode  int    `json:"statusCode"`
}

// ServeHTTP handles HTTP requests.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/api/v1/check":
		p.handleCookieCheck(w, r)
	}
}

// handleCookieCheck function is used to handle the checks for the existense of the Canary cookie.
func (p *Plugin) handleCookieCheck(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	response := apiResponse{"", "", http.StatusOK}
	cookie, err := r.Cookie("canary")
	if err != nil {
		p.API.LogDebug("Canary cookie does not exist for user " + userID)
		response, err = p.addCanaryCookieBasedOnPercentage(w, userID)
		if err != nil {
			p.API.LogError("Canary percentage configuration missing." + err.Error())
			http.Error(w, "Canary percentage configuration missing.", http.StatusInternalServerError)
			return
		}
	} else if cookie.Value == "never" {
		p.API.LogDebug("Canary cookie is set to never for user " + userID)
		response, err = p.addCanaryCookieBasedOnPercentage(w, userID)
		if err != nil {
			p.API.LogError("Canary percentage configuration missing." + err.Error())
			http.Error(w, "Canary percentage configuration missing.", http.StatusInternalServerError)
			return
		}
	}
	b, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(b)
}

// canaryPercentage function is used to check the canary percentage and call the addCookie function.
func (p *Plugin) addCanaryCookieBasedOnPercentage(w http.ResponseWriter, userID string) (apiResponse, error) {
	config := p.getConfiguration()
	percentage, err := strconv.Atoi(config.CanaryPercentage)
	if err != nil {
		return apiResponse{"", "", http.StatusInternalServerError}, err
	}
	randomNumber := rand.Intn(100)
	if randomNumber < percentage {
		p.API.LogDebug("Setting Canary cookie to always for user " + userID)
		response := p.addCookie(w, "canary", "always")
		return response, nil
	} else {
		p.API.LogDebug("Setting Canary cookie to never for user " + userID)
		response := p.addCookie(w, "canary", "never")
		return response, nil
	}
}

// addCookie function is used to set the canary cookie.
func (p *Plugin) addCookie(w http.ResponseWriter, cookieName, cookieValue string) apiResponse {
	expire := time.Now().AddDate(0, 0, 1)
	canaryCookie := http.Cookie{
		Name:    cookieName,
		Value:   cookieValue,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &canaryCookie)
	response := apiResponse{cookieName, cookieValue, http.StatusOK}
	return response
}
