package handler

import (
	"net/http"
	"net/url"

	"github.com/moq77111113/circuit/internal/actions"
)

func (h *Handler) executeAction(w http.ResponseWriter, r *http.Request, actionName string) {
	if h.readOnly {
		http.Error(w, "Actions not allowed in read-only mode", http.StatusForbidden)
		return
	}

	var found *actions.Def
	for i := range h.actions {
		if h.actions[i].Name == actionName {
			found = &h.actions[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Action not found", http.StatusNotFound)
		return
	}

	if err := actions.Execute(r.Context(), *found); err != nil {
		errMsg := url.QueryEscape(err.Error())
		http.Redirect(w, r, "/?error="+errMsg, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
