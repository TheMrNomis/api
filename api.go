package main

import (
    "log"
    "time"
    "io"
    "bytes"
    "os"

    "net/http"
    "github.com/gorilla/mux"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "encoding/json"
)

type DatabaseSettings struct {
    Driver  string
    Path    string
}

type WebserverSettings struct {
    Address string
}

type Settings struct {
    Database    *DatabaseSettings
    Webserver   *WebserverSettings
}

var settings Settings
var db *sql.DB

func writeErrorCode(w http.ResponseWriter, code int) {
    w.WriteHeader(code)

    switch code {
    case 404:
        io.WriteString(w, "404 page not found")
    case 500:
        io.WriteString(w, "500 internal server error")
    }
}

func GetDataEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    stmt, err := db.Prepare("SELECT data FROM script_data WHERE userID=? AND dataName=?")
    if err != nil {
        log.Println(err)
        writeErrorCode(w, 500)
        return
    }
    defer stmt.Close()

    var data string
    err = stmt.QueryRow(params["UserID"], params["DataName"]).Scan(&data)
    if err != nil {
        writeErrorCode(w, 404)
        return
    }
    io.WriteString(w, data)
}

func SetDataEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    data := buf.String()

    stmt, err := db.Prepare("INSERT OR REPLACE INTO script_data (userID, dataName, data) VALUES (?,?,?)")
    if err != nil {
        log.Println(err)
        writeErrorCode(w, 500)
    }
    defer stmt.Close()

    _, err = stmt.Exec(params["UserID"], params["DataName"], data)
    if err != nil {
        log.Println(err)
        writeErrorCode(w, 500)
    }
}

func LoadConf(path string) Settings {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    conf := Settings{}
    err = decoder.Decode(&conf)
    if err != nil {
        log.Fatal(err)
    }
    return conf
}

func OpenDatabase(conf *DatabaseSettings) *sql.DB {
    db, err := sql.Open(conf.Driver, conf.Path)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func init() {
    settings = LoadConf("conf.json")
    db = OpenDatabase(settings.Database)
}

func main() {
    defer db.Close()

    router := mux.NewRouter()

    scriptData := router.PathPrefix("/script-data/{UserID}/{DataName}").Subrouter()
    scriptData.HandleFunc("/", GetDataEndpoint).Methods("GET")
    scriptData.HandleFunc("/", SetDataEndpoint).Methods("POST")

    srv := &http.Server{
        Handler:        router,
        Addr:           settings.Webserver.Address,
        WriteTimeout:   15 * time.Second,
        ReadTimeout:    15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}
