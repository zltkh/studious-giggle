package main

import (
    "fmt"
    "log"
    "net/http"
    "math/rand"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// PASSWORD 
// DBNAME 

func RandStringRunes(n int) string {
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }

    return string(b)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    http.HandleFunc("/generate", generate)
    http.HandleFunc("/get_poem", getHandler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeToDB(token, prompt string) error {
    db, err := sql.Open("mysql", "l905412p_project:4V6NR5Cg@tcp(l905412p.beget.tech)/l905412p_project")
    if err != nil {
        return err
    }
    defer db.Close()
    _, err = db.Exec("INSERT INTO requests(token, prompt) VALUES(?, ?)", token, prompt)
    return err 
}

func getFromDB(token string) (string, error) {
    db, err := sql.Open("mysql", "l905412p_project:4V6NR5Cg@tcp(l905412p.beget.tech)/l905412p_project")
    if err != nil {
        return "", err
    }
    defer db.Close()
    row := db.QueryRow("SELECT * FROM requests WHERE token=?", token)
    var id int
    var tokenn string
    var result string
    var pprompt string
    err = row.Scan(&id, &tokenn, &result, &pprompt)
    if err != nil {
        return "", err
    }
    return result, nil
}

func generate(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        return
    }
    prompt := r.FormValue("prompt")
    token := RandStringRunes(32)
    if len(prompt) > 400 {
        http.Error(w, "Prompt is too large. Max length is 400 symbols.", http.StatusBadRequest)
        return
    }

    if len(prompt) == 0 {
        http.Error(w, "Promp must be specified", http.StatusBadRequest)
        return
    }

    if err := writeToDB(token, prompt); err != nil {
        http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, token)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    fmt.Println(token)
    if len(token) != 32 {
        http.Error(w, "Token is wrong", http.StatusBadRequest)
        return
    }
    res, err := getFromDB(token)
    if err != nil {
        http.Error(w, "Internal error :(", http.StatusInternalServerError)
        return
    }
    if len(res) == 0 {
        http.Error(w, "Result not ready :(", http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, res)
}
