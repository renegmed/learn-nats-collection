package web

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	log.Println("...home")

	files := []string{
		"ui/html/home.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	//ts, err := template.ParseFiles(files...)
	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getAPIContent(url string, templateData interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(bodyBytes, templateData)
	return nil
}

func (app *Application) static(dir string) http.Handler {
	dirCleaned := filepath.Clean(dir)
	return http.StripPrefix("/static/", http.FileServer(http.Dir(dirCleaned)))
}
