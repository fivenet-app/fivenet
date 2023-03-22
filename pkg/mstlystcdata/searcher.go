package mstlystcdata

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/proto/resources/jobs"
)

type Searcher struct {
	c *Cache

	docCategories bleve.Index
	jobs          bleve.Index
}

func NewSearcher(c *Cache) (*Searcher, error) {
	s := &Searcher{
		c: c,
	}

	docCategories, err := s.newDocsCategoriesIndex()
	if err != nil {
		return nil, err
	}
	s.docCategories = docCategories

	jobs, err := s.newJobsIndex()
	if err != nil {
		return nil, err
	}
	s.jobs = jobs

	return s, nil
}

func (s *Searcher) newDocsCategoriesIndex() (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultField = "name"

	jobMapping := bleve.NewDocumentDisabledMapping()
	categoryMapping := bleve.NewDocumentMapping()
	categoryMapping.AddSubDocumentMapping("Job", jobMapping)

	indexMapping.AddDocumentMapping("category", categoryMapping)

	return bleve.NewMemOnly(indexMapping)
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
	// Fill document categories search from cache
	for _, k := range s.c.docCategories.Keys() {
		cat, ok := s.c.docCategories.Get(k)
		if !ok {
			continue
		}

		id := strconv.Itoa(int(cat.Id))
		s.docCategories.Index(id, cat)
	}

	// Fill jobs search from cache
	for _, k := range s.c.jobs.Keys() {
		job, ok := s.c.jobs.Get(k)
		if !ok {
			continue
		}

		s.jobs.Index(k, job)
	}
}

func (s *Searcher) SearchDocumentCategories(ctx context.Context, search string, searchJobs []string) ([]*documents.DocumentCategory, error) {
	if len(searchJobs) == 0 {
		return []*documents.DocumentCategory{}, nil
	}

	var userSearch query.Query
	if search == "" {
		userSearch = bleve.NewMatchAllQuery()
	} else {
		userSearch = bleve.NewMatchPhraseQuery(search)
	}

	queries := make([]query.Query, len(searchJobs))
	for i := 0; i < len(searchJobs); i++ {
		jobsQuery := bleve.NewTermQuery(searchJobs[i])
		jobsQuery.SetField("job")

		queries[i] = jobsQuery
	}

	searchQuery := bleve.NewBooleanQuery()
	searchQuery.Must = bleve.NewDisjunctionQuery(queries...)
	searchQuery.Should = userSearch

	searchRequest := bleve.NewSearchRequest(searchQuery)
	searchRequest.Size = 10
	searchRequest.Fields = []string{"id", "name", "description", "job"}
	searchRequest.SortBy([]string{"job", "name", "_id"})

	searchResult, err := s.docCategories.SearchInContext(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	categories := make([]*documents.DocumentCategory, len(searchResult.Hits))
	for i, result := range searchResult.Hits {
		id, err := strconv.Atoi(result.ID)
		if err != nil {
			return nil, err
		}

		category, ok := s.c.docCategories.Get(uint64(id))
		if !ok {
			return nil, fmt.Errorf("no document category found for search result id %s (%d)", result.ID, id)
		}

		categories[i] = category
	}

	return categories, nil
}

func (s *Searcher) SearchJobs(ctx context.Context, search string) ([]*jobs.Job, error) {
	var searchQuery query.Query
	if search == "" {
		searchQuery = bleve.NewMatchAllQuery()
	} else {
		searchQuery = bleve.NewQueryStringQuery(strings.ToLower(search) + "*")
	}

	searchRequest := bleve.NewSearchRequest(searchQuery)
	searchRequest.Size = 15
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

func CleanupSearchInput(search string) string {
	return strings.ToLower(strings.Trim(search, "*")) + "*"
}
