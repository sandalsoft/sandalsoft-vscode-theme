package search

import (
	"sync"

	"github.com/sandalsoft/rubrik-search/models"
)

type SiteSearcher interface {
	SearchSite(query string, wg *sync.WaitGroup, outputChan chan SearchChanData)
}

type SearchRepository struct {
	APIEndpoint       string
	RepositoryBaseUrl string
}
type SearchChanData struct {
	Err           error
	SearchResults []models.SearchResult
	HasMoreData   bool
}

type SearchService struct {
	query         string
	searchSources []SiteSearcher
	OutputChan    chan SearchChanData
	SearchWG      *sync.WaitGroup
}

func NewSearchService(sources ...SiteSearcher) SearchService {
	var wg sync.WaitGroup
	outputChan := make(chan SearchChanData)
	return SearchService{
		searchSources: sources,
		OutputChan:    outputChan,
		SearchWG:      &wg,
	}
}

func (svc *SearchService) Search(query string) ([]models.SearchResult, error) {
	for _, source := range svc.searchSources {
		svc.SearchWG.Add(1)
		go source.SearchSite(query, svc.SearchWG, svc.OutputChan)
	}

	go func() {
		svc.SearchWG.Wait()
		close(svc.OutputChan)
	}()

	var allResults []models.SearchResult

	for {
		chanData, hasMoreData := <-svc.OutputChan
		if hasMoreData == false {
			break
		}
		allResults = append(allResults, chanData.SearchResults...)
	}
	return allResults, nil
}
