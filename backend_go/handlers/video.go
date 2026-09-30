package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const pythonStreamURL = "http://localhost:8001/api/video"

// RegisterVideoHandler registers the raw video stream proxy route on Chi
func RegisterVideoHandler(r chi.Router) {
	r.Get("/api/video", func(w http.ResponseWriter, r *http.Request) {
		// 1. Connect to FastAPI's video stream
		resp, err := http.Get(pythonStreamURL)
		if err != nil {
			http.Error(w, "Failed to connect to OCR video stream", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 2. Forward streaming headers (multipart/x-mixed-replace)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// 3. Pipe the frame bytes straight to the client
		_, _ = io.Copy(w, resp.Body)
	})
}
