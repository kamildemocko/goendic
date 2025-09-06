package main

import (
	"log"
	"test/internal/database"
)

const downloadUrl = `https://en-word.net/static/english-wordnet-2024.xml.gz`
const downloadFilename = "db.gz"

func main() {
	db := database.NewDatabase(downloadUrl, downloadFilename)
	file, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println(file)
}
