import (
    "fmt"
    "net/http"
    "math/rand"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func RandStringRunes(n int) string {
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }

    return string(b)
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


func Handler(w http.ResponseWriter, r *http.Request) {
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