package web

import (
	"fmt"
	"net/http"

	"github.com/servntire/car-ownership/web/controllers"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	http.HandleFunc("/query.html", app.QueryHandler)
	http.HandleFunc("/create.html", app.CreateHandler)
	http.HandleFunc("/update.html", app.UpdateHandler)
	http.HandleFunc("/history.html", app.HistoryHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
