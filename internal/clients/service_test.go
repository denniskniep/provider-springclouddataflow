package clients

import "testing"

func skipIfIsShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func getJsonConfig() string {
	return `{
		"url": "http://localhost:9393"
	}`
}
