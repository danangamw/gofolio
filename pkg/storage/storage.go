package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	s3Client         *s3.Client
	bucketName       string
	internalEndpoint string
	publicEndpoint   string
)

// InitStorage menginisialisasi client dan bucket S3 sekali saja saat aplikasi startup.
func InitStorage(client *s3.Client, bucket string, internalEP, publicEP string) {
	s3Client = client
	bucketName = bucket
	internalEndpoint = internalEP
	publicEndpoint = publicEP
}

// InitStorageFromConfig menginisialisasi S3/MinIO client menggunakan konfigurasi mentah.
func InitStorageFromConfig(endpoint, publicEP, accessKey, secretKey, bucket, region string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	// Buat client S3 dengan custom endpoint (MinIO)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	// Cek apakah bucket sudah ada, jika belum buat baru secara otomatis.
	// Menggunakan retry loop untuk mengantisipasi jika container MinIO masih dalam proses booting.
	var lastErr error
	for i := 1; i <= 10; i++ {
		checkCtx, cancelCheck := context.WithTimeout(ctx, 2*time.Second)
		_, err = client.HeadBucket(checkCtx, &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		cancelCheck()

		if err == nil {
			InitStorage(client, bucket, endpoint, publicEP)
			return nil
		}

		// Coba buat bucket baru
		createCtx, cancelCreate := context.WithTimeout(ctx, 2*time.Second)
		_, createErr := client.CreateBucket(createCtx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		cancelCreate()

		if createErr == nil {
			InitStorage(client, bucket, endpoint, publicEP)
			return nil
		}

		lastErr = createErr
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to auto-create S3/MinIO bucket %s after retries: %w", bucket, lastErr)
}

// UploadFile mengunggah file ke S3/MinIO bucket menggunakan client global yang sudah terinisialisasi.
func UploadFile(ctx context.Context, file io.Reader, folder string, customName string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("storage client has not been initialized. Call storage.InitStorage() first")
	}

	fileName := customName
	if fileName == "" {
		fileName = fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	}

	key := filepath.Join(folder, fileName)

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

// DeleteFile menghapus file dari S3/MinIO bucket.
func DeleteFile(ctx context.Context, filePath string) error {
	if s3Client == nil {
		return fmt.Errorf("storage client has not been initialized. Call storage.InitStorage() first")
	}

	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	})
	return err
}

// DeleteFiles menghapus banyak file dari S3/MinIO bucket.
func DeleteFiles(ctx context.Context, filePaths []string) error {
	for _, path := range filePaths {
		_ = DeleteFile(ctx, path)
	}
	return nil
}

// GetPresignedURL menghasilkan temporary signed URL dari S3/MinIO untuk file path/key tertentu.
// Masa kedaluwarsa URL ditentukan oleh parameter expires.
func GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("storage client has not been initialized. Call storage.InitStorage() first")
	}

	if key == "" {
		return "", nil
	}

	presignClient := s3.NewPresignClient(s3Client)
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	urlStr := presignedReq.URL
	// Jika publicEndpoint di-set dan berbeda dengan internalEndpoint, ganti host & scheme-nya
	if publicEndpoint != "" && internalEndpoint != "" {
		urlStr = strings.Replace(urlStr, internalEndpoint, publicEndpoint, 1)
	}

	return urlStr, nil
}

// GetPublicURL menghasilkan URL publik langsung (non-presigned) jika bucket terkonfigurasi publik.
func GetPublicURL(key string) string {
	if key == "" {
		return ""
	}

	base := publicEndpoint
	if base == "" {
		base = internalEndpoint
	}

	base = strings.TrimSuffix(base, "/")
	// Pastikan jika reverse-proxy langsung ke bucket, key ditambahkan langsung.
	// Jika host adalah domain global MinIO server biasa, tambahkan nama bucket.
	// Kita serahkan format publicEndpoint agar fleksibel, jadi jika publicEndpoint berisi path bucket (misal http://localhost:9000/go-cms),
	// key akan diappend setelahnya.
	// Jika tidak mengandung nama bucket dan kita tidak memakai subdomain per-bucket, defaultnya:
	if !strings.Contains(base, bucketName) && !strings.Contains(base, "danangamw.com") {
		// fallback standard path-style
		return fmt.Sprintf("%s/%s/%s", base, bucketName, key)
	}

	return fmt.Sprintf("%s/%s", base, key)
}
