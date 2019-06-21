package config

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/integrii/flaggy"
	"github.com/mattn/go-isatty"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ValidDomains []string
	Apps         map[string]string
	LogLevel     logger.LogLevel
	LogColor     int
	Addr         net.Addr
}

const version = "1.0"

func (c *Config) String() string {
	return fmt.Sprintf(
		"{ValidDomains: %v, Apps: %v, LogLevel: %s, LogColor: %v, Addr: %v}",
		c.ValidDomains,
		c.Apps,
		[]string{
			"CRITICAL",
			"ERROR",
			"WARNING",
			"NOTICE",
			"INFO",
			"DEBUG",
		}[c.LogLevel-1],
		[]bool{false, true}[c.LogColor],
		c.Addr,
	)
}

var Global = &Config{
	Apps:     make(map[string]string),
	LogColor: 0,
}

func ParseConfig() *Config {
	flaggy.SetVersion(version)

	var validDomains []string
	flaggy.StringSlice(&validDomains, "d", "domains", "valid domains for l2 jump")
	var apps []string
	flaggy.StringSlice(&apps, "a", "apps", "weixin appId:appSecret pairs")
	var logLevel = "INFO"
	flaggy.String(&logLevel, "l", "log-level", "log level")
	var bindIP = ""
	flaggy.String(&bindIP, "b", "bind", "binding ip address default: [::]")
	var port uint16 = 0
	flaggy.UInt16(&port, "p", "port", "binding port number default: 3001")
	var color = false
	flaggy.Bool(&color, "c", "color", "force use color for log output")

	flaggy.Parse()

	for i := range validDomains {
		domain := strings.Trim(validDomains[i], "\r\n \t")
		if domain == "" {
			continue
		}
		domain = strings.ToLower(domain)
		Global.ValidDomains = append(Global.ValidDomains, domain)
	}
	for i := 0; i < len(apps); i++ {
		app := apps[i]
		pair := strings.SplitN(app, ":", 2)
		Global.Apps[pair[0]] = pair[1]
	}
	logLevel = strings.ToUpper(logLevel)
	switch logLevel {
	case "DEBUG":
		Global.LogLevel = logger.DebugLevel
	case "INFO":
		Global.LogLevel = logger.InfoLevel
	case "NOTICE":
		Global.LogLevel = logger.NoticeLevel
	case "WARNING":
		Global.LogLevel = logger.WarningLevel
	case "ERROR":
		Global.LogLevel = logger.ErrorLevel
	case "CRITICAL":
		Global.LogLevel = logger.CriticalLevel
	default:
		Global.LogLevel = logger.InfoLevel
	}
	if color || isatty.IsTerminal(os.Stderr.Fd()) {
		Global.LogColor = 1
	}
	addr := &net.TCPAddr{}
	if bindIP != "" {
		addr.IP = net.ParseIP(bindIP)
	}
	if addr.IP == nil {
		addr.IP = net.ParseIP("[::]")
	}

	if port > 0 {
		addr.Port = int(port)
	} else if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		addr.Port = port
	} else {
		addr.Port = 3001
	}

	Global.Addr = addr

	return Global
}
