package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed web/*
var webFS embed.FS

func main() {
	listen := flag.String("listen", ":3456", "HTTP listen address")
	home := flag.String("home", "", "canined home (default ~/.canine-nico)")
	chainID := flag.String("chain-id", "jackal-nico-1", "chain id")
	rpc := flag.String("rpc", "http://127.0.0.1:26657", "comet RPC")
	rest := flag.String("rest", "http://127.0.0.1:1317", "REST API")
	provider := flag.String("provider", "http://127.0.0.1:3333", "nico-provider base URL")
	canined := flag.String("canined", "canined", "canined binary")
	key := flag.String("key", "validator", "keyring key for txs")
	keyring := flag.String("keyring-backend", "test", "keyring backend")
	flag.Parse()

	if *home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		*home = filepath.Join(h, ".canine-nico")
	}

	lab := &Lab{
		Home:    *home,
		ChainID: *chainID,
		RPC:     *rpc,
		REST:    *rest,
		Provider: *provider,
		Canined: *canined,
		Key:     *key,
		Keyring: *keyring,
		Client:  &http.Client{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", lab.handleHealth)
	mux.HandleFunc("/api/status", lab.handleStatus)
	mux.HandleFunc("/api/balance", lab.handleBalance)
	mux.HandleFunc("/api/providers", lab.handleProviders)
	mux.HandleFunc("/api/storage/params", lab.handleStorageParams)
	mux.HandleFunc("/api/storage/payment", lab.handleStoragePayment)
	mux.HandleFunc("/api/storage/buy", lab.handleBuyStorage)
	mux.HandleFunc("/api/storage/upload", lab.handleUpload)
	mux.HandleFunc("/api/storage/files", lab.handleChainFiles)
	mux.HandleFunc("/api/provider/status", lab.handleProviderStatus)
	mux.HandleFunc("/api/provider/files", lab.handleProviderFiles)
	mux.HandleFunc("/api/provider/download/", lab.handleProviderDownload)
	mux.HandleFunc("/api/filetree/params", lab.handleFiletreeParams)
	mux.HandleFunc("/api/filetree/provision", lab.handleFiletreeProvision)
	mux.HandleFunc("/api/rns/params", lab.handleRNSParams)
	mux.HandleFunc("/api/oracle/feeds", lab.handleOracleFeeds)
	mux.HandleFunc("/api/oracle/feed", lab.handleOracleCreateFeed)
	mux.HandleFunc("/api/jklmint/params", lab.handleJklmintParams)
	mux.HandleFunc("/api/wasm/codes", lab.handleWasmCodes)
	mux.HandleFunc("/api/ibc/clients", lab.handleIBCClients)

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	log.Printf("nico-lab listening on http://127.0.0.1%s", *listen)
	log.Printf("  home=%s chain=%s key=%s provider=%s", *home, *chainID, *key, *provider)
	log.Fatal(http.ListenAndServe(*listen, withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
