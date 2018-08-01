package main

import (
	"context"
	"reflect"
	"time"

	"github.com/monochromegane/synapse"
	"github.com/olivere/elastic"
)

type searcher struct {
	client *elastic.Client
}

func (s *searcher) Search(ctx synapse.Context, profile synapse.Profile, association synapse.Association) ([]string, []float64, error) {
	ids := []string{}
	scores := []float64{}

	termQuery := elastic.NewTermQuery("category_id", association["category_id"])
	hits, err := s.client.Search().
		Index("synapse").
		Query(termQuery).
		Sort("created_at", false).
		From(0).Size(3).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return ids, scores, err
	}

	var product Product
	for i, item := range hits.Each(reflect.TypeOf(product)) {
		if p, ok := item.(Product); ok {
			ids = append(ids, p.Name)
			scores = append(scores, float64(i))
		}
	}

	return ids, scores, nil
}

func (s *searcher) Initialize() error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

func (s *searcher) Finalize() error {
	return nil
}

func NewSearcher() synapse.Searcher {
	return &searcher{}
}

type Product struct {
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}
