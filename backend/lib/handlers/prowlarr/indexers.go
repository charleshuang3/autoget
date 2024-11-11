package prowlarr

import (
	"encoding/json"
	"net/http"
)

type Indexer struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Indexers []Indexer `json:"indexers"`
}

func (p *Prowlarr) IndexersHandler(w http.ResponseWriter, r *http.Request) {
	indexers, err := p.client.GetIndexers()
	if err != nil {
		http.Error(w, "Failed to get indexers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Indexers: []Indexer{},
	}

	for _, indexer := range indexers {
		response.Indexers = append(response.Indexers, Indexer{
			ID:   indexer.ID,
			Name: indexer.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
