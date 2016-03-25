package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// Paste retrieves the nth recently copied object given the user's id.
func (a *API) Paste(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if id == "" {
		BadRequest.ServeHTTP(w, r)
		return
	}

	index, err := strconv.ParseUint(params.ByName("index"), 10, 64)
	if err != nil {
		BadRequest.ServeHTTP(w, r)
		return
	}

	obj := a.retrieve(id, index)
	if obj == "" {
		NotFound.ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, "https://pbcp.s3-us-west-2.amazonaws.com/"+obj, http.StatusFound)
}

// Copy adds the uploaded object for a user
func (a *API) Copy(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if id == "" {
		BadRequest.ServeHTTP(w, r)
		return
	}

	size := r.ContentLength
	if size <= 0 {
		BadRequest.ServeHTTP(w, r)
		return
	}

	// Check if size is larger than limit
	if size > 10000000 {
		BadRequest.ServeHTTP(w, r)
		log.Printf("File too large! %d\n", size)
		return
	}

	obj := a.id()

	err := upload(obj, r.Body)
	if err != nil {
		InternalError.ServeHTTP(w, r)
		log.Printf("Couldn't upload file: %s\n", err)
		return
	}

	err = a.add(id, obj)
	if err != nil {
		InternalError.ServeHTTP(w, r)
		log.Printf("Couldn't upload file: %s\n", err)
		return
	}

	a.respond(obj, w, r)
}

// Register outputs a unique user ID
func (a *API) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(a.id()))
}

// Debug
func (a *API) Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users := map[string][]string{}
	a.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user := parse(v)
			users[string(k)] = user
		}
		return nil
	})
	a.respond(users, w, r)
}

// retrieve retrieves the id of the user's nth object
func (a *API) retrieve(id string, n uint64) string {
	var obj string = ""

	err := a.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		if b == nil {
			return errors.New("Users bucket does not exist!")
		}

		data := b.Get([]byte(id))
		if data == nil {
			return errors.New("User does not exist")
		}

		user := parse(data)
		if user == nil {
			return errors.New("Invalid user?")
		}

		obj = user[len(user)-1-int(n)]
		return nil
	})
	if err != nil {
		log.Printf("Retrieve: %s\n", err)
	}

	return obj
}

// list a user's objects
func (a *API) list(id string) {
	a.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		if b == nil {
			return errors.New("Users bucket does not exist!")
		}

		user := b.Bucket([]byte(id))
		if user == nil {
			return errors.New("User does not exist")
		}

		c := user.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Printf("===> key=%s, value=%s\n", k, v)
		}

		return nil
	})
}

// parse parses user JSON data
func parse(data []byte) []string {
	var user []string

	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil
	}

	return user
}

// unparse converts a user into JSON data
func unparse(user []string) ([]byte, error) {
	return json.Marshal(user)
}

// add adds an object id to a user's bucket
func (a *API) add(id string, obj string) error {
	return a.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		if b == nil {
			return errors.New("User bucket does not exist!")
		}

		var user []string

		data := b.Get([]byte(id))
		if data == nil {
			user = []string{}
		} else {
			user = parse(data)
			if user == nil {
				return errors.New("Invalid user?")
			}
		}

		// Limit to 10 objects
		user = append(user, obj)
		for len(user) > 10 {
			var excess string
			excess, user = user[0], user[1:]
			err := delete(excess)
			if err != nil {
				return err
			}
		}

		// Store into DB
		json, err := unparse(user)
		if err != nil {
			return err
		}
		b.Put([]byte(id), json)

		return nil
	})
}

// id generates service-wide unique IDs (used for both users and files)
func (a *API) id() string {
	byt := make([]byte, 17)
	_, err := rand.Read(byt)
	if err != nil {
		return ""
	}

	var index uint64
	a.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Meta"))
		index, _ = b.NextSequence()
		return nil
	})

	return (hex.EncodeToString(byt) + strconv.FormatUint(index, 10))
}

// respond writes out a JSON serializer of whatever is passed in to the HTTP
// response body.
func (a *API) respond(data interface{}, w http.ResponseWriter, r *http.Request) {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		InternalError.ServeHTTP(w, r)
		log.Println("Failed to serialize JSON")
		return
	}

	w.Write(json)
}
