package notifyssh

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseGlob("app/templates/*.html"))

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	data := struct {
		Title string
	}{
		Title: "NotifySSH",
	}

	if err := tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Errorf(c, "%v", err)
	}
}
