package cloud

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v6"
)

type CloudStorage interface {
	Upload(ctx context.Context, bucketName string, byte []byte, objectName string) error
	Download(ctx context.Context, bucketName, objectName, destination string) error
}

type CloudStorageConf struct {
	StorageEndpoint        string
	AccessKeyID            string
	SecretAccessKey        string
	UseSSL                 bool
	UploadContextTimeout   time.Duration
	DownloadContextTimeout time.Duration
}

type cloudStorage struct {
	conf   *CloudStorageConf
	client *minio.Client
}

func NewCloudStorage(conf *CloudStorageConf) (CloudStorage, error) {

	if conf == nil {
		conf = getDefaultCloudStorageConf()
	}

	client, err := minio.New(conf.StorageEndpoint, conf.AccessKeyID, conf.SecretAccessKey, conf.UseSSL)
	if err != nil {
		return nil, err
	}

	return &cloudStorage{
		conf:   conf,
		client: client,
	}, nil
}

func getDefaultCloudStorageConf() *CloudStorageConf {
	conf := new(CloudStorageConf)
	conf.UploadContextTimeout = 10 * time.Minute
	conf.DownloadContextTimeout = 15 * time.Minute

	return conf
}

func (c *cloudStorage) Upload(ctx context.Context, bucketName string, byte []byte, objectName string) (err error) {
	_, err = c.client.PutObjectWithContext(ctx, bucketName, objectName, bytes.NewReader(byte), int64(len(byte)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})

	return err
}

func (c *cloudStorage) Download(ctx context.Context, bucketName, objectName, destination string) error {
	object, err := c.client.GetObjectWithContext(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(destination), os.ModePerm)
	if err != nil {
		return err
	}

	localObject, err := os.Create(destination)
	if err != nil {
		return err
	}

	if _, err = io.Copy(localObject, object); err != nil {
		return err
	}

	return nil
}
