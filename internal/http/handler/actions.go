package handler

import "github.com/moq77111113/circuit/internal/http/form"

func (h *Handler) handleSave(formData map[string][]string) error {
	var err error
	h.loader.WithLock(func() {
		err = form.Apply(h.cfg, h.schema, formData)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}

func (h *Handler) handleAdd(fieldName string) error {
	var err error
	h.loader.WithLock(func() {
		err = form.AddSliceItem(h.cfg, h.schema, fieldName)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}

func (h *Handler) handleRemove(fieldName string, index int) error {
	var err error
	h.loader.WithLock(func() {
		err = form.RemoveSliceItem(h.cfg, h.schema, fieldName, index)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}
