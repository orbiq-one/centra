package content

import (
	"encoding/json"
	"fmt"

	"github.com/cheetahbyte/centra/internal/cache"
	"github.com/cheetahbyte/centra/internal/config"
)

func GetEntry(collection, slug string) (config.ContentRecord, error) {
	lookupPath := collection + "/" + slug
	fmt.Printf("looking for %s in cache", lookupPath)
	bytes := cache.Get(lookupPath)

	var data map[string]any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return config.ContentRecord{}, err
	}

	return data, nil
}
