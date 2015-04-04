package notifyssh

import (
	"encoding/json"
	"fmt"
	//"github.com/alexjlockwood/gcm"
	//"appengine"
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseGlob("app/templates/*.html"))

func deviceKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "deivce", "public", 0, nil)
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/device", device)
	http.HandleFunc("/notify", notify)
	http.HandleFunc("/sign", handleSign)
}

func notify(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	if r.Method != "POST" {
		return
	}

	q := datastore.NewQuery("device").Ancestor(deviceKey(c))
	var dd []*Device
	if _, err := q.GetAll(c, &dd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof(c, "Rendering %d devices", len(dd))

	ids := make([]string, len(dd))
	for i, v := range dd {
		if len(v.Id) == 0 {
			continue
		}
		ids[i] = v.Id
	}

	data := struct {
		Ids []string `json:"registration_ids"`
	}{
		Ids: ids,
	}

	res, _ := json.Marshal(data)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(string(res))
	if err != nil {
		fmt.Println(err)
		return
	}

	client := urlfetch.Client(c)
	req, _ := http.NewRequest("POST", "https://android.googleapis.com/gcm/send", buf)
	req.Header.Add("Authorization", "key=AIzaSyCCDovuMzTMspo8CgahN1Wew3Y_DZHQWdU")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Infof(c, "%s", resp.Body)
	w.Write([]byte("OK"))
}

func handleSign(w http.ResponseWriter, r *http.Request) {
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

	c := appengine.NewContext(r)

	if r.Method == "GET" {
		q := datastore.NewQuery("device").Ancestor(deviceKey(c))
		var dd []*Device
		if _, err := q.GetAll(c, &dd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Devices []*Device
		}{
			Devices: dd,
		}

		log.Infof(c, "Rendering %d devices", len(dd))

		if err := tpl.ExecuteTemplate(w, "device.html", data); err != nil {
			log.Errorf(c, "%v", err)
		}

		return
	}

	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 404)
		return
	}

	d := &Device{
		Id:       r.FormValue("id"),
		Endpoint: r.FormValue("endpoint"),
	}

	log.Infof(c, "devices id %s", r.FormValue("id"))

	key := datastore.NewIncompleteKey(c, "device", deviceKey(c))
	if _, err := datastore.Put(c, key, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
