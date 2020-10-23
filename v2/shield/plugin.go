package shield

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Plugin struct {
	ID   string
	Name string
}

func ParseAgentMetadata(in map[string]interface{}) ([]*Plugin, error) {
	var metadata struct {
		Plugins map[string]struct {
			Name string `json:"name"`
		} `json:"plugins"`
	}

	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &metadata)
	if err != nil {
		return nil, err
	}

	if metadata.Plugins != nil {
		var keys []string
		for k := range metadata.Plugins {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var l []*Plugin
		for _, id := range keys {
			p := metadata.Plugins[id]
			l = append(l, &Plugin{
				ID:   id,
				Name: p.Name,
			})
		}
		return l, nil
	}

	return nil, fmt.Errorf("unable to detect plugins in agent metadata")
}
