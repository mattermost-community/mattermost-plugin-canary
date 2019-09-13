package main

import (
	"fmt"
	"math/rand"
	"net/http"
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

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("canary")
	if err != nil {
		p.handleCookie(w)
	} else if cookie.Value == "never" {
		p.handleCookie(w)
	}
}

func (p *Plugin) handleCookie(w http.ResponseWriter) {
	config := p.getConfiguration()
	randomNumber := rand.Intn(100)
	if randomNumber < config.CanaryPercentage {
		p.addCookie(w, "always")
	} else {
		p.addCookie(w, "never")
	}
}

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
