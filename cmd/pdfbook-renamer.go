package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmikoss/pdfbook-renamer/isbn"
	"github.com/klippa-app/go-pdfium/webassembly"
)

const pagesToScan = 15

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

	// initialize pdfium
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		return err
	}
	defer pool.Close()

	pdfiumInst, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return err
	}
	defer pdfiumInst.Close()

	// pdf processing
	fetcher := isbn.NewProviderGoogleBooks(http.DefaultClient)
	for _, path := range pdflist {
		isbn, err := isbn.FindPdfISBN(path, pagesToScan, pdfiumInst)
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
