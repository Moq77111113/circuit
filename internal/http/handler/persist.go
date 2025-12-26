package handler

import (
	"github.com/moq77111113/circuit/internal/sync"
)

func (h *Handler) writeConfig() error {
	h.store.MarkFormSubmit()

	if h.store.AutoSave() {
		if err := h.store.Save(); err != nil {
			return err
		}
	}

	h.store.EmitChange(sync.SourceFormSubmit)

	return nil
}
