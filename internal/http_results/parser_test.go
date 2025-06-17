package http_results

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var post = []byte(
	`# POST /api/v1/users/{uuid}/settings 200 application/json
{
	"id": "{1}",
	"name": "test name",
	"email": "test@example.com",
}
`)

var postResult = []byte(`{
	"id": "{1}",
	"name": "test name",
	"email": "test@example.com",
}`)

func TestParser(t *testing.T) {
	result, err := Parse(post)
	require.Nil(t, err, "error parsing")
	require.NotNil(t, result, "result is nil")

	assert.Equal(t, 200, result.Code, "code is not 200")
	assert.Equal(t, "POST", result.Method, "method is not POST")
	assert.Equal(t, "/api/v1/users/{uuid}/settings", result.Path, "path is not correct")
	assert.Equal(t, "application/json", result.ContentType, "content type is not correct")
	assert.Equal(t, postResult, result.Data, "data is not correct")
}
