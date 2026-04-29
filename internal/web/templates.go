package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/justinas/nosurf"
	"github.com/truby4/gofasting/internal/fasts"
)

type templateData struct {
	Flash           string
	Form            any
	IsAuthenticated bool
	Fast            *fasts.Fast
	Fasts           []fasts.Fast
	CSRFToken       string
}

func (h *Handler) newTemplateData(r *http.Request) templateData {
	return templateData{
		Flash:           h.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: h.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		h.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		h.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 15:04")
}

func humanDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return humanDate(*t)
}

func timeAgo(t time.Time) string {
	return humanize.Time(t)
}

func goalText(goalSeconds int) string {
	hours := goalSeconds / 3600
	if hours == 1 {
		return "1 hour"
	}
	return fmt.Sprintf("%d hours", hours)
}

func expectedEnd(start time.Time, goalSeconds int) string {
	end := start.Add(time.Duration(goalSeconds) * time.Second)
	return humanDate(end)
}
func completedDurationText(f fasts.Fast) string {
	if f.EndTime == nil {
		return "In progress"
	}
	return durationText(int(f.EndTime.Sub(f.StartTime).Seconds()))
}

func fastDurationText(f fasts.Fast) string {
	end := time.Now()
	if f.EndTime != nil {
		end = *f.EndTime
	}
	return durationText(int(end.Sub(f.StartTime).Seconds()))
}

func fastRangeText(f fasts.Fast) string {
	if f.EndTime == nil {
		return fmt.Sprintf("%s → now", humanDate(f.StartTime))
	}
	return fmt.Sprintf("%s → %s", humanDate(f.StartTime), humanDate(*f.EndTime))
}

var functions = template.FuncMap{
	"humanDate":             humanDate,
	"humanDatePtr":          humanDatePtr,
	"timeAgo":               timeAgo,
	"goalText":              goalText,
	"expectedEnd":           expectedEnd,
	"fastDurationText":      fastDurationText,
	"fastRangeText":         fastRangeText,
	"completedDurationText": completedDurationText,
	"durationText":          durationText,
}

func durationText(seconds int) string {
	if seconds <= 0 {
		return "0m"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	switch {
	case hours == 0:
		return fmt.Sprintf("%dm", minutes)
	case minutes == 0:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}
