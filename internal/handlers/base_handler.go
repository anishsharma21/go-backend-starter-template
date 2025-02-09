package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func BaseHandler(dbConn *pgx.Conn, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			slog.Error("Failed to execute template", "error", err, "template", "index.html")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
