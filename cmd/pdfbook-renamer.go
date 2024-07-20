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
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/webassembly"
)

const pagesToScan = 15

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal error" + err.Error())
	}
}

type pdfengine struct {
	instance pdfium.Pdfium
	pool     pdfium.Pool
}

func NewPdfEngine() (*pdfengine, error) {
	// initialize pdfium
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		return nil, err
	}
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return nil, err
	}
	return &pdfengine{instance, pool}, nil
}

func (p *pdfengine) Destroy() {
	p.instance.Close()
	p.pool.Close()
}

func run() error {
	// initialize pdfium (with wasm runtime)
	pdfengine, err := NewPdfEngine()
	if err != nil {
		return err
	}
	defer pdfengine.Destroy()

	// get pdf file list
	var pdflist []string
	err = filepath.Walk("/data/",
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
		isbn, err := isbn.FindPdfISBN(path, pagesToScan, pdfengine.instance)
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
