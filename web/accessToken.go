package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
	"wx/cache"
	"wx/config"
)

var log *logger.Logger

func initLogger() *logger.Logger {
	if log == nil {
		log, _ = logger.New("accessToken", config.Global.LogLevel, config.Global.LogColor)
	}
	return log
}

var aTokStore = cache.NewStore()

type WeixinAccessTokenResponse struct {
	ErrorCode    int    `json:"errcode,omitempty"`
	ErrorMessage string `json:"errmsg,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    uint   `json:"expires_in,omitempty"`
}

func (p WeixinAccessTokenResponse) String() string {
	return fmt.Sprintf(
		"{ErrorCode: %d, ErrorMessage: %s, AccessToken: %s, ExpiresIn: %d}",
		p.ErrorCode,
		p.ErrorMessage,
		p.AccessToken,
		p.ExpiresIn,
	)
}

var requestAccessToken cache.FetchFromSource = func(key string, timeout uint) (value interface{}, ttl uint, err error) {
	var log = initLogger()
	log.NoticeF("request access token for appid: %s", key)
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
			key,
			config.Global.Apps[key],
		),
		nil,
	)
	if err != nil {
		return
	}
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	log.Debug(string(body))

	var o WeixinAccessTokenResponse
	err = json.Unmarshal(body, &o)
	if err != nil {
		return
	}
	if o.ErrorCode > 0 {
		err = errors.New(o.ErrorMessage)
		return
	}
	log.DebugF("%v", o)
	value = o.AccessToken
	ttl = o.ExpiresIn
	return
}

func AccessToken(c *gin.Context) {
	var log = initLogger()
	q := c.Request.URL.Query()
	appId := q.Get("appid")
	if appId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "appid is required",
		})
		log.Warning("missing appid")
		return
	}

	_, ok := config.Global.Apps[appId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no such appid",
		})
		log.Warning("appid is not configured")
		return
	}

	var (
		token interface{}
		ttl   uint
		err   error
	)
	if q.Get("force") != "" {
		token, ttl, err = aTokStore.ForceUpdate(appId, requestAccessToken, 10)
	} else {
		token, ttl, err = aTokStore.Get(appId, requestAccessToken, 10)
	}

	if err != nil {
		aTokStore.Del(appId)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "service temporary unavailable",
		})
		log.WarningF("fetch token error: %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"expires_in":   ttl,
	})
}
