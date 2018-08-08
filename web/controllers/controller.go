package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/servntire/car-ownership/blockchain"
)

type Application struct {
	Fabric *blockchain.FabricSetup
}

func add(x, y int) int {
	return x + y
}

func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	lp := filepath.Join("web", "templates", "layout.html")
	tp := filepath.Join("web", "templates", templateName)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(tp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	funcs := template.FuncMap{"add": add}
	resultTemplate := template.Must(template.New(templateName).Funcs(funcs).ParseFiles(tp, lp))
	/*if err != nil {
		// Log the detailed error
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}*/
	if err := resultTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
