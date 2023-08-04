package timelines

import (
	"encoding/json"
	"net/http"
	"context"
	"strconv"
)

type ContextKey string

const (
    ContextOnlyMedia ContextKey = "only_media"
    ContextMaxID     ContextKey = "max_id"
    ContextSinceID   ContextKey = "since_id"
    ContextLimit     ContextKey = "limit"

	ContextAr ContextKey = "ar"
)

func (h *handler) PublicTimelines(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	onlyMediastr := r.URL.Query().Get("only_media")
    maxIDstr := r.URL.Query().Get("max_id")
    sinceIDstr := r.URL.Query().Get("since_id")
    limitstr := r.URL.Query().Get("limit")

	maxID, _ := strconv.ParseInt(maxIDstr, 10, 64)
    sinceID, _ := strconv.ParseInt(sinceIDstr, 10, 64)
    limit, _ := strconv.ParseInt(limitstr, 10, 64)
	onlyMedia, _ := strconv.ParseBool(onlyMediastr)
	

    ctx = context.WithValue(ctx, ContextOnlyMedia, onlyMedia)
    ctx = context.WithValue(ctx, ContextMaxID, maxID)
    ctx = context.WithValue(ctx, ContextSinceID, sinceID)
    ctx = context.WithValue(ctx, ContextLimit, limit)

	ctx = context.WithValue(ctx, ContextAr, h.ar)

	timelines, err := h.sr.GetTimelines(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timelines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
