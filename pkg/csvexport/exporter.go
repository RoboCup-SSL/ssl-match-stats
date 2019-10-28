package csvexport

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

func writeCsv(header []string, data [][]string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Could not create CSV output file")
	}

	if _, err := f.WriteString("#" + strings.Join(header, ",") + "\n"); err != nil {
		return err
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(data); err != nil {
		return err
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return f.Close()
}

func uintToStr(n uint32) string {
	return strconv.FormatUint(uint64(n), 10)
}

func intToStr(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func floatToStr(n float32) string {
	return strconv.FormatFloat(float64(n), 'f', 2, 32)
}
