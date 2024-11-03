package mstlystcdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"go.uber.org/multierr"
)

type Searcher struct {
	cache *Cache

	jobs bleve.Index
}

func NewSearcher(cache *Cache) (*Searcher, error) {
	s := &Searcher{
		cache: cache,
	}

	jobs, err := s.newJobsIndex()
	if err != nil {
		return nil, err
	}
	s.jobs = jobs

	return s, nil
}

func (s *Searcher) newJobsIndex() (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()

	jobMapping := bleve.NewDocumentMapping()
	gradesMapping := bleve.NewDocumentMapping()
	jobMapping.AddSubDocumentMapping("grades", gradesMapping)
	indexMapping.AddDocumentMapping("job", jobMapping)

	return bleve.NewMemOnly(indexMapping)
}

func (s *Searcher) addDataToIndex(ctx context.Context) error {
	errs := multierr.Combine()

	batch := s.jobs.NewBatch()
	// Fill jobs search from cache
	s.cache.jobs.Range(ctx, func(key string, value *users.Job) bool {
		batch.Delete(key)

		if err := batch.Index(key, value); err != nil {
			errs = multierr.Append(errs, err)
		}

		return true
	})

	if err := s.jobs.Batch(batch); err != nil {
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (s *Searcher) SearchJobs(ctx context.Context, search string, exactMatch bool) ([]*users.Job, error) {
	var searchQuery query.Query
	if search == "" {
		searchQuery = bleve.NewMatchAllQuery()
	} else {
		if exactMatch {
			searchQuery = bleve.NewMatchQuery(strings.ToLower(search))
		} else {
			searchQuery = bleve.NewQueryStringQuery(strings.ToLower(search) + "*")
		}
	}

	searchRequest := bleve.NewSearchRequest(searchQuery)
	if exactMatch {
		searchRequest.Size = 1
	} else {
		searchRequest.Size = 32
	}
	searchRequest.Fields = []string{"label", "name"}
	searchRequest.SortBy([]string{"label", "_id"})

	searchResult, err := s.jobs.SearchInContext(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	jobs := make([]*users.Job, len(searchResult.Hits))
	for i, result := range searchResult.Hits {
		job, ok := s.cache.jobs.Get(result.ID)
		if !ok {
			return nil, fmt.Errorf("no job found for search result id %s", result.ID)
		}

		jobs[i] = job
	}

	return jobs, nil
}
