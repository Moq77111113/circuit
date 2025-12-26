package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/http/action"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	act := action.Parse(r.Form)

	switch act.Type {
	case action.ActionAdd:
		if err := h.handleAdd(act.Field); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, h.path+"?focus="+act.Field, http.StatusSeeOther)

	case action.ActionRemove:
		if err := h.handleRemove(act.Field, act.Index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, h.path+"?focus="+act.Field, http.StatusSeeOther)

	case action.ActionConfirm:
		if err := h.Apply(r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, h.path, http.StatusSeeOther)

	case action.ActionSave:
		preview, err := h.handleSave(r.Form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if preview {
			h.renderPreview(w, r)
			return
		}
		http.Redirect(w, r, h.path, http.StatusSeeOther)
	}
}
