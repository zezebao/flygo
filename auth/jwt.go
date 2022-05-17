//认证相关逻辑
//处理用户名、有效期、校验等逻辑...
package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/zezebao/flygo/logger"
	"go.uber.org/zap"
)

var (
	key []byte = []byte("-jwt-lkpush.com")
)

var Jwt *jwtStruct = &jwtStruct{}

type jwtStruct struct{}

// 产生json web token
func (j *jwtStruct) GenToken(c *gin.Context, userName string) string {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + 60*30), //seconds
		Issuer:    userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logger.ErrLogger.Error("auth.token", zap.String("err", err.Error()))
		return ""
	}
	return ss
}

// 校验token是否有效
func (j *jwtStruct) CheckToken(token string) bool {
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false
	}
	return true
}
