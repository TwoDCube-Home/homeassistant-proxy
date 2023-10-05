package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	HomeAssistantURL string `envconfig:"HOME_ASSISTANT_URL"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		logger.Fatal("Failed processing environment variables", zap.Error(err))
	}

	url, err := url.Parse(c.HomeAssistantURL)
	if err != nil {
		logger.Fatal("Failed parsing url HOME_ASSISTANT_URL", zap.Error(err))
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		req.Header.Set("X-Forwarded-Proto", "https")
		logger.Info(
			"Forwarding request",
			zap.String("url", req.URL.String()),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("x_forwarded_for", req.Header.Get("X-Forwarded-For")),
			zap.Any("x_forwarded_proto", req.Header["X-Forwarded-Proto"]),
		)
		director(req)
	}
	http.Handle("/", proxy)
	logger.Fatal("Failed to start http server", zap.Error(http.ListenAndServe(":6969", nil)))
}
