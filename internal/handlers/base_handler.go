package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-backend-starter-template/internal/types/selectors"
)

type indexPageModel struct {
	Login bool
}

func BaseHandler(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.BaseHTML, indexPageModel{Login: true})
		if err != nil {
			slog.Error("Failed to execute template", "error", err, "template", selectors.IndexPage.BaseHTML)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func LoginHandler(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.LoginComponent, indexPageModel{Login: true})
		if err != nil {
			slog.Error("Failed to execute template", "error", err, "template", selectors.IndexPage.LoginComponent)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func SignUpHandler(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, selectors.IndexPage.SignUpComponent, indexPageModel{Login: false})
		if err != nil {
			slog.Error("Failed to execute template", "error", err, "template", selectors.IndexPage.SignUpComponent)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
