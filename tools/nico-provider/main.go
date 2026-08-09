// nico-provider is a lightweight Jackal storage provider for private jackal-nico-1 testing.
// It speaks the Sequoia-compatible /v2/upload API used by `canined tx storage post`,
// without depending on Sequoia's cosmos-sdk 0.45 wallet stack.
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

type acceptedUploadResponse struct {
	JobID string `json:"job_id"`
}

type statusResponse struct {
	Status  string `json:"status"`
	Address string `json:"address"`
	ChainID string `json:"chain_id"`
	Files   int    `json:"files"`
}

type versionResponse struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	ChainID string `json:"chain-id"`
}

type chainFile struct {
	Merkle   string `json:"merkle"`
	Owner    string `json:"owner"`
	Start    string `json:"start"`
	FileSize string `json:"file_size"`
}

type fileQueryResponse struct {
	File chainFile `json:"file"`
}

func main() {
	addr := flag.String("address", "jkl-provider-local", "provider bech32 address (for / status)")
	listen := flag.String("listen", ":3333", "HTTP listen address")
	dataDir := flag.String("data", "", "directory to store uploaded files (default ~/.nico-provider/data)")
	rest := flag.String("rest", "http://127.0.0.1:1317", "canined REST API base URL")
	chainID := flag.String("chain-id", "jackal-nico-1", "chain id reported by /version")
	skipChain := flag.Bool("skip-chain-check", false, "accept uploads without verifying MsgPostFile on-chain (dev only)")
	flag.Parse()

	if *dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		*dataDir = filepath.Join(home, ".nico-provider", "data")
	}
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		log.Fatal(err)
	}

	s := &server{
		addr:      *addr,
		dataDir:   *dataDir,
		rest:      strings.TrimRight(*rest, "/"),
		chainID:   *chainID,
		skipChain: *skipChain,
		jobs:      sync.Map{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleStatus)
	mux.HandleFunc("/version", s.handleVersion)
	mux.HandleFunc("/v2/upload", s.handleUploadV2)
	mux.HandleFunc("/upload", s.handleUploadV2)
	mux.HandleFunc("/download/", s.handleDownload)
	mux.HandleFunc("/list", s.handleList)

	log.Printf("nico-provider listening on %s (data=%s rest=%s chain=%s)", *listen, *dataDir, *rest, *chainID)
	log.Fatal(http.ListenAndServe(*listen, withCORS(mux)))
}

type server struct {
	addr      string
	dataDir   string
	rest      string
	chainID   string
	skipChain bool
	jobs      sync.Map
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

func writeErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func (s *server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	files, _ := os.ReadDir(s.dataDir)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(statusResponse{
		Status:  "online",
		Address: s.addr,
		ChainID: s.chainID,
		Files:   len(files),
	})
}

func (s *server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(versionResponse{
		Version: "nico-provider-dev",
		Build:   "local",
		ChainID: s.chainID,
	})
}

func (s *server) handleList(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"files": files})
}

func (s *server) handleDownload(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/download/")
	name = filepath.Base(name)
	if name == "" || name == "." {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("missing merkle"))
		return
	}
	path := filepath.Join(s.dataDir, name)
	f, err := os.Open(path)
	if err != nil {
		writeErr(w, http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	_, _ = io.Copy(w, f)
}

func (s *server) handleUploadV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("cannot parse form: %w", err))
		return
	}

	sender := r.FormValue("sender")
	merkleHex := r.FormValue("merkle")
	startStr := r.FormValue("start")
	if sender == "" || merkleHex == "" || startStr == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("sender, merkle, and start are required"))
		return
	}
	merkle, err := hex.DecodeString(merkleHex)
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("cannot parse merkle: %w", err))
		return
	}
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("cannot parse start: %w", err))
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("file field required: %w", err))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("cannot read file: %w", err))
		return
	}

	h := sha256.New()
	_, _ = h.Write(merkle)
	_, _ = h.Write([]byte(sender))
	_, _ = h.Write([]byte(strconv.FormatInt(start, 10)))
	jobID := hex.EncodeToString(h.Sum(nil))

	if !s.skipChain {
		ok, size, err := s.verifyOnChain(merkleHex, sender, start)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, fmt.Errorf("chain verify failed: %w", err))
			return
		}
		if !ok {
			writeErr(w, http.StatusInternalServerError, fmt.Errorf("file not found on chain"))
			return
		}
		if size > 0 && int64(len(data)) != size {
			writeErr(w, http.StatusInternalServerError, fmt.Errorf("size mismatch got=%d want=%d", len(data), size))
			return
		}
	}

	path := filepath.Join(s.dataDir, merkleHex)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		writeErr(w, http.StatusInternalServerError, fmt.Errorf("write failed: %w", err))
		return
	}
	s.jobs.Store(jobID, time.Now())
	log.Printf("stored %s (%d bytes) owner=%s start=%d", merkleHex, len(data), sender, start)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(acceptedUploadResponse{JobID: jobID})
}

func (s *server) verifyOnChain(merkleHex, owner string, start int64) (bool, int64, error) {
	merkleBytes, err := hex.DecodeString(merkleHex)
	if err != nil {
		return false, 0, err
	}
	// gRPC-gateway expects URL-safe base64 for bytes path params.
	merkleB64 := base64.URLEncoding.EncodeToString(merkleBytes)
	u := fmt.Sprintf("%s/jackal/canine-chain/storage/files/%s/%s/%d",
		s.rest, merkleB64, url.PathEscape(owner), start)

	cli := &http.Client{Timeout: 15 * time.Second}
	res, err := cli.Get(u)
	if err != nil {
		return false, 0, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return false, 0, fmt.Errorf("status %d: %s", res.StatusCode, string(body))
	}

	var out fileQueryResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return false, 0, err
	}
	size, _ := strconv.ParseInt(out.File.FileSize, 10, 64)
	return true, size, nil
}
