package isbn

import (
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

func FindPdfISBN(path string, numpages int) ( /*isbn*/ string /*err*/, error) {
	patterns := [...]string{
		`(?:ISBN:?\s?)?(?:978-|979-)?(?:[0-9]{1,5}-)(?:[0-9]{1,7}-)(?:[0-9]{1,6}-)(?:[0-9X]{1})`,
		`(?:97[89]([- ])?)?(?:\d{10}|(?=(?:\d\1?){9}\1?[xX\d]$)\d+(?:\1?\d+){2}\1?[\dxX])`,
		`(?:ISBN(?:-13)?:?\ )?(?=[0-9]{13}$|(?=(?:[0-9]+[-\ ]){4})[-\ 0-9]{17}$)97[89][-\ ]?[0-9]{1,5}[-\ ]?[0-9]+[-\ ]?[0-9]+[-\ ]?[0-9]`,
	}

	var regexps []*regexp2.Regexp
	for _, rstr := range patterns {
		regexps = append(regexps, regexp2.MustCompile(rstr, 0))
	}

	cmd := exec.Command("python", "/app/pdf-to-text.py", path, strconv.Itoa(numpages))
	pipe, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		if err != nil {
			return "", err
		}
	}

	bytes, err := io.ReadAll(pipe)
	if err != nil {
		return "", err
	}
	str := string(bytes[:])
	for _, rg := range regexps {
		if match, _ := rg.FindStringMatch(str); match != nil {
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
