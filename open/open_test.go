package open

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/skratchdot/open-golang/opts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name string

	target      string
	options     []opts.Option
	opener      Opener
	expectError bool
}

func TestOpen(t *testing.T) {
	t.Parallel()
	// Run a local HTTP server to ensure tests are able to executable
	// without a live 3rd party dependency.
	m := &sync.RWMutex{}
	actualVisitsByPath := make(map[string]int)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		defer m.Unlock()
		actualVisitsByPath[r.URL.Path] += 1
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		w.Header().Add(`Content-Type`, `text/html; charset=utf-8`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
		<html>
		<body>server is ok</body>
		<script>window.close()</script>
		</html>
		`))
	}))
	defer server.Close()

	testcases := []TestCase{
		{
			name:        "[Open] URL with wait should not error",
			target:      server.URL,
			options:     []opts.Option{opts.Wait()},
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] xxxxxxxxxxxxxxx with wait should error",
			target:      "xxxxxxxxxxxxxxx",
			options:     []opts.Option{opts.Wait()},
			opener:      OpenerFunc(Open),
			expectError: true,
		},
		{
			name:        "[Open] URL should not error",
			target:      server.URL,
			options:     nil,
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] xxxxxxxxxxxxxxx should not error",
			target:      "xxxxxxxxxxxxxxx",
			options:     nil,
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] URL with appname and wait should not error",
			target:      server.URL,
			options:     []opts.Option{opts.Wait(), opts.WithApp("curl")},
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] xxxxxxxxxxxxxxx with appname and wait should error",
			target:      "xxxxxxxxxxxxxxx",
			options:     []opts.Option{opts.Wait(), opts.WithApp("curl")},
			opener:      OpenerFunc(Open),
			expectError: true,
		},
		{
			name:        "[Open] URL with appname should not error",
			target:      server.URL,
			options:     []opts.Option{opts.WithApp("curl")},
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] xxxxxxxxxxxxxxx with appname should not error",
			target:      "xxxxxxxxxxxxxxx",
			options:     []opts.Option{opts.WithApp("curl")},
			opener:      OpenerFunc(Open),
			expectError: false,
		},
		{
			name:        "[Open] URL with appname that doesn't exist should error",
			target:      server.URL,
			options:     []opts.Option{opts.WithApp("this-does-not-exist")},
			opener:      OpenerFunc(Open),
			expectError: true,
		},
		{
			name:        "[Open] URL with wait and with appname that doesn't exist should error",
			target:      server.URL,
			options:     []opts.Option{opts.Wait(), opts.WithApp("this-does-not-exist")},
			opener:      OpenerFunc(Open),
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
			defer cancel()

			if err := tc.opener.Open(ctx, tc.target, tc.options...); tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	t.Run("path visits - xdg-open curl", func(t *testing.T) {
		expectedVisits := map[string]int{
			"/xdg-open/curl/foo":             1,
			"/xdg-open/curl/bar":             3,
			"/xdg-open/curl/baz":             3,
			"/xdg-open/curl/ctrl/alt/delete": 7,
			"/xdg-open/curl/none":            0,
		}
		for expectedPath, expectedCount := range expectedVisits {
			t.Run(expectedPath, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
				defer cancel()

				target := server.URL + expectedPath
				for range expectedCount {
					err := Open(ctx, target, opts.Wait())
					require.NoError(t, err)
				}
				m.RLock()
				defer m.RUnlock()
				actualVisits := actualVisitsByPath[expectedPath]
				assert.Equal(t, expectedCount, actualVisits, "path count did not match")
			})
		}
	})

	t.Run("path visits - appname curl", func(t *testing.T) {
		expectedVisits := map[string]int{
			"/curl/foo":             1,
			"/curl/bar":             3,
			"/curl/baz":             3,
			"/curl/ctrl/alt/delete": 7,
			"/curl/none":            0,
		}
		for expectedPath, expectedCount := range expectedVisits {
			t.Run(expectedPath, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
				defer cancel()

				target := server.URL + expectedPath
				for range expectedCount {
					err := Open(ctx, target, opts.Wait(), opts.WithApp("curl"))
					require.NoError(t, err)
				}
				m.RLock()
				defer m.RUnlock()
				actualVisits := actualVisitsByPath[expectedPath]
				assert.Equal(t, expectedCount, actualVisits, "path count did not match")
			})
		}
	})
}
