package proxy

import (
	"net/url"
	"testing"
	"time"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/proxy/mocks"
)

type Tests struct {
	name           string
	url            *url.URL
	expectedStatus bool
	expectedLoad   int32
}

func TestNewProxy(t *testing.T) {

	validUrl, _ := url.Parse("http://valid")
	invalidUrl, _ := url.Parse("http://invalid")

	tests := []Tests{
		{
			name:           "valid url",
			url:            validUrl,
			expectedStatus: true,
			expectedLoad:   0,
		},
		{
			name:           "invalid url",
			url:            invalidUrl,
			expectedStatus: false,
			expectedLoad:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			h := mocks.NewHealth(test.url)
			p := NewProxy(test.url, test.name)
			p.SetHealthCheck(h.CheckFunction, time.Second*1)

			// Check name
			if p.GetName() != test.name {
				t.Errorf("FAILED: Expected name %s, got %s", test.name, p.GetName())
			}

			// Check health
			if p.Status() != test.expectedStatus {
				t.Errorf("FAILED: Expected status %v, got %v", test.expectedStatus, p.Status())
			}

			// Check load
			if p.GetLoad() != test.expectedLoad {
				t.Errorf("FAILED: Expected load %d, got %d", test.expectedLoad, p.GetLoad())
			}

		})
	}
}
