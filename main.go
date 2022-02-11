package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

var ctx = context.Background()
var redisConn *redis.Client

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/stat", StatIndexHandler).Methods("GET")
	router.HandleFunc("/stat", StatCreateHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8999", nil)
}

func getRedis() *redis.Client {
	if redisConn != nil {
		return redisConn
	}
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return redisConn
}

func StatIndexHandler(w http.ResponseWriter, r *http.Request) {
	stats := getRedis().HGetAll(ctx, "stat")
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]int)
	for key, value := range stats.Val() {
		response[key], _ = strconv.Atoi(value)
	}

	jsonResponse, _ := json.Marshal(response)
	io.WriteString(w, string(jsonResponse))
}

type statCreateParams struct {
	Country string
}

func StatCreateHandler(w http.ResponseWriter, r *http.Request) {
	var body statCreateParams
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go getRedis().HIncrBy(ctx, "stat", body.Country, 1)
}
