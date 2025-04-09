package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("hd_psi_secret_key") // 在实际应用中应该从环境变量或配置文件中读取

// Claims 定义JWT令牌中包含的用户信息声明
// 包含用户ID、用户名和角色等基本信息，以及标准的JWT声明
type Claims struct {
	// UserID 用户唯一标识
	UserID uint `json:"user_id"`
	// Username 用户登录名
	Username string `json:"username"`
	// Role 用户角色
	Role string `json:"role"`
	// RegisteredClaims 包含标准JWT声明，如过期时间、签发时间等
	jwt.RegisteredClaims
}

// TokenExpiration 定义不同类型令牌的过期时间
const (
	// StandardTokenExpiration 标准令牌过期时间 - 24小时
	StandardTokenExpiration = 24 * time.Hour
	// RememberMeTokenExpiration 记住我令牌过期时间 - 7天
	RememberMeTokenExpiration = 7 * 24 * time.Hour
	// RefreshTokenExpiration 刷新令牌过期时间 - 30天
	RefreshTokenExpiration = 30 * 24 * time.Hour
	// PasswordResetTokenExpiration 密码重置令牌过期时间 - 1小时
	PasswordResetTokenExpiration = 1 * time.Hour
)

// GenerateToken 生成包含用户信息的JWT令牌
// 参数：
//   - userID: 用户唯一标识
//   - username: 用户登录名
//   - role: 用户角色
//   - rememberMe: 是否记住用户，如果为 true，则令牌过期时间更长
//
// 返回：
//   - string: 生成的JWT令牌字符串
//   - time.Time: 令牌过期时间
//   - error: 如果生成过程出错则返回错误，否则返回 nil
func GenerateToken(userID uint, username string, role string, rememberMe bool) (string, time.Time, error) {
	// 根据是否记住用户设置过期时间
	var expiration time.Duration
	if rememberMe {
		expiration = RememberMeTokenExpiration
	} else {
		expiration = StandardTokenExpiration
	}

	// 设置过期时间
	expirationTime := time.Now().Add(expiration)

	// 创建JWT声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hd_psi",
		},
	}

	// 使用声明创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ParseToken 解析和验证JWT令牌
// 参数：
//   - tokenString: 要解析的JWT令牌字符串
//
// 返回：
//   - *Claims: 如果令牌有效，返回包含用户信息的Claims对象
//   - error: 如果令牌无效或解析过程出错则返回错误
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
