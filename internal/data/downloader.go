package data

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DataLoader struct {
	url               string
	tempDir           string
	tempFilePath      string
	extractedFileName string
}

func NewDataLoader(url string) DataLoader {
	return DataLoader{
		url:          url,
		tempFilePath: "db.gz",
	}
}

func (d *DataLoader) downloadData(url string) error {
	err := d.createTempDir()
	if err != nil {
		return err
	}

	log.Println("downloading")

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", response.StatusCode)
	}

	log.Println("writing")

	d.tempFilePath = filepath.Join(d.tempDir, d.tempFilePath)
	out, err := os.Create(d.tempFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	log.Println("success")

	return nil
}

func (d *DataLoader) extractGzFileXml() error {
	log.Println("extracting file")
	gzipFile, err := os.Open(d.tempFilePath)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	reader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	tempFileSuffix := filepath.Ext(d.tempFilePath)
	d.extractedFileName = fmt.Sprintf(
		"%s.%s",
		strings.TrimSuffix(d.tempFilePath, tempFileSuffix),
		"xml",
	)
	outputFile, err := os.Create(d.extractedFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		return err
	}

	log.Println("success")

	return nil
}

func (d *DataLoader) Get() (string, error) {
	var err error

	err = d.downloadData(d.url)
	if err != nil {
		return "", err
	}

	err = d.extractGzFileXml()
	if err != nil {
		return "", err
	}

	return d.extractedFileName, nil
}

func (d *DataLoader) Close() {
	d.cleanTempDir()
}
