package filetype

import (
	"errors"
	"os"
	"unicode/utf8"

	"github.com/h2non/filetype"
	"github.com/liamg/magic"
)

const maxBytes = 1024

func Type(filename string) (ft string, err error) {
	ft = ""
	f, err := os.Open(filename)
	if err != nil {
		return ft, err
	}
	defer f.Close()

	buf := make([]byte, maxBytes)
	n, err := f.Read(buf)
	if err != nil {
		return ft, err
	}

	match, err := filetype.Match(buf)
	if err != nil {
		return ft, err
	}
	if match.MIME.Value != "" {
		return match.MIME.Value, nil
	}

	ma, err := magic.Lookup(buf)
	if err != nil {
		if !errors.Is(err, magic.ErrUnknown) {
			return ft, err
		}
	} else {
		if ma.Description != "" {
			return ma.Description, nil
		}
	}

	if utf8.ValidString(string(buf)) {
		return "text/UTF-8", nil
	}

	for _, b := range buf[:n] {
		if b < 32 || b >= 127 {
			return "unknown", nil
		}
	}

	return "text/ASCII", nil
}
