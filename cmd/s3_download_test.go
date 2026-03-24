package cmd_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/APTrust/apt-cmd/cmd"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// These must match the unexported constants in s3_download.go.
	chunkTestMin = int64(64 * 1024 * 1024)       // 64 MiB
	chunkTestMax = int64(5 * 1024 * 1024 * 1024) // 5 GiB
)

func TestComputeChunkSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected int64
	}{
		{name: "zero size", size: 0, expected: chunkTestMin},
		{name: "1 GiB", size: 1 * 1024 * 1024 * 1024, expected: chunkTestMin},
		{name: "640 GiB boundary", size: 10000 * chunkTestMin, expected: chunkTestMin},
		{name: "just over 640 GiB", size: 10000*chunkTestMin + 10000, expected: chunkTestMin + 1},
		{name: "2 TiB", size: 2 * 1024 * 1024 * 1024 * 1024, expected: 2 * 1024 * 1024 * 1024 * 1024 / 10000},
		{name: "100 TiB", size: 100 * 1024 * 1024 * 1024 * 1024, expected: chunkTestMax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, cmd.ComputeChunkSize(tt.size))
		})
	}
}

// newMockS3Server returns an httptest.Server that serves range GET requests
// from data and handles the minio client's bucket-location probes.
func newMockS3Server(data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.RawQuery, "location") {
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><LocationConstraint></LocationConstraint>")
			return
		}
		if r.Method == http.MethodHead {
			w.Header().Set("ETag", "\"test-etag\"")
			w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Accept-Ranges", "bytes")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet {
			rangeHeader := r.Header.Get("Range")
			if rangeHeader == "" {
				w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Header().Set("ETag", "\"test-etag\"")
				w.Header().Set("Last-Modified", "Mon, 15 Jan 2024 00:00:00 GMT")
				w.WriteHeader(http.StatusOK)
				w.Write(data) //nolint:errcheck
				return
			}
			var start, end int64
			fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
			dataLen := int64(len(data))
			if start >= dataLen {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			if end >= dataLen {
				end = dataLen - 1
			}
			chunk := data[start : end+1]
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, dataLen))
			w.Header().Set("Content-Length", strconv.Itoa(len(chunk)))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("ETag", "\"test-etag\"")
			w.Header().Set("Last-Modified", "Mon, 15 Jan 2024 00:00:00 GMT")
			w.WriteHeader(http.StatusPartialContent)
			w.Write(chunk) //nolint:errcheck
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}

// newTestMinioClient returns a minio client pointed at the given httptest server URL.
func newTestMinioClient(t *testing.T, serverURL string) *minio.Client {
	t.Helper()
	u, err := url.Parse(serverURL)
	require.NoError(t, err)
	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4("test-key", "test-secret", ""),
		Secure: false,
	})
	require.NoError(t, err)
	return client
}

func TestGetLargeObject_SmallObject(t *testing.T) {
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}
	server := newMockS3Server(data)
	defer server.Close()
	client := newTestMinioClient(t, server.URL)
	rc, err := cmd.GetLargeObject(client, "test-bucket", "test-key", int64(len(data)))
	require.NoError(t, err)
	require.NotNil(t, rc)
	defer rc.Close()
	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestGetLargeObject_MultiChunk(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large-object multi-chunk test in short mode")
	}
	// 2*minChunkSize+1 forces three separate range requests.
	size := 2*chunkTestMin + 1
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 251) // prime modulus avoids alignment with chunk boundaries
	}
	server := newMockS3Server(data)
	defer server.Close()
	client := newTestMinioClient(t, server.URL)
	rc, err := cmd.GetLargeObject(client, "test-bucket", "test-key", size)
	require.NoError(t, err)
	require.NotNil(t, rc)
	defer rc.Close()
	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestGetLargeObject_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.RawQuery, "location") {
			fmt.Fprint(w, "<?xml version=\"1.0\"?><LocationConstraint></LocationConstraint>")
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<?xml version=\"1.0\"?><Error><Code>InternalError</Code><Message>Internal Server Error</Message></Error>")
	}))
	defer server.Close()
	client := newTestMinioClient(t, server.URL)
	// GetLargeObject only starts a goroutine; errors surface when reading.
	rc, err := cmd.GetLargeObject(client, "test-bucket", "test-key", 100)
	require.NoError(t, err)
	require.NotNil(t, rc)
	defer rc.Close()
	_, readErr := io.ReadAll(rc)
	assert.Error(t, readErr)
}
