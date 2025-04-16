// internal/utils/checksum.go
package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"sort"
	"strings"

	model "github.com/pennsieve/drs-service/internal/models"
)

// 校验和类型常量
const (
	ChecksumTypeMD5    = "md5"
	ChecksumTypeSHA1   = "sha1"
	ChecksumTypeSHA256 = "sha-256"
	ChecksumTypeSHA512 = "sha-512"
)

// CalculateChecksum 计算数据的校验和
func CalculateChecksum(r io.Reader, checksumType string) (string, error) {
	var hasher hash.Hash
	
	switch strings.ToLower(checksumType) {
	case ChecksumTypeMD5:
		hasher = md5.New()
	case ChecksumTypeSHA1:
		hasher = sha1.New()
	case ChecksumTypeSHA256:
		hasher = sha256.New()
	case ChecksumTypeSHA512:
		hasher = sha512.New()
	default:
		return "", fmt.Errorf("不支持的校验和类型: %s", checksumType)
	}
	
	if _, err := io.Copy(hasher, r); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// CalculateMultipleChecksums 计算数据的多种校验和
func CalculateMultipleChecksums(r io.Reader, checksumTypes []string) ([]model.Checksum, error) {
	// 需要读取所有数据以计算多个校验和
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	
	checksums := make([]model.Checksum, 0, len(checksumTypes))
	
	for _, checksumType := range checksumTypes {
		checksum, err := CalculateChecksum(strings.NewReader(string(data)), checksumType)
		if err != nil {
			return nil, err
		}
		
		checksums = append(checksums, model.Checksum{
			Type:     checksumType,
			Checksum: checksum,
		})
	}
	
	return checksums, nil
}

// CalculateBundleChecksum 计算包含多个对象的校验和
// 按照DRS规范，这是通过对所有对象的校验和进行排序、连接，然后再计算校验和实现的
func CalculateBundleChecksum(checksums []model.Checksum, checksumType string) (string, error) {
	// 提取所有校验和值
	checksumValues := make([]string, 0, len(checksums))
	for _, checksum := range checksums {
		checksumValues = append(checksumValues, checksum.Checksum)
	}
	
	// 排序
	sort.Strings(checksumValues)
	
	// 连接
	concatenated := strings.Join(checksumValues, "")
	
	// 计算新的校验和
	return CalculateChecksum(strings.NewReader(concatenated), checksumType)
}
