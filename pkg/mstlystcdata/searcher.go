package mstlystcdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
)

type Searcher struct {
	c *Cache

	jobs bleve.Index
}

func NewSearcher(c *Cache) (*Searcher, error) {
	s := &Searcher{
		c: c,
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

func (s *Searcher) addDataToIndex() {
	// Fill jobs search from cache
	for _, k := range s.c.jobs.Keys() {
		job, ok := s.c.jobs.Get(k)
		if !ok {
			continue
		}

		s.jobs.Index(k, job)
	}
}

func (s *Searcher) SearchJobs(ctx context.Context, search string, exactMatch bool) ([]*jobs.Job, error) {
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

	jobs := make([]*jobs.Job, len(searchResult.Hits))
	for i, result := range searchResult.Hits {
		job, ok := s.c.jobs.Get(result.ID)
		if !ok {
			return nil, fmt.Errorf("no job found for search result id %s", result.ID)
		}

		jobs[i] = job
	}

	return jobs, nil
}
