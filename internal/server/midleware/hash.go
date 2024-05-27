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

func WithHash(h http.HandlerFunc, log *logger.Logger) http.HandlerFunc {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		log.Sugar.Infof("HASH")
		if config.Key == "" {
			h.ServeHTTP(w, r)
			return
		}
		reqHash := r.Header.Get("HashSHA256")

		if reqHash == "" {
			log.Sugar.Errorf("Empty hash")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Sugar.Errorf("Error reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		calcHash, err := hash.GetSHA256Hash(config.Key, reqData)
		if err != nil {
			log.Sugar.Errorf("Error calculating hash: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Sugar.Infof("calcHash: %s", calcHash)

		decodeHash, err := base64.URLEncoding.DecodeString(reqHash)
		log.Sugar.Infof("decodeHash: %s", string(decodeHash))
		if err != nil {
			log.Sugar.Errorf("Error decoding hash: %v", err)
		}
		if !bytes.Equal(calcHash, decodeHash) {
			log.Sugar.Errorf("Invalid hash")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.ServeHTTP(w, r)
	}
	return hashFn
}
