package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"io"
	"net/http"
	"strings"
)

func UpdateHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// разрешаем только POST-запросы
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		urlPath := strings.Split(r.URL.Path, "/")
		urlPathLen := len(urlPath)

		//if urlPathLen != 5 {
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		if urlPathLen < 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mtype := urlPath[2]
		mtitle := urlPath[3]

		//check type in path
		if mtype == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//check title in path
		if mtitle == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if urlPathLen < 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mval := urlPath[4]

		if mval == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := app.Save(mtype, mtitle, mval); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
