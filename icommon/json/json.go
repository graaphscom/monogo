package json

import (
	"encoding/json"
	"os"
)

type IcoManifest struct {
	BasePath     string
	VendorsPaths map[string]struct {
		Icons    string
		Metadata string
	}
}

func ReadJson[T any](path string) (T, error) {
	contents, err := os.ReadFile(path)

	var result T

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(contents, &result)

	return result, err
}
