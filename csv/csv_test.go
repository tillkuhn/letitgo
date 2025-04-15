package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/stretchr/testify/assert"
)

// thanks https://stackoverflow.com/a/58841827/4292075

func TestName(t *testing.T) {
	records, err := readCsvFile("travel.csv")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(records))
	assert.Equal(t, "Bangkok", records[1][0])
	assert.Equal(t, "Cuba", records[3][1])
	t.Logf("Hacker news of the day: %s", gofakeit.HackerPhrase())
}

// readCsvFile format: Bangkok,Thailand
func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file %s: %w", filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse csv for %s: %w", filePath, err)
	}

	return records, nil
}
