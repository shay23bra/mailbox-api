package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseCSV(csvData string) ([]map[string]string, error) {
	reader := csv.NewReader(strings.NewReader(csvData))

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	result := []map[string]string{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		rowMap := make(map[string]string)
		for i, field := range row {
			if i < len(header) {
				rowMap[header[i]] = field
			}
		}

		result = append(result, rowMap)
	}

	return result, nil
}

func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
