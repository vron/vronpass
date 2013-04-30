package vronpass

import (
	"log"
	"net/http"
	"appengine"
	"io/ioutil"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/get/", fg)
	http.HandleFunc("/set/", fs)
}

// This is for returning the data

type Entity struct {
	Value []byte
}

func fg(w http.ResponseWriter, r *http.Request) {
	// Start by checking so that the user is logged in!
	c := appengine.NewContext(r)
	c.Infof("THis is info")
	u := user.Current(c)
	if u == nil {
		c.Infof("User not signed in")
		// Generate a response that may be used to login the user
		// since this might be needed!
		str, er := user.LoginURL(c, "https://vronpass.appspot.com/")
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(str))
		return
	}
	
	// We simply store the information under the username!
	path := u.Email

	k := datastore.NewKey(c, "Entity", path, 0, nil)
	e := new(Entity)
	if err := datastore.Get(c, k, e); err!= nil {
		// This does not exist so simply return a empty response
		w.Write([]byte("{}"))
		return
	}
	// Otherwise return the string
	w.Write([]byte(e.Value))
}


func fs(w http.ResponseWriter, r *http.Request) {
	// Start by checking so that the user is logged in!
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		c.Infof("User not signed in")
		// Generate a response that may be used to login the user
		// since this might be needed!
		str, er := user.LoginURL(c, "https://vronpass.appspot.com/")
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(str))
		return
	}
	
	// We simply store the information under the username!
	path := u.Email

	// Read the body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	k := datastore.NewKey(c, "Entity", path, 0, nil)
	e := Entity{b}
	if _, err := datastore.Put(c, k, &e); err!= nil {
		// This does not exist so simply return a empty response
		log.Println("Could not put: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
