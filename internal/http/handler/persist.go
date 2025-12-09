package handler

import (
	"os"

	"github.com/moq77111113/circuit/internal/yaml"
)

func (h *Handler) writeConfig() error {
	var data []byte
	var err error
	h.loader.WithLock(func() {
		data, err = yaml.Encode(h.cfg)
	})

	if err != nil {
		return err
	}

	return os.WriteFile(h.path, data, 0644)
}
