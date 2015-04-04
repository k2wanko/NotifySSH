package notifyssh

import (
	//"github.com/alexjlockwood/gcm"
	"google.golang.org/appengine"
	//"google.golang.org/appengine/datastore"
	"fmt"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseGlob("app/templates/*.html"))

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/device", device)
	http.HandleFunc("/sign", handleSign)
}

func handleSign(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		url, _ := user.LoginURL(c, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}

	url, _ := user.LogoutURL(c, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func device(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 404)
		return
	}

	w.Write([]byte("OK"))
}

func handler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	data := struct {
		Title string
		Name  string
	}{
		Title: "NotifySSH",
		Name:  "",
	}

	if u != nil {
		data.Name = u.String()
	}

	if err := tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Errorf(c, "%v", err)
	}
}
