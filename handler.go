package circuit

import (
	"net/http"
	"net/url"

	"github.com/moq77111113/circuit/internal/http/handler"
)

// Handler serves the Circuit UI and provides manual control methods.
//
// The handler implements http.Handler and can be mounted on any mux.
//
// Manual control methods are available for advanced workflows:
//   - Apply(formData) - manually apply form changes (preview mode with WithAutoApply(false))
//   - Save() - manually persist to disk (manual save with WithAutoSave(false))
//
// Example (preview mode):
//
//	h, _ := circuit.From(&cfg, circuit.WithAutoApply(false))
//
//	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
//	    if r.Method == "POST" && r.FormValue("confirm") == "yes" {
//	        r.ParseForm()
//	        if err := h.Apply(r.Form); err != nil {
//	            http.Error(w, err.Error(), 500)
//	            return
//	        }
//	    }
//	    h.ServeHTTP(w, r)
//	})
type Handler struct {
	h *handler.Handler
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}

// Apply manually applies form data to the in-memory config.
//
// Used in preview mode (WithAutoApply(false)) to confirm changes after user review.
// The typical workflow:
//  1. User submits form (POST)
//  2. Circuit renders preview with submitted values (doesn't modify memory yet)
//  3. User reviews and confirms
//  4. Your code calls Apply(formData) to commit changes
//
// Respects WithAutoSave setting: if enabled, changes are saved to disk after apply.
// If WithAutoSave(false), call Save() separately to persist.
func (h *Handler) Apply(formData url.Values) error {
	return h.h.Apply(formData)
}

// Save manually persists the current in-memory config to disk.
//
// Used in manual save mode (WithAutoSave(false)) to control when persistence happens.
// Common use cases:
//   - Batch multiple changes before saving
//   - Add validation or approval workflows before persisting
//   - Trigger saves on external events (timers, signals)
//
// Uses custom SaveFunc if provided via WithSaveFunc, otherwise writes YAML to the path
// specified in WithPath.
func (h *Handler) Save() error {
	return h.h.Save()
}
