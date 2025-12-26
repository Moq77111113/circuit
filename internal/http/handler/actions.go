package handler

import "github.com/moq77111113/circuit/internal/http/form"

func (h *Handler) handleSave(formData map[string][]string) (bool, error) {
	if !h.store.AutoApply() {
		return true, nil
	}

	var err error
	h.store.WithLock(func() {
		err = form.Apply(h.cfg, h.schema, formData)
	})

	if err != nil {
		return false, err
	}

	return false, h.writeConfig()
}

func (h *Handler) handleAdd(fieldName string) error {
	var err error
	h.store.WithLock(func() {
		err = form.AddSliceItemNode(h.cfg, h.schema.Nodes, fieldName)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}

func (h *Handler) handleRemove(fieldName string, index int) error {
	var err error
	h.store.WithLock(func() {
		err = form.RemoveSliceItemNode(h.cfg, h.schema.Nodes, fieldName, index)
	})

	if err != nil {
		return err
	}

	return h.writeConfig()
}
