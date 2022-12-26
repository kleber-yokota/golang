package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type USDBRL struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	http.HandleFunc("/cotacao", Cotacao)
	http.ListenAndServe(":8080", nil)

}

func Cotacao(w http.ResponseWriter, r *http.Request) {
	var result USDBRL

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/dolar")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = GetApi(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = InsertTable(db, &result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.Usdbrl.Bid)

}

func InsertTable(db *sql.DB, value *USDBRL) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, "insert into exchange(Code, Codein, Name, High, Low, VarBid, PctChange, Bid, Ask, Timestamp, CreateDate) values(?,?,?,?,?,?,?,?,?,?,?) ")

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		value.Usdbrl.Code,
		value.Usdbrl.Codein,
		value.Usdbrl.Name,
		value.Usdbrl.High,
		value.Usdbrl.Low,
		value.Usdbrl.VarBid,
		value.Usdbrl.PctChange,
		value.Usdbrl.Bid,
		value.Usdbrl.Ask,
		value.Usdbrl.Timestamp,
		value.Usdbrl.CreateDate,
	)
	if err != nil {
		return err
	}
	defer cancel()
	return nil

}

func GetApi(res *USDBRL) error {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	req, err := client.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	err = GetJson(req.Body, res)
	if err != nil {
		return err
	}
	return nil

}

func GetJson(req io.Reader, result *USDBRL) error {

	res, err := io.ReadAll(req)
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return err
	}

	return nil
}
