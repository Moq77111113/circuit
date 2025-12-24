package handler

import (
	"os"

	"github.com/moq77111113/circuit/internal/sync"
	"github.com/moq77111113/circuit/internal/yaml"
)

func (h *Handler) writeConfig() error {
	var data []byte
	var err error
	h.store.WithLock(func() {
		data, err = yaml.Encode(h.cfg)
	})

	if err != nil {
		return err
	}

	err = os.WriteFile(h.path, data, 0644)
	if err != nil {
		return err
	}

	h.store.EmitChange(sync.SourceFormSubmit)

	return nil
}
