package main

import (
	"goendic/internal/database"
	"log"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`

func main() {
	db := database.NewDatabase(downloadUrl)
	file, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println(file)

	database.ParseXML(file)
}
