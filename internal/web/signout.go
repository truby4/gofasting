package web

import "net/http"

func (h *Handler) signoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	h.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
