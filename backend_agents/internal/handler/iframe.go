package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

type IframeResponse struct {
	ID           string  `json:"id"`
	WorkspaceID  string  `json:"workspace_id"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	Icon         *string `json:"icon"`
	IframeUrl    *string `json:"iframe_url"`
	IframeScript *string `json:"iframe_script"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func iframeToResponse(i db.Iframe) IframeResponse {
	return IframeResponse{
		ID:           uuidToString(i.ID),
		WorkspaceID:  uuidToString(i.WorkspaceID),
		Title:        i.Title,
		Description:  textToPtr(i.Description),
		Icon:         textToPtr(i.Icon),
		IframeUrl:    textToPtr(i.IframeUrl),
		IframeScript: textToPtr(i.IframeScript),
		CreatedAt:    timestampToString(i.CreatedAt),
		UpdatedAt:    timestampToString(i.UpdatedAt),
	}
}

type CreateIframeRequest struct {
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	Icon         *string `json:"icon"`
	IframeUrl    *string `json:"iframe_url"`
	IframeScript *string `json:"iframe_script"`
}

type UpdateIframeRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	Icon         *string `json:"icon"`
	IframeUrl    *string `json:"iframe_url"`
	IframeScript *string `json:"iframe_script"`
}

func (h *Handler) ListIframes(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	iframes, err := h.Queries.ListIframes(r.Context(), parseUUID(workspaceID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list iframes")
		return
	}

	resp := make([]IframeResponse, len(iframes))
	for idx, iframe := range iframes {
		resp[idx] = iframeToResponse(iframe)
	}
	writeJSON(w, http.StatusOK, map[string]any{"iframes": resp, "total": len(resp)})
}

func (h *Handler) GetIframe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	workspaceID := h.resolveWorkspaceID(r)
	iframe, err := h.Queries.GetIframeInWorkspace(r.Context(), db.GetIframeInWorkspaceParams{
		ID:          parseUUID(id),
		WorkspaceID: parseUUID(workspaceID),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "iframe not found")
		return
	}
	writeJSON(w, http.StatusOK, iframeToResponse(iframe))
}

func (h *Handler) CreateIframe(w http.ResponseWriter, r *http.Request) {
	var req CreateIframeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	workspaceID := h.resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	iframe, err := h.Queries.CreateIframe(r.Context(), db.CreateIframeParams{
		WorkspaceID:  parseUUID(workspaceID),
		Title:        req.Title,
		Description:  ptrToText(req.Description),
		Icon:         ptrToText(req.Icon),
		IframeUrl:    ptrToText(req.IframeUrl),
		IframeScript: ptrToText(req.IframeScript),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create iframe")
		return
	}
	resp := iframeToResponse(iframe)
	h.publish(protocol.EventIframeCreated, workspaceID, "member", userID, map[string]any{"iframe": resp})
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateIframe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	workspaceID := h.resolveWorkspaceID(r)
	prevIframe, err := h.Queries.GetIframeInWorkspace(r.Context(), db.GetIframeInWorkspaceParams{
		ID:          parseUUID(id),
		WorkspaceID: parseUUID(workspaceID),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "iframe not found")
		return
	}
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	var req UpdateIframeRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var rawFields map[string]json.RawMessage
	json.Unmarshal(bodyBytes, &rawFields)

	params := db.UpdateIframeParams{
		ID:           prevIframe.ID,
		Description:  prevIframe.Description,
		Icon:         prevIframe.Icon,
		IframeUrl:    prevIframe.IframeUrl,
		IframeScript: prevIframe.IframeScript,
	}
	if req.Title != nil {
		params.Title = pgtype.Text{String: *req.Title, Valid: true}
	}
	if _, ok := rawFields["description"]; ok {
		if req.Description != nil {
			params.Description = pgtype.Text{String: *req.Description, Valid: true}
		} else {
			params.Description = pgtype.Text{Valid: false}
		}
	}
	if _, ok := rawFields["icon"]; ok {
		if req.Icon != nil {
			params.Icon = pgtype.Text{String: *req.Icon, Valid: true}
		} else {
			params.Icon = pgtype.Text{Valid: false}
		}
	}
	if _, ok := rawFields["iframe_url"]; ok {
		if req.IframeUrl != nil {
			params.IframeUrl = pgtype.Text{String: *req.IframeUrl, Valid: true}
		} else {
			params.IframeUrl = pgtype.Text{Valid: false}
		}
	}
	if _, ok := rawFields["iframe_script"]; ok {
		if req.IframeScript != nil {
			params.IframeScript = pgtype.Text{String: *req.IframeScript, Valid: true}
		} else {
			params.IframeScript = pgtype.Text{Valid: false}
		}
	}

	iframe, err := h.Queries.UpdateIframe(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update iframe")
		return
	}
	resp := iframeToResponse(iframe)
	h.publish(protocol.EventIframeUpdated, workspaceID, "member", userID, map[string]any{"iframe": resp})
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteIframe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	workspaceID := h.resolveWorkspaceID(r)
	if _, err := h.Queries.GetIframeInWorkspace(r.Context(), db.GetIframeInWorkspaceParams{
		ID:          parseUUID(id),
		WorkspaceID: parseUUID(workspaceID),
	}); err != nil {
		writeError(w, http.StatusNotFound, "iframe not found")
		return
	}
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	if err := h.Queries.DeleteIframe(r.Context(), parseUUID(id)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete iframe")
		return
	}
	h.publish(protocol.EventIframeDeleted, workspaceID, "member", userID, map[string]any{"iframe_id": id})
	w.WriteHeader(http.StatusNoContent)
}
