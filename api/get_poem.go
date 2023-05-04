package handler

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

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
    err = row.Scan(&id, &tokenn, &pprompt, &result)
    if err != nil {
        return "", err
    }
    return result, nil
}
 
func Handler(w http.ResponseWriter, r *http.Request) {
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