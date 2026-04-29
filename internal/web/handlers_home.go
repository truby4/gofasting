package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/truby4/gofasting/internal/fasts"
)

type fastStartForm struct {
	Goal int `form:"goal_hours" validate:"required,gte=5,lte=168"`

	Validator
}

var validate = validator.New()

func (f *fastStartForm) validateForm() error {
	err := validate.Struct(f)
	if err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			verrs := make(map[string]string)
			for _, e := range validateErrs {
				verrs[strings.ToLower(e.Field())] = "goal must be over 5 and less than 168"
			}
			f.AddFieldErrors(verrs)
		}
		return err
	}
	return nil
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	data := h.newTemplateData(r)
	data.Form = fastStartForm{
		Goal: 8,
	}

	fast, err := h.fasts.GetActiveFast(userID)
	if errors.Is(err, fasts.ErrNoRecord) {
		err = nil
	} else if err != nil {
		h.serverError(w, r, err)
		return
	} else {
		data.Fast = &fast
	}

	fasts, err := h.fasts.GetHistory(userID)
	if err != nil {
		h.serverError(w, r, err)
		return
	} else {
		h.logger.Debug("Fasts download successful?")
		data.Fasts = fasts
	}

	h.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (h *Handler) fastStartPost(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	var form fastStartForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form.validateForm()

	_, err = h.fasts.Start(form.Goal*3600, userID)
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	h.sessionManager.Put(r.Context(), "flash", "You've started a new fast!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) fastEndPost(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	err := h.fasts.End(userID)
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	h.sessionManager.Put(r.Context(), "flash", "Well done, another fast down!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
