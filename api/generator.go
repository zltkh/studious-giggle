package handler

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net/http"
	"strconv"
)

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func writeToDB(token, prompt, num_beams, temperature, top_p, num_beam_groups string) error {
	db, err := sql.Open("mysql", "l905412p_project:4V6NR5Cg@tcp(l905412p.beget.tech)/l905412p_project")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO requests(token, prompt, num_beams, temperature, top_p, num_beam_groups) VALUES(?, ?, ?, ?, ?, ?)", token, prompt, num_beams, temperature, top_p, num_beam_groups)
	return err
}

func isPositiveInteger(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return false
	}
	return true
}

func isPositiveNumber(s string) bool {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil || num <= 0 {
		return false
	}
	return true
}

func isPositiveNumber2(s string) bool {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil || num <= 0 || num > 1 {
		return false
	}
	return true
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                // Разрешение запросов со всех источников
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Разрешение методов запроса
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Разрешение заголовков запроса

	if r.Method != "POST" {
		return
	}
	prompt := r.FormValue("prompt")
	num_beams := r.FormValue("num_beams")
	temperature := r.FormValue("temperature")
	top_p := r.FormValue("top_p")
	num_beam_groups := r.FormValue("num_beam_groups")

	token := RandStringRunes(32)

	if len(prompt) > 400 {
		http.Error(w, "Prompt is too large. Max length is 400 symbols.", http.StatusBadRequest)
		return
	}
	if len(num_beams) == 0 {
		num_beams = "4"
	}

	if len(temperature) == 0 {
		temperature = "1"
	}

	if len(top_p) == 0 {
		top_p = "0.8"
	}

	if len(num_beam_groups) == 0 {
		num_beam_groups = "4"
	}

	if !isPositiveInteger(num_beams) {
		http.Error(w, "num_beams should be positive integer", http.StatusBadRequest)
		return
	}

	if !isPositiveInteger(num_beam_groups) {
		http.Error(w, "num_beam_group should be positive integer", http.StatusBadRequest)
		return
	}

	if !isPositiveNumber(temperature) {
		http.Error(w, "temperature should be positive number", http.StatusBadRequest)
		return
	}

	if !isPositiveNumber2(top_p) {
		http.Error(w, "top_p should be positive number greater than 0 and less then 1", http.StatusBadRequest)
		return
	}

	if len(prompt) == 0 {
		http.Error(w, "Promp must be specified", http.StatusBadRequest)
		return
	}

	if err := writeToDB(token, prompt, num_beams, temperature, top_p, num_beam_groups); err != nil {
		http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, token)
}
