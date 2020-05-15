/**
 * @author: D-S
 * @date: 2020/5/12 2:13 下午
 */

package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var (
	key = "xxxxxxxxxxxxx"
)

func TestToken() {
	type Info map[string]interface{}
	var expTime int64 = 1000 * 60 * 10
	var tokenState string
	info := make(map[string]interface{})
	now := strconv.FormatInt(time.Now().Local().UnixNano(), 10)
	info["exp"] = now
	info["iat"] = "0"

	tokenString := CreateToken(key, info)
	claims, ok := ParseToken(tokenString, key)
	if ok {
		oldT, _ := strconv.ParseInt(claims.(jwt.MapClaims)["exp"].(string), 10, 64)
		ct := time.Now().Local().UnixNano()
		c := ct - oldT
		fmt.Println(ct)
		if c > expTime {
			ok = false
			tokenState = "Token 已过期"
		} else {
			tokenState = "Token 正常"
		}
	} else {
		tokenState = "token无效"
	}
	fmt.Println("======")
	fmt.Println(tokenString)
	fmt.Println("======")
	fmt.Println(tokenState)
	fmt.Println(claims)

}

func CreateToken(key string, m map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	for index, val := range m {
		claims[index] = val
	}
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}

func ParseToken(tokenString string, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return "", false
	}
}
