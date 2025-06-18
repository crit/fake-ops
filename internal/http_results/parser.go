package http_results

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Result is parsed from a http response file.
type Result struct {
	Code        int
	Method      string
	Path        string
	ContentType string
	Data        []byte
}

// Parse takes in the content of a response file and creates a Result.
func Parse(data []byte) (*Result, error) {
	var result Result

	// get first line of data
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	line := scanner.Text()

	// # GET /api/v1/users 200 application/json => ["#", "GET", "/api/v1/users", "200", "application/json"]
	parts := strings.Split(line, " ")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}

	if parts[0] != "#" {
		return nil, fmt.Errorf("invalid line: %s", line)
	}

	result.Method = parts[1]
	result.Path = parts[2]

	var err error
	result.Code, err = strconv.Atoi(parts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid code: %s", parts[3])
	}

	result.ContentType = parts[4]

	// rest of the data is put into result.Data
	result.Data = bytes.TrimSpace(data[len(scanner.Bytes()):])

	return &result, nil
}
