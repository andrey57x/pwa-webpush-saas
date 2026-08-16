package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/spf13/viper"
)

type PizzaConfig struct {
	Host       string `mapstructure:"PIZZA_HOST"`
	Port       string `mapstructure:"PIZZA_PORT"`
	SaaSAPIURL string `mapstructure:"PIZZA_SAAS_API_URL"`
	AppCode    string `mapstructure:"PIZZA_APP_CODE"`
}

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := PizzaConfig{
		Host:       viper.GetString("PIZZA_HOST"),
		Port:       viper.GetString("PIZZA_PORT"),
		SaaSAPIURL: viper.GetString("PIZZA_SAAS_API_URL"),
		AppCode:    viper.GetString("PIZZA_APP_CODE"),
	}

	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}
	if cfg.Port == "" {
		cfg.Port = "8085"
	}
	if cfg.SaaSAPIURL == "" {
		cfg.SaaSAPIURL = "http://localhost:8080/api/v1"
	}
	if cfg.AppCode == "" {
		cfg.AppCode = "pizza_app"
	}

	staticDir := "./cmd/pizza/public"

	http.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fmt.Fprintf(w, "self.PIZZA_CONFIG = { saasApiUrl: %q, appCode: %q };\n", cfg.SaaSAPIURL, cfg.AppCode)
	})

	http.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Service-Worker-Allowed", "/")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, filepath.Join(staticDir, "sw.js"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)
		if path == "/" || path == "/index.html" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, path))
	})

	listenAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("Luigi's Pizza App listening on [%s]", listenAddr)
	log.Printf("Target SaaS API URL: %s", cfg.SaaSAPIURL)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Pizza app runtime error: %v", err)
	}
}
