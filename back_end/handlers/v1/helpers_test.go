package v1

import (
	"testing"
)

func TestCheckWerkzeugPassword(t *testing.T) {
	// 由 Python werkzeug.security.generate_password_hash('test123', method='pbkdf2:sha256') 生成
	werkzeugHash := "pbkdf2:sha256:1000000$pRWshiUTIcfFPWH7$6c2c343d9040812137169b6115dfba15c002d1d916434d123187ca9c82c8d697"

	if !checkWerkzeugPassword(werkzeugHash, "test123") {
		t.Error("checkWerkzeugPassword should return true for correct password")
	}

	if checkWerkzeugPassword(werkzeugHash, "wrong_password") {
		t.Error("checkWerkzeugPassword should return false for wrong password")
	}
}

func TestCheckPassword_Werkzeug(t *testing.T) {
	// 同时测试 checkPassword 分发逻辑
	werkzeugHash := "pbkdf2:sha256:1000000$pRWshiUTIcfFPWH7$6c2c343d9040812137169b6115dfba15c002d1d916434d123187ca9c82c8d697"

	if !checkPassword(werkzeugHash, "test123") {
		t.Error("checkPassword should support werkzeug pbkdf2 format")
	}
}

func TestCheckPassword_Bcrypt(t *testing.T) {
	// 测试 bcrypt 路径
	hash, err := hashPassword("hello")
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	if !checkPassword(hash, "hello") {
		t.Error("checkPassword should work with bcrypt hash")
	}

	if checkPassword(hash, "wrong") {
		t.Error("checkPassword should reject wrong password for bcrypt")
	}
}

func TestPbkdf2Hash(t *testing.T) {
	// 验证 pbkdf2Hash 输出格式
	// werkzeug 格式: pbkdf2:sha256:iterations$salt$hex_hash
	// 手动验证：password="test123", salt="pRWshiUTIcfFPWH7", iterations=1000000, sha256
	expected := "6c2c343d9040812137169b6115dfba15c002d1d916434d123187ca9c82c8d697"
	actual := pbkdf2Hash("test123", "pRWshiUTIcfFPWH7", 1000000, "sha256")
	if actual != expected {
		t.Errorf("pbkdf2Hash mismatch:\n  expected: %s\n  actual:   %s", expected, actual)
	}
}

func TestCheckPassword_Scrypt(t *testing.T) {
	// Werkzeug 3.x 默认使用 scrypt 格式: "scrypt:N:r:p$salt$hash"
	// 由 Python werkzeug 3.0.6 generate_password_hash('test123') 生成
	scryptHash := "scrypt:32768:8:1$1gvr7PSyGGBUywLv$61cab99f212e857c374824081fd3fc3047dbcb84a553c4cb60e6402ce78bc939762d73fedd3eafe592546202c378582da6569c2aa57760963f1054cbc2f2916a"

	if !checkPassword(scryptHash, "test123") {
		t.Error("checkPassword should support werkzeug scrypt format")
	}

	if checkPassword(scryptHash, "wrong_password") {
		t.Error("checkPassword should reject wrong password for scrypt")
	}
}
