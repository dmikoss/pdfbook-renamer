package isbn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Provider interface {
	Fetch(ctx context.Context, isbn string) ProviderInfo
}

type ProviderInfo struct {
	Title         string
	YearOfPublish int
	Authors       []string
}

type ProviderGoogleBooks struct {
	client  *http.Client
	baseURL *url.URL
}

type response struct {
	TotalItems int            `json:"totalItems,omitempty"`
	Items      []responseItem `json:"items,omitempty"`
}
type responseItem struct {
	VolumeInfo responseVolumeInfo `json:"volumeInfo,omitempty"`
}

type responseVolumeInfo struct {
	Title         string   `json:"title,omitempty"`
	PublishedDate string   `json:"publishedDate,omitempty"`
	Authors       []string `json:"authors,omitempty"`
}

const urlGoogleBooksApi = "https://www.googleapis.com/books/v1/volumes"

func NewProviderGoogleBooks(httpclient *http.Client) *ProviderGoogleBooks {
	serverUrl, _ := url.Parse(urlGoogleBooksApi)
	return &ProviderGoogleBooks{
		client:  httpclient,
		baseURL: serverUrl,
	}
}

func (p *ProviderGoogleBooks) Fetch(ctx context.Context, isbn string) (ProviderInfo, error) {
	querypath := *p.baseURL
	q := querypath.Query()
	q.Set("q", fmt.Sprintf("isbn:%v", isbn))
	querypath.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, querypath.String(), nil)
	if err != nil {
		return ProviderInfo{}, fmt.Errorf("http.NewRequestWithContext error: %w", err)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return ProviderInfo{}, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProviderInfo{}, fmt.Errorf("client.Do: %w", err)
	}

	var value response
	err = json.Unmarshal(body, &value)
	if err != nil {
		return ProviderInfo{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if len(value.Items) > 0 {
		item := value.Items[0]
		return ProviderInfo{
			item.VolumeInfo.Title,
			getYear(item.VolumeInfo.PublishedDate),
			item.VolumeInfo.Authors}, nil
	}

	return ProviderInfo{}, nil
}

func getYear(date string) int {
	year := 0
	date1, err := time.Parse("2006-01-02", date)
	if err == nil {
		year = date1.Year()
	} else {
		date2, err := time.Parse("2006", date)
		if err == nil {
			year = date2.Year()
		}
	}
	return year
}
