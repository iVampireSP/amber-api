package file

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/codingsince1985/checksum"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"os"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
	"strconv"
	"time"
)

const TMPDIR = ""
const RootDir = "files"
const MaxSize = 5 * 1024 * 1024

func (s *Service) CreateFileFromUrl(ctx context.Context, url string) (*entity.File, error) {
	var urlHash = s.sha256String(url)

	fileEntity := &entity.File{}

	// 如果已经存在
	exists, err := s.URLExists(ctx, urlHash)
	if err != nil {
		return nil, err
	}
	if exists {
		fileEntity, err = s.GetFileByUrlHash(ctx, urlHash)
		return fileEntity, err
	}

	fileEntity.Url = &url
	fileEntity.UrlHash = &urlHash

	// 获取内容
	// 先验证大小
	err = s.validateRemoteFileSize(url)
	if err != nil {
		return nil, err
	}

	// 下载文件
	path, err := s.downloadRemoteFile(url)
	if err != nil {
		return nil, err
	}

	fileSha256, err := checksum.SHA256sum(path)
	if err != nil {
		return nil, err
	}

	fileMimeType, err := mimetype.DetectFile(path)
	if err != nil {
		return nil, consts.ErrMimeTypeNotFound
	}
	fileMimeTypeString := fileMimeType.String()

	// 只允许 jpg(jpeg),png,webp
	if fileMimeTypeString != "image/jpeg" && fileMimeTypeString != "image/png" && fileMimeTypeString != "image/webp" {
		return nil, consts.ErrMimeTypeNotAllowed
	}

	// 上传文件到 S3
	fileName, filePath := s.GenerateFilePath(fileSha256)

	err = s.uploadToBucket(ctx, filePath+"/"+fileName, path)
	if err != nil {
		return nil, err
	}

	// 全部成功的情况
	fileEntity.FileHash = fileSha256
	fileEntity.MimeType = fileMimeType.String()
	fileEntity.Path = filePath + "/" + fileName

	// 删除临时文件
	defer s.deleteTmpFile(path)

	_, err = s.x.Insert(fileEntity)
	if err != nil {
		return nil, err
	}

	return fileEntity, nil
}

func (s *Service) CreateFile(ctx context.Context, file io.ReadSeekCloser) (*entity.File, error) {
	fileEntity := &entity.File{
		Url:     nil,
		UrlHash: nil,
	}

	fileSha256, err := checksum.MD5sumReader(file)
	if err != nil {
		return nil, err
	}

	// 如果已经存在
	exists, err := s.FileHashExists(ctx, fileSha256)
	if err != nil {
		return nil, err
	}
	if exists {
		fileEntity, err = s.GetFileByFileHash(ctx, fileSha256)
		return fileEntity, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// 获取路径
	fileMimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, consts.ErrMimeTypeNotFound
	}
	fileMimeTypeString := fileMimeType.String()

	// 只允许 jpg(jpeg),png,webp
	if fileMimeTypeString != "image/jpeg" && fileMimeTypeString != "image/png" && fileMimeTypeString != "image/webp" {
		return nil, consts.ErrMimeTypeNotAllowed
	}

	// 上传文件到 S3
	fileName, filePath := s.GenerateFilePath(fileSha256)

	err = s.uploadToBucketIO(ctx, filePath+"/"+fileName, file)
	if err != nil {
		return nil, err
	}

	// 全部成功的情况
	fileEntity.FileHash = fileSha256
	fileEntity.MimeType = fileMimeType.String()
	fileEntity.Path = filePath + "/" + fileName

	_, err = s.x.Insert(fileEntity)
	if err != nil {
		return nil, err
	}

	return fileEntity, nil
}

func (s *Service) uploadToBucket(ctx context.Context, filename string, localPath string) error {
	_, err := s.s3.Client.FPutObject(ctx, s.s3.Bucket, filename, localPath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) uploadToBucketIO(ctx context.Context, filename string, file io.ReadSeeker) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	size, err := io.Copy(io.Discard, file)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = s.s3.Client.PutObject(ctx, s.s3.Bucket, filename, file, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GenerateFilePath(sha256FileName string) (filename string, path string) {
	// 取前 2 个字符
	p2 := sha256FileName[0:2]
	// 取第 3 到 4 个字符
	p34 := sha256FileName[2:4]

	// 生成路径
	return sha256FileName, fmt.Sprintf("/%s/%s/%s", RootDir, p2, p34)
}

// GetBucketFile 获取 S3 中的文件，并返回一个 io.ReadCloser
func (s *Service) GetBucketFile(ctx context.Context, fileEntity *entity.File) (size int64, io io.ReadCloser, err error) {
	obj, err := s.s3.Client.GetObject(ctx, s.s3.Bucket, fileEntity.Path, minio.GetObjectOptions{})
	if err != nil {
		return 0, nil, err
	}
	stat, err := obj.Stat()

	if err != nil {
		return 0, nil, err
	}

	return stat.Size, obj, nil
}

func (s *Service) downloadRemoteFile(url string) (path string, err error) {
	// 随机一个文件名
	var filename = strconv.FormatInt(time.Now().UnixNano(), 10)

	// 保存到系统临时文件
	file, err := os.CreateTemp(TMPDIR, filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// 下载文件
	rsp, err := http.Get(url)
	defer func() {
		_ = rsp.Body.Close()
	}()
	if err != nil {
		return "", err
	}

	// 限制读取的最大长度
	limitedReader := io.LimitedReader{R: rsp.Body, N: MaxSize}
	_, err = io.Copy(file, &limitedReader)
	if err != nil {
		// 如果读取长度超过 MaxSize，则直接关闭文件和响应体
		if errors.Is(err, io.ErrShortWrite) {
			_ = file.Close()
			_ = rsp.Body.Close()
			return "", consts.ErrFileSizeTooLarge
		}
		return "", err
	}

	// 获取文件路径
	path = file.Name()

	return path, nil
}

// ValidateRemoteFileSize 验证远程文件大小，如果超过限制则返回错误
func (s *Service) validateRemoteFileSize(url string) error {
	// Create a new HTTP client
	client := http.Client{
		Timeout: time.Second * 10,
	}

	// Create a new HTTP request
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	// Send the HTTP request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Get the content length from the response headers
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return consts.ErrContentLengthHeaderMissing
	}

	// Convert the content length to an integer
	fileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return err
	}

	// Check if the file size exceeds the limit
	if fileSize > MaxSize {
		return consts.ErrFileSizeTooLarge
	}

	return nil
}

// delete tmp file
func (s *Service) deleteTmpFile(path string) {
	_ = os.Remove(path)
}

func (s *Service) sha256String(str string) string {
	h := sha256.New()
	h.Write([]byte(str))

	return fmt.Sprintf("%x", h.Sum(nil))
}
