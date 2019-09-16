package main

import (
	"fmt"
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

// ServeHTTP demonstrates a plugin that handles HTTP requests.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("canary")
	if err != nil {
		p.API.LogInfo("Canary cookie does not exist. Setting up.")
		err := p.handleCookie(w)
		if err != nil {
			p.API.LogError("Unable to set cookie" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if cookie.Value == "never" {
		p.API.LogInfo("Canary cookie is set to never.")
		err := p.handleCookie(w)
		if err != nil {
			p.API.LogError("Unable to set cookie" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// handleCookie function is used to check the canary percentage and call the addCookie function.
func (p *Plugin) handleCookie(w http.ResponseWriter) error {
	config := p.getConfiguration()
	percentage, err := strconv.Atoi(config.CanaryPercentage)
	if err != nil {
		return err
	}
	randomNumber := rand.Intn(100)
	if randomNumber < percentage {
		p.API.LogInfo("Setting Canary cookie to always.")
		p.addCookie(w, "always")
	} else {
		p.API.LogInfo("Setting Canary cookie to never.")
		p.addCookie(w, "never")
	}
	return nil
}

// addCookie function is used to set the canary cookie.
func (p *Plugin) addCookie(w http.ResponseWriter, cookieValue string) {
	expire := time.Now().AddDate(0, 0, 1)
	canaryCookie := http.Cookie{
		Name:    "canary",
		Value:   cookieValue,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &canaryCookie)
	fmt.Fprint(w, "Setting cookie for Canary build!")
}
