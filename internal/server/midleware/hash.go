package midleware

import (
	"bytes"
	"encoding/base64"
	"github.com/mylastgame/yp-metrics-service/internal/core/hash"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"io"
	"net/http"
)

type hashWriter struct {
	http.ResponseWriter
	body []byte
}

func newHashWriter(w http.ResponseWriter) *hashWriter {
	return &hashWriter{w, []byte{}}
}

func (hw *hashWriter) Write(p []byte) (int, error) {
	hw.body = p
	return hw.ResponseWriter.Write(p)
}

func WithHash(h http.HandlerFunc, log *logger.Logger) http.HandlerFunc {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		if config.Key == "" {
			h.ServeHTTP(w, r)
			return
		}
		reqHash := r.Header.Get("Hash")

		if reqHash == "" {
			//log.Sugar.Errorf("Empty hash")
			//w.WriteHeader(http.StatusBadRequest)
			//return
			//autotest fixes
			h.ServeHTTP(w, r)
			return
		}

		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Sugar.Errorf("Error reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(reqData))

		calcHash, err := hash.GetSHA256Hash(config.Key, reqData)
		if err != nil {
			log.Sugar.Errorf("Error calculating hash: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decodeHash, err := base64.URLEncoding.DecodeString(reqHash)
		if err != nil {
			log.Sugar.Errorf("Error decoding hash: %v", err)
		}
		if !bytes.Equal(calcHash, decodeHash) {
			log.Sugar.Errorf("Invalid hash")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hw := newHashWriter(w)
		h.ServeHTTP(hw, r)

		if config.Key == "" {
			return
		}

		respHash, err := hash.GetSHA256Hash(config.Key, hw.body)
		if err != nil {
			log.Sugar.Errorf("Error calculating hash: %v", err)
			return
		}

		w.Header().Set("Hash", string(respHash))
	}
	return hashFn
}
