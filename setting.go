package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

type Setting struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPass     string
	PgName     string
	Data       string
	Assets     string
	Html       string
	Email      string
}

var cfg Setting

func init() {
	file, err := os.Open("setting.cfg")
	if err != nil {
		fmt.Println("Error - open file setting.cfg", err)
	}
	defer file.Close()

	var stat fs.FileInfo
	stat, err = file.Stat()
	if err != nil {
		fmt.Println("Error - stat file setting.cfg", err)
	}
	byteFile := make([]byte, stat.Size())

	_, err = file.Read(byteFile)
	if err != nil {
		fmt.Println("Error- read file setting.cfg", err)
	}
	err = json.Unmarshal(byteFile, &cfg)
	if err != nil {
		fmt.Println("Error - unmarshal file setting.cfg", err)
	}

}
