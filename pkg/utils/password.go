package utils

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

const saltCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// 生成随机6位数 盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, 6)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	for i := range salt {
		salt[i] = saltCharset[salt[i]%byte(len(saltCharset))]
	}

	return string(salt), nil
}

// 使用 盐值 加密密码
func HashPasswordWithSalt(password string, salt string) (string, error) {
	saltedPassword := password + salt

	// 使用 单向hash bcrypt 算法加密  默认难度为 DefaultCost: 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 校验 密码 + 盐值
func CheckPasswordWithSalt(password string, hashedPassword string, salt string) bool {

	saltedPassword := password + salt
	// 使用 bcrypt 验证
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))

	return err == nil
}
