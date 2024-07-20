package isbn

import (
	"os"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
)

func FindPdfISBN(path string, numpages int, pdfiumInst pdfium.Pdfium) ( /*isbn*/ string /*err*/, error) {
	patterns := [...]string{
		`(?:ISBN:?\s?)?(?:978-|979-)?(?:[0-9]{1,5}-)(?:[0-9]{1,7}-)(?:[0-9]{1,6}-)(?:[0-9X]{1})`,
		`(?:97[89]([- ])?)?(?:\d{10}|(?=(?:\d\1?){9}\1?[xX\d]$)\d+(?:\1?\d+){2}\1?[\dxX])`,
		`(?:ISBN(?:-13)?:?\ )?(?=[0-9]{13}$|(?=(?:[0-9]+[-\ ]){4})[-\ 0-9]{17}$)97[89][-\ ]?[0-9]{1,5}[-\ ]?[0-9]+[-\ ]?[0-9]+[-\ ]?[0-9]`,
	}

	var regexps []*regexp2.Regexp
	for _, rstr := range patterns {
		regexps = append(regexps, regexp2.MustCompile(rstr, 0))
	}

	pdfBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	// Open the PDF using PDFium (and claim a worker)
	doc, err := pdfiumInst.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return "", err
	}

	// Get the number of pages in the PDF
	pageCount, err := pdfiumInst.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return "", err
	}

	var pagestext string
	// Get the text from the pages
	for i := 0; i < min(pageCount.PageCount, numpages); i++ {

		textpage, err := pdfiumInst.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    i, // 0-indexed
				},
			},
		})
		if err != nil {
			return "", err
		}

		text, err := pdfiumInst.FPDFText_GetText(&requests.FPDFText_GetText{
			TextPage:   textpage.TextPage,
			StartIndex: 0,
			Count:      10000,
		})
		if err != nil {
			return "", err
		}
		pagestext += text.Text
	}

	// Always close the document, this will release its resources.
	pdfiumInst.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	for _, rg := range regexps {
		if match, _ := rg.FindStringMatch(pagestext); match != nil {
			return cleanISBN(match.String()), nil
		}
	}
	return "", nil
}

func cleanISBN(isbn string) string {
	str := strings.ReplaceAll(isbn, "ISBN", "")
	str = strings.ReplaceAll(str, "-", "")
	str = strings.ReplaceAll(str, ":", "")
	str = strings.TrimSpace(str)
	return str
}
