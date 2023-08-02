package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmikoss/pdfbook-renamer/isbn"
)

const pagesToScan = 15

type rawBook struct {
	path string
	isbn string
	info isbn.ProviderInfo
}

type Book struct {
	rawBook
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal error" + err.Error())
	}
}

func run() error {

	// get pdf file list
	var pdflist []string
	err := filepath.Walk("/data/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.Contains(info.Name(), ".pdf") || info.IsDir() {
				return nil
			}
			pdflist = append(pdflist, path)
			return nil
		})

	if err != nil {
		return err
	}

	// pdf processing
	fetcher := isbn.NewProviderGoogleBooks(http.DefaultClient)
	for _, path := range pdflist {
		isbn, err := isbn.FindPdfISBN(path, pagesToScan)
		if err != nil || isbn == "" {
			continue
		} else {
			info, err := fetcher.Fetch(context.Background(), isbn)
			if err != nil {
				continue
			}

			if info.Title != "" && len(info.Authors) > 0 {
				dir := filepath.Dir(path)
				renamed := fmt.Sprintf("%v - %v - %d.pdf", info.Title, info.Authors[0], info.YearOfPublish)
				dstpath := filepath.Join(dir, renamed)

				if err := os.Rename(path, dstpath); err != nil {
					continue
				}
			}
			fmt.Println(info)
		}
	}

	return err
}
