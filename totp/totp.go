package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// 提取常量
const (
	DefaultDigitCount = 6
	DefaultTimeStep   = 30
	DefaultT0         = 0
)

// TOTP 表示基于时间的一次性密码生成器
type TOTP struct {
	Secret     string
	DigitCount int
	TimeStep   uint64
	T0         uint64
}

// NewTOTP 创建一个新的TOTP实例，使用标准参数
func NewTOTP(secret string) *TOTP {
	return &TOTP{
		Secret:     secret,
		DigitCount: DefaultDigitCount,
		TimeStep:   DefaultTimeStep,
		T0:         DefaultT0,
	}
}

// GenerateCode generates a TOTP code based on the specified Unix time.
// It uses the TOTP instance's secret, time step, and digit count for calculation.
// Returns the generated TOTP code as a string or an error if generation fails.
func (t *TOTP) GenerateCode(unixTime uint64) (string, error) {
	tValue := calculateT(unixTime, t.T0, t.TimeStep)
	hash, err := computeHmacSha1(t.Secret, tValue)
	if err != nil {
		return "", err
	}
	return truncate(hash, t.DigitCount), nil
}

// calculateT T totp's first t
func calculateT(unixTime, t0, timeStep uint64) uint64 {
	return uint64(float64(unixTime-t0) / float64(timeStep))
}

// computeHmacSha1 generates an HMAC using the SHA-1 hashing algorithm with the given key and message.
// It takes a base32-encoded key string and a uint64 message as input.
// It returns the computed HMAC as a byte slice or an error if key decoding fails.
func computeHmacSha1(key string, msg uint64) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("the key cannot be empty")
	}
	upperedLetterKey := strings.ToUpper(key)
	keyByte, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(upperedLetterKey)
	if err != nil {
		return nil, fmt.Errorf("解码base32密钥失败: %w", err)
	}

	msgByte := make([]byte, 8)
	binary.BigEndian.PutUint64(msgByte, msg)

	hash := hmac.New(sha1.New, keyByte)
	hash.Write(msgByte)

	return hash.Sum(nil), nil
}

// truncate extracts a numeric code from a hashed byte sequence and formats it with leading zeros to match digitCount.
func truncate(hash []byte, digitCount int) string {
	// 计算距离哈希最后一个字节的偏移量
	offset := uintptr(hash[len(hash)-1] & 0xf)

	// 从计算出的偏移量开始从哈希中提取 4 个字节
	extractedValue := int(hash[offset])<<24 | int(hash[offset+1])<<16 | int(hash[offset+2])<<8 | int(hash[offset+3])

	// 通过按位与运算确保提取的值为正数
	extractedPositiveValue := int64(extractedValue & 0x7fffffff)

	otp := int(int64(extractedPositiveValue) % int64(math.Pow10(digitCount)))

	// 使用前导零格式化 OTP 以适合数字计数
	formattedOTP := fmt.Sprintf(fmt.Sprintf("%%0%dd", digitCount), otp)

	return formattedOTP
}
