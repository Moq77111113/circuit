package circuit

import (
	"net/http"
	"net/url"

	"github.com/moq77111113/circuit/internal/http/handler"
)

// Handler serves the Circuit UI and provides manual control methods.
type Handler struct {
	h *handler.Handler
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}

// Apply manually applies form data to the config.
// Used in preview mode (WithAutoApply(false)) to confirm changes after review.
// Respects WithAutoSave setting: saves to disk if enabled.
func (h *Handler) Apply(formData url.Values) error {
	return h.h.Apply(formData)
}

// Save manually saves the current config to disk.
// Uses custom SaveFunc if provided via WithSaveFunc, otherwise writes YAML.
func (h *Handler) Save() error {
	return h.h.Save()
}
