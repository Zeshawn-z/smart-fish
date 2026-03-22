package v1

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
	"time"

	"smart-fish/back_end/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// generateFlaskCompatToken 生成兼容 Flask flask-jwt-extended 的 JWT
// Flask 的 JWT payload 中 sub = user_id (数字)
func generateFlaskCompatToken(userID uint) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":     fmt.Sprintf("%d", userID), // Flask 把 identity 存在 sub 中
		"user_id": userID,                    // 同时添加 user_id 字段便于 Go 端解析
		"iat":     now.Unix(),
		"exp":     now.Add(24 * time.Hour).Unix(), // 默认 24 小时
		"type":    "access",
		"fresh":   false,
		"nbf":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// hashPassword 使用 bcrypt 哈希密码
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword 验证密码，支持 bcrypt 和 werkzeug 两种格式
func checkPassword(hashed, password string) bool {
	// 尝试 bcrypt（Go 后端格式）
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err == nil {
		return true
	}

	// 尝试 werkzeug 格式: "method$salt$hash"
	// werkzeug 使用 pbkdf2:sha256 等方法
	if strings.Contains(hashed, "$") {
		parts := strings.SplitN(hashed, "$", 3)
		if len(parts) == 3 {
			method := parts[0]
			if strings.HasPrefix(method, "pbkdf2:") || strings.HasPrefix(method, "scrypt:") {
				// 需要使用对应的算法验证
				// 由于 werkzeug 的哈希格式较复杂，这里简化处理
				// 实际生产中可能需要引入 Python 的 werkzeug 格式解析
				return checkWerkzeugPassword(hashed, password)
			}
		}
	}

	return false
}

// checkWerkzeugPassword 验证 werkzeug 格式的密码哈希
// werkzeug 格式: "pbkdf2:sha256:600000$salt$hash" 或 "scrypt:32768:8:1$salt$hash"
func checkWerkzeugPassword(hashed, password string) bool {
	parts := strings.SplitN(hashed, "$", 3)
	if len(parts) != 3 {
		return false
	}

	method := parts[0]
	salt := parts[1]
	expectedHash := parts[2]

	// 解析方法
	methodParts := strings.Split(method, ":")
	if len(methodParts) < 2 {
		return false
	}

	switch methodParts[0] {
	case "pbkdf2":
		hashFunc := methodParts[1] // sha256, sha512, etc.
		iterations := 600000       // 默认迭代次数
		if len(methodParts) >= 3 {
			fmt.Sscanf(methodParts[2], "%d", &iterations)
		}
		actualHash := pbkdf2Hash(password, salt, iterations, hashFunc)
		return actualHash == expectedHash

	case "scrypt":
		// scrypt:N:r:p$salt$hash
		// Werkzeug 3.x 默认: scrypt:32768:8:1
		n := 32768
		r := 8
		p := 1
		if len(methodParts) >= 2 {
			fmt.Sscanf(methodParts[1], "%d", &n)
		}
		if len(methodParts) >= 3 {
			fmt.Sscanf(methodParts[2], "%d", &r)
		}
		if len(methodParts) >= 4 {
			fmt.Sscanf(methodParts[3], "%d", &p)
		}
		actualHash := scryptHash(password, salt, n, r, p)
		return actualHash == expectedHash
	}

	return false
}

// pbkdf2Hash 使用 PBKDF2 算法生成哈希（兼容 werkzeug 格式）
// werkzeug 的哈希输出为小写十六进制字符串
func pbkdf2Hash(password, salt string, iterations int, hashFunc string) string {
	var keyLen int
	var h func() hash.Hash

	switch hashFunc {
	case "sha256":
		keyLen = 32
		h = sha256.New
	case "sha512":
		keyLen = 64
		h = sha512.New
	default:
		keyLen = 32
		h = sha256.New
	}

	dk := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLen, h)
	return fmt.Sprintf("%x", dk)
}

// scryptHash 使用 scrypt 算法生成哈希（兼容 werkzeug 3.x 格式）
// werkzeug scrypt 格式: "scrypt:N:r:p$salt$hex_hash"
// 默认参数: N=32768, r=8, p=1, keyLen=64
func scryptHash(password, salt string, n, r, p int) string {
	dk, err := scrypt.Key([]byte(password), []byte(salt), n, r, p, 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", dk)
}

// parseFlexibleTime 解析多种时间格式（兼容 Flask datetime 传参）
func parseFlexibleTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+08:00",
		time.RFC3339,
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse time: %s", s)
}

// parseUintParam 从路径参数字符串解析 uint
func parseUintParam(s string) (uint, error) {
	var n uint
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("invalid uint: %s", s)
		}
		n = n*10 + uint(ch-'0')
	}
	return n, nil
}
