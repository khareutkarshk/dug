package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Echo is a small configurable upstream used by DUG examples.
//
// Environment:
//
//	PORT          listen port (default 3000)
//	SERVICE_NAME  logical service name in responses
//	VERSION       service version label (default v1)
func main() {
	port := envOr("PORT", "3000")
	name := envOr("SERVICE_NAME", "echo")
	version := envOr("VERSION", "v1")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, r, name, version)
	})

	addr := ":" + port
	log.Printf("%s (%s) listening on %s", name, version, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeResponse(w http.ResponseWriter, r *http.Request, name, version string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Service", name)
	w.Header().Set("X-Version", version)

	switch name {
	case "openai-mock":
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "chatcmpl-dug-mock",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-mock",
			"choices": []map[string]any{{
				"index": 0,
				"message": map[string]string{
					"role":    "assistant",
					"content": "Hello from the OpenAI mock behind DUG",
				},
				"finish_reason": "stop",
			}},
			"service": name,
			"path":    r.URL.Path,
		})
	case "ollama-mock":
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model":      "llama-mock",
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"message": map[string]string{
				"role":    "assistant",
				"content": "Hello from the Ollama mock behind DUG",
			},
			"done":    true,
			"service": name,
			"path":    r.URL.Path,
		})
	default:
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "hello from " + name,
			"service": name,
			"version": version,
			"path":    r.URL.Path,
			"method":  r.Method,
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
