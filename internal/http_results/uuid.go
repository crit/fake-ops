package http_results

import (
	"bytes"
	"strconv"

	"github.com/google/uuid"
)

// FillUUID attempts to replace any instance of $1, $2, etc in some content
// with new UUIDs.
func FillUUID(content []byte, count int) []byte {
	for i := 0; i < count; i++ {
		// id: "$1", => id: "8BC48765-6456-4F70-9B73-E03CE3760F44",
		token := "$" + strconv.Itoa(i+1)
		content = bytes.ReplaceAll(content, []byte(token), []byte(uuid.NewString()))
	}

	return content
}
