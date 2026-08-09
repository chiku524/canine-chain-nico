package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Lab struct {
	Home     string
	ChainID  string
	RPC      string
	REST     string
	Provider string
	Canined  string
	Key      string
	Keyring  string
	Client   *http.Client
}

type apiError struct {
	Error string `json:"error"`
}

type apiOK struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Raw     string          `json:"raw,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, apiError{Error: err.Error()})
}

func (l *Lab) getJSON(url string) (json.RawMessage, error) {
	res, err := l.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.RawMessage(body), nil
}

func (l *Lab) keyAddr() (string, error) {
	out, err := l.runCaninedNoNode("keys", "show", l.Key, "-a", "--keyring-backend", l.Keyring)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (l *Lab) runCaninedNoNode(args ...string) (string, error) {
	base := []string{
		"--home", l.Home,
		"--chain-id", l.ChainID,
	}
	return l.execCanined(append(base, args...)...)
}

func (l *Lab) runCanined(args ...string) (string, error) {
	base := []string{
		"--home", l.Home,
		"--chain-id", l.ChainID,
		"--node", "tcp://" + strings.TrimPrefix(strings.TrimPrefix(l.RPC, "http://"), "https://"),
	}
	return l.execCanined(append(base, args...)...)
}

func (l *Lab) execCanined(args ...string) (string, error) {
	cmd := exec.Command(l.Canined, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return stdout.String(), fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}

func (l *Lab) runTx(args ...string) (string, error) {
	txArgs := append([]string{}, args...)
	txArgs = append(txArgs,
		"--from", l.Key,
		"--keyring-backend", l.Keyring,
		"--fees", "5000ujkl",
		"--gas", "auto",
		"--gas-adjustment", "1.5",
		"-y",
		"--output", "json",
	)
	return l.runCanined(txArgs...)
}

func (l *Lab) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{
		"ok":       true,
		"chain_id": l.ChainID,
		"home":     l.Home,
		"provider": l.Provider,
		"key":      l.Key,
	})
}

func (l *Lab) handleStatus(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.RPC + "/status")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleBalance(w http.ResponseWriter, _ *http.Request) {
	addr, err := l.keyAddr()
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	data, err := l.getJSON(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", l.REST, addr))
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: addr, Data: data})
}

func (l *Lab) handleProviders(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine-chain/storage/providers")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleStorageParams(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine-chain/storage/params")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleStoragePayment(w http.ResponseWriter, _ *http.Request) {
	addr, err := l.keyAddr()
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	data, err := l.getJSON(fmt.Sprintf("%s/jackal/canine-chain/storage/payment_info/%s", l.REST, addr))
	if err != nil {
		// not found is fine before buy
		writeJSON(w, 200, apiOK{OK: true, Message: addr, Raw: err.Error()})
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: addr, Data: data})
}

func (l *Lab) handleBuyStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, 405, fmt.Errorf("POST required"))
		return
	}
	addr, err := l.keyAddr()
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	days := r.URL.Query().Get("days")
	if days == "" {
		days = "30"
	}
	bytesAmt := r.URL.Query().Get("bytes")
	if bytesAmt == "" {
		bytesAmt = "1000000000" // 1 GB
	}
	out, err := l.runTx("tx", "storage", "buy-storage", addr, days, bytesAmt, "ujkl")
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: "buy-storage submitted", Raw: strings.TrimSpace(out)})
}

func (l *Lab) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, 405, fmt.Errorf("POST required"))
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeErr(w, 400, err)
		return
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		// fallback: text body field
		text := r.FormValue("text")
		if text == "" {
			writeErr(w, 400, fmt.Errorf("file or text required"))
			return
		}
		tmp, err := os.CreateTemp("", "nico-lab-*.txt")
		if err != nil {
			writeErr(w, 500, err)
			return
		}
		defer os.Remove(tmp.Name())
		if _, err := tmp.WriteString(text); err != nil {
			tmp.Close()
			writeErr(w, 500, err)
			return
		}
		tmp.Close()
		l.doUpload(w, tmp.Name())
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp("", "nico-lab-*"+filepath.Ext(hdr.Filename))
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		writeErr(w, 500, err)
		return
	}
	tmp.Close()
	l.doUpload(w, tmp.Name())
}

func (l *Lab) doUpload(w http.ResponseWriter, path string) {
	statusBody, err := l.getJSON(l.RPC + "/status")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	var st struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}
	_ = json.Unmarshal(statusBody, &st)
	height, _ := strconv.ParseInt(st.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if height == 0 {
		height = time.Now().Unix() % 1_000_000
	}
	expire := strconv.FormatInt(height+100000, 10)

	out, err := l.runTx(
		"tx", "storage", "post", path, expire,
		"--dest", l.Provider,
		"--max_proofs", "1",
	)
	if err != nil {
		writeErr(w, 500, fmt.Errorf("%w\n%s", err, out))
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: "posted + uploaded", Raw: strings.TrimSpace(out)})
}

func (l *Lab) handleChainFiles(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine-chain/storage/files")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleProviderStatus(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.Provider + "/")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleProviderFiles(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.Provider + "/list")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleProviderDownload(w http.ResponseWriter, r *http.Request) {
	merkle := strings.TrimPrefix(r.URL.Path, "/api/provider/download/")
	if merkle == "" {
		writeErr(w, 400, fmt.Errorf("missing merkle"))
		return
	}
	res, err := l.Client.Get(l.Provider + "/download/" + merkle)
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	defer res.Body.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+merkle+"\"")
	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)
}

func (l *Lab) handleFiletreeParams(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine/filetree/params")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleFiletreeProvision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, 405, fmt.Errorf("POST required"))
		return
	}
	out, err := l.runTx("tx", "filetree", "provision")
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: "filetree provision submitted", Raw: strings.TrimSpace(out)})
}

func (l *Lab) handleRNSParams(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine-chain/rns/params")
	if err != nil {
		// alternate path prefix used by some modules
		data, err = l.getJSON(l.REST + "/jackal/canine/rns/params")
		if err != nil {
			writeErr(w, 502, err)
			return
		}
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleOracleFeeds(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/jackal/canine-chain/oracle/feeds")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleOracleCreateFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, 405, fmt.Errorf("POST required"))
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "jklprice"
	}
	out, err := l.runTx("tx", "oracle", "create-feed", name)
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Message: "create-feed submitted", Raw: strings.TrimSpace(out)})
}

func (l *Lab) handleJklmintParams(w http.ResponseWriter, _ *http.Request) {
	out, err := l.runCanined("query", "jklmint", "params", "--output", "json")
	if err != nil {
		writeErr(w, 500, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Raw: strings.TrimSpace(out), Data: json.RawMessage(strings.TrimSpace(out))})
}

func (l *Lab) handleWasmCodes(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/cosmwasm/wasm/v1/code")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}

func (l *Lab) handleIBCClients(w http.ResponseWriter, _ *http.Request) {
	data, err := l.getJSON(l.REST + "/ibc/core/client/v1/client_states")
	if err != nil {
		writeErr(w, 502, err)
		return
	}
	writeJSON(w, 200, apiOK{OK: true, Data: data})
}
