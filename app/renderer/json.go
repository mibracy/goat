package renderer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

func PrettyJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		render.JSON(w, r, err) // fallback to default render
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(jsonBytes)
}
