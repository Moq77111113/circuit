package handler

import (
	"net/url"

	"github.com/moq77111113/circuit/internal/http/form"
)

// Apply manually applies form data to the config.
// Used in preview mode (autoApply=false) to confirm changes after user review.
// Respects autoSave setting: saves to disk if enabled.
func (h *Handler) Apply(formData url.Values) error {
	var err error
	h.store.WithLock(func() {
		err = form.Apply(h.cfg, h.schema, formData)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}

// Save manually saves the current config to disk.
// Uses custom saveFunc if provided, otherwise writes YAML.
func (h *Handler) Save() error {
	return h.store.Save()
}
