package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func authMiddleware(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if strings.ToLower(token) == "token example-token" {
		return nil
	}
	w.WriteHeader(http.StatusUnauthorized)
	return errors.New("401 Unauthorized")
}

func pick(l []string) string {
	index := rand.Int() % len(l)
	return l[index]
}

func generateList(min, max, length int) []int {
	var res []int
	division := max - min + 1
	for i := 0; i < length; i++ {
		item := rand.Int()%division + min
		res = append(res, item)
	}
	return res
}

func fruitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	if err := authMiddleware(w, r); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	data := map[string]string{
		"fruit": pick([]string{"apple", "orange", "banana"}),
	}
	resp, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, string(resp))
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	if err := authMiddleware(w, r); err != nil {
		return
	}
	data := generateList(0, 10, 10)
	resp, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, string(resp))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	if err := authMiddleware(w, r); err != nil {
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req interface{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Submit with body %v\n", req)
	data := map[string]string{
		"message": "ok",
	}
	resp, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, string(resp))
}

func main() {
	http.HandleFunc("/api/examples/fruit", fruitHandler)
	http.HandleFunc("/api/examples/weight", weightHandler)
	http.HandleFunc("/api/examples/submit", submitHandler)
	fmt.Println("Start server http://127.0.0.1:4983")
	log.Fatal(http.ListenAndServe(":4983", nil))
}
