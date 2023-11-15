package testutils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

// allow alphanumerics, underscores and dashes
var testNameCleanRegex = regexp.MustCompile(`[^a-zA-Z0-9 \-_]+`)

func getRaftStateTraceLocation() string {
	if v, enabled := os.LookupEnv("TEST_TRACE_RAFT_STATE"); enabled && len(v) > 0 {
		return v
	}

	return ""
}

func GetRaftStateTraceFilename(t testing.TB) string {
	location := getRaftStateTraceLocation()
	if len(location) == 0 {
		return ""
	}

	filename := fmt.Sprintf("%s-%s.ndjson", testNameCleanRegex.ReplaceAllString(t.Name(), ""), time.Now().Format("20060102150405"))
	return filepath.Join(location, filename)
}

func CompressRaftStateTraceFile(filename string) error {
	if len(filename) == 0 {
		return nil
	}

	zipFilename := fmt.Sprintf("%s.gz", filename)
	writer, err := os.OpenFile(zipFilename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer writer.Close()

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	zip := gzip.NewWriter(writer)
	defer zip.Close()

	_, err = io.Copy(zip, reader)
	if err != nil {
		return err
	}

	zip.Flush()
	os.Remove(filename)

	return nil
}
