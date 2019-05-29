package login

import (
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"wx/config"
)

var log *logger.Logger

func initLogger() *logger.Logger {
	if log == nil {
		log, _ = logger.New(config.Global.LogColor, config.Global.LogLevel)
	}
	return log
}

func inWhiteList(host string, list []string) bool {
	for _, domain := range list {
		if strings.Contains(host, domain) {
			return true
		}
	}
	return false
}

const typeHtml = "text/html; charset=utf-8"

func L2(c *gin.Context) {
	var log = initLogger()

	q := c.Request.URL.Query()
	l2 := q.Get("l2")
	if l2 == "" {
		c.Data(http.StatusBadRequest, typeHtml, []byte("invalid l2"))
		log.Warning("Missing l2")
		return
	}
	q.Del("l2")
	log.InfoF("l2: %s, other params: %v", l2, q)
	parsedL2Url, err := url.Parse(l2)
	if err != nil {
		c.Data(http.StatusBadRequest, typeHtml, []byte("invalid l2"))
		log.Warning("Invalid l2")
		return
	}

	host := strings.ToLower(parsedL2Url.Hostname())
	if host == "" {
		c.Data(http.StatusBadRequest, typeHtml, []byte("invalid l2"))
		log.Warning("l2 missing host")
		return
	}

	log.DebugF("host: %s", host)
	if !inWhiteList(host, config.Global.ValidDomains) {
		c.Data(http.StatusBadRequest, typeHtml, []byte("invalid l2"))
		log.Warning("l2 not in domain white list")
		return
	}

	sep := "?"
	if strings.Contains(l2, "?") {
		sep = "&"
	}

	redirectUrl := l2 + sep + q.Encode()
	log.InfoF("Redirect to %s", redirectUrl)
	c.Header("Content-Type", typeHtml)
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
