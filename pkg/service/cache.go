package service

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
	"github.com/secmon-as-code/overseer/pkg/logging"
)

type gzipWriter struct {
	w  io.WriteCloser
	gz *gzip.Writer
}

func newGzipWriter(w io.WriteCloser) *gzipWriter {
	return &gzipWriter{w: w, gz: gzip.NewWriter(w)}
}

func (x *gzipWriter) Write(p []byte) (int, error) {
	return x.gz.Write(p)
}

func (x *gzipWriter) Close() error {
	rtErr := x.gz.Close()
	if rtErr != nil {
		logging.Default().Error("fail to close gzip writer", "err", rtErr)
	}

	if err := x.w.Close(); err != nil {
		rtErr = goerr.Wrap(err, "fail to close writer").With("err", rtErr)
	}

	return rtErr
}

type gzipReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create gzip reader")
	}
	return &gzipReader{r: r, gz: gz}, nil
}

func (x *gzipReader) Read(p []byte) (int, error) {
	return x.gz.Read(p)
}

func (x *gzipReader) Close() error {
	rtErr := x.gz.Close()
	if rtErr != nil {
		logging.Default().Error("fail to close gzip reader", "err", rtErr)
	}

	if err := x.r.Close(); err != nil {
		rtErr = goerr.Wrap(err, "fail to close reader").With("err", rtErr)
	}

	return rtErr
}

type FileCacheService struct {
	id      model.JobID
	baseDir string
}

func NewFileCacheService(id model.JobID, baseDir string) (*FileCacheService, error) {
	dirPath := filepath.Dir(fromIDtoFilePath(baseDir, id, "x"))
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return nil, goerr.Wrap(err, "fail to create baseDir for cache").With("baseDir", baseDir)
	}

	return &FileCacheService{id: id, baseDir: baseDir}, nil
}

func fromIDtoFilePath(baseDir string, jobID model.JobID, queryID model.QueryID) string {
	return filepath.Join(baseDir, string(jobID), string(queryID)+".json")
}

func (x *FileCacheService) NewWriter(_ context.Context, ID model.QueryID) (io.WriteCloser, error) {
	fpath := fromIDtoFilePath(x.baseDir, x.id, ID)

	fd, err := os.Create(filepath.Clean(fpath))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create file").With("path", fpath)
	}

	return newGzipWriter(fd), nil
}

func (x *FileCacheService) NewReader(_ context.Context, ID model.QueryID) (io.ReadCloser, error) {
	fpath := fromIDtoFilePath(x.baseDir, x.id, ID)
	println(fpath)

	fd, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to open file").With("path", fpath)
	}

	return newGzipReader(fd)
}

type CloudStorageCacheService struct {
	id     model.JobID
	bucket string
	prefix string
	client interfaces.CloudStorageClient
}

func fromIDtoCloudStoragePath(prefix string, jobID model.JobID, queryID model.QueryID) string {
	return strings.Join([]string{prefix, string(jobID), string(queryID) + ".json.gz"}, "/")
}

func NewCloudStorageCacheService(id model.JobID, bucket string, prefix string, client interfaces.CloudStorageClient) *CloudStorageCacheService {
	return &CloudStorageCacheService{
		id:     id,
		bucket: bucket,
		prefix: prefix,
		client: client,
	}
}

func (x *CloudStorageCacheService) NewWriter(ctx context.Context, ID model.QueryID) (io.WriteCloser, error) {
	w, err := x.client.PutObject(ctx, x.bucket, fromIDtoCloudStoragePath(x.prefix, x.id, ID))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create writer")
	}

	return newGzipWriter(w), nil
}

func (x *CloudStorageCacheService) NewReader(ctx context.Context, ID model.QueryID) (io.ReadCloser, error) {
	r, err := x.client.GetObject(ctx, x.bucket, fromIDtoCloudStoragePath(x.prefix, x.id, ID))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create reader")
	}

	return newGzipReader(r)
}
