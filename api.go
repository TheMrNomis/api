package main

import (
    "log"
    "time"
    "io"
    "bytes"

    "net/http"
    "github.com/gorilla/mux"
)

type MockDB struct {
    UserID      string
    DataName    string
    Data        string
}

var db []MockDB

func GetDataEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    for _, item := range db {
        if item.UserID == params["UserID"] && item.DataName == params["DataName"] {
            io.WriteString(w,item.Data)
            return
        }
    }

    w.WriteHeader(404)
    io.WriteString(w, "404 page not found")
}

func SetDataEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    data := buf.String()

    for _, item := range db {
        if item.UserID == params["UserID"] && item.DataName == params["DataName"] {
            item.Data = data
            return
        }
    }
    db = append(db, MockDB{UserID: params["UserID"], DataName: params["DataName"], Data: data})
}

func main() {
    router := mux.NewRouter()

    db = append(db, MockDB{UserID: "n0m1s", DataName: "data1", Data: `{foo:"data1"}`})
    db = append(db, MockDB{UserID: "n0m1s", DataName: "data2", Data: `{bar:"data2"}`})

    scriptData := router.PathPrefix("/script-data/{UserID}/{DataName}").Subrouter()
    scriptData.HandleFunc("/", GetDataEndpoint).Methods("GET")
    scriptData.HandleFunc("/", SetDataEndpoint).Methods("POST")

    srv := &http.Server{
        Handler:        router,
        Addr:           ":12345",
        WriteTimeout:   15 * time.Second,
        ReadTimeout:    15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}
