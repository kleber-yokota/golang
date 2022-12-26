package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	price, err := GetPrice()
	if err != nil {
		panic(err)
	}
	err = SavePrice(price)
	if err != nil {
		panic(err)
	}

}

func SavePrice(price string) error {
	value := FormatString(price)
	f, err := CreateFile("cotacao.txt")
	if err != nil {
		return err
	}
	err = Save(f, value)
	if err != nil {
		return err
	}
	return nil
}

func FormatString(price string) string {
	return "Dolar:" + price
}

func CreateFile(name string) (*os.File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Save(f *os.File, value string) error {
	_, err := f.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

func GetPrice() (string, error) {
	url := "http://localhost:8080/cotacao"

	client := http.Client{
		Timeout: 300 * time.Millisecond,
	}

	req, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil

}
