package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	TotpPeriod = 30 // 每 30 秒一个时间窗口
	totpDigits = 6
	totpSkew   = 1 // 允许前后 ±1 个窗口（共 90 秒容差）
)

// GenerateECCKey 生成新的 ECDSA P-256 私钥，返回 base64 编码的私钥标量。
func GenerateECCKey() (string, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key.D.Bytes()), nil
}

// ParseECCKey 解码 base64 编码的 ECC 私钥标量。
func ParseECCKey(b64 string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 解码失败: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("密钥为空")
	}
	return data, nil
}

// GenerateTOTP 根据 ECC 私钥和当前时间生成 6 位动态码。
func GenerateTOTP(key []byte, t time.Time) string {
	counter := t.Unix() / TotpPeriod
	return computeHOTP(key, uint64(counter))
}

// VerifyTOTP 校验 6 位动态码，允许前后各 1 个时间窗口。
func VerifyTOTP(key []byte, code string, t time.Time) bool {
	counter := t.Unix() / TotpPeriod
	for i := int64(-totpSkew); i <= int64(totpSkew); i++ {
		c := counter + i
		if c < 0 {
			continue
		}
		if computeHOTP(key, uint64(c)) == code {
			return true
		}
	}
	return false
}

// computeHOTP 基于 HMAC-SHA256 + 动态截断生成 6 位码（RFC 4226 变体）。
func computeHOTP(key []byte, counter uint64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	mac := hmac.New(sha256.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	binCode := (uint32(sum[offset])&0x7f)<<24 |
		uint32(sum[offset+1])<<16 |
		uint32(sum[offset+2])<<8 |
		uint32(sum[offset+3])

	otp := binCode % 1000000
	return fmt.Sprintf("%06d", otp)
}
