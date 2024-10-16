package service_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-as-code/overseer/pkg/adaptor/cs"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
	"github.com/secmon-as-code/overseer/pkg/service"
)

func TestCacheFile(t *testing.T) {
	d := os.TempDir()
	id1, id2 := model.NewJobID(), model.NewJobID()
	svc1, err := service.NewFileCacheService(id1, d)
	gt.NoError(t, err)
	svc2, err := service.NewFileCacheService(id2, d)
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func TestCacheCloudStorage(t *testing.T) {
	bucketName, ok := os.LookupEnv("TEST_CLOUD_STORAGE_BUCKET_NAME")
	if !ok {
		t.Skip("TEST_CLOUD_STORAGE_BUCKET_NAME is not set")
	}

	client, err := cs.NewClient(context.Background())

	id1, id2 := model.NewJobID(), model.NewJobID()
	svc1 := service.NewCloudStorageCacheService(id1, bucketName, "overseer-test", client)
	gt.NoError(t, err)
	svc2 := service.NewCloudStorageCacheService(id2, bucketName, "overseer-test", client)
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func testCache(t *testing.T, svc1, svc2 interfaces.CacheService) {
	ctx := context.Background()
	t.Run("write data to cache", func(t *testing.T) {
		w, err := svc1.NewWriter(ctx, "test1")
		gt.NoError(t, err)

		_, err = w.Write([]byte("blue"))
		gt.NoError(t, err)
		gt.NoError(t, w.Close())
	})

	t.Run("read data from cache", func(t *testing.T) {
		r, err := svc1.NewReader(ctx, "test1")
		gt.NoError(t, err)

		buf, err := io.ReadAll(r)
		gt.NoError(t, err)
		gt.Equal(t, string(buf), "blue")
	})

	t.Run("can not read data from other id", func(t *testing.T) {
		r, err := svc1.NewReader(ctx, "test2")
		gt.Error(t, err)
		gt.Equal(t, r, nil)
	})

	t.Run("can not read data from other job", func(t *testing.T) {
		r, err := svc2.NewReader(ctx, "test1")
		gt.Error(t, err)
		gt.Equal(t, r, nil)
	})
}
