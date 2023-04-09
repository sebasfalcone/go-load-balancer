package health

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type Tests struct {
	name           string
	server         *httptest.Server
	response       *http.Response
	expectedStatus bool
}

func TestNewHealth(t *testing.T) {

	tests := []Tests{
		{
			name: "Server alive",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

			})),
			response: &http.Response{
				Status: "200 OK",
			},
			expectedStatus: true,
		},
		{
			name: "Server dead",
			server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// This handler will never be called because the server is unsarted
			})),
			response:       nil,
			expectedStatus: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			// After creating a new health object, the check function will be called
			serverUrl, _ := url.Parse(test.server.URL)
			health := NewHealth(serverUrl)

			// Test Status function
			if health.Status() != test.expectedStatus {
				t.Errorf("FAILED: Expected status %v, got %v", test.expectedStatus, health.alive)
			}

			// Test SetHealthCheck function
			// Is the same as the original check function but it returns the opposite values
			var checkFunction = func(addr *url.URL) bool {
				conn, err := net.DialTimeout("tcp", addr.Host, defaultCheckTimeout)
				if err != nil {
					return true
				}
				_ = conn.Close()
				return false
			}

			health.SetHealthCheck(checkFunction, time.Second*1)
			expectedStatus := !test.expectedStatus
			if health.Status() != expectedStatus {
				t.Errorf("FAILED: Expected status %v, got %v", test.expectedStatus, health.alive)
			}
		})
	}
}
