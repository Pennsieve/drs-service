// internal/utils/auth.go
package utils

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// ExtractBasicAuth 从请求中提取基本认证信息
func ExtractBasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", "", false
	}
	
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return "", "", false
	}
	
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}
	
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return "", "", false
	}
	
	return cs[:s], cs[s+1:], true
}

// ExtractBearerToken 从请求中提取Bearer令牌
func ExtractBearerToken(r *http.Request) (token string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", false
	}
	
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", false
	}
	
	return auth[len(prefix):], true
}
