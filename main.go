/**
* @Author: D-S
* @Date: 2020/3/20 10:31 上午
 */

package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"game-test/CSGO"
	"game-test/constant"
	jwtN "game-test/jwt"
	"game-test/library/log"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	//GetToken()
	//TestVerifySign()
	//engine := gin.Default()
	go CSGO.WS()
	//go CSGO.WS1()
	//go websocket.Ray57()
	//router.LoadRouter(engine)
	//engine.Run(":50052")
	//toTest()
	select {}
	//ToString()
}

func GetToken() {
	url := "https://api.abiosgaming.com/v2/oauth/access_token"

	payload := strings.NewReader("grant_type=client_credentials&client_id=chunzhi_c9425&client_secret=a345702d-5491-423e-8efe-71f00ca8a88d")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func TestVerifySign() {
	sign, errSign := Sign("game_id=1&limit=3000&offset=0&request_time=12123123&status=1&tenant_id=19", constant.Private_key19)
	errVerify := VerifySign("game_id=1&limit=3000&offset=0&request_time=12123123&status=1&tenant_id=19", sign, constant.Public_key19)

	fmt.Println(sign)
	fmt.Println(errSign)
	fmt.Println(errVerify)
	//
	//ass := assert.New(t)
	//ass.True(len(sign) > 0 && errVerify == nil)

} //620568336701.dkr.ecr.ap-southeast-1.amazonaws.com/risewinter/data-result-statistics:master-14-25aa0c5d058cb9b7b721f68c6693d8a6606a5da5
//620568336701.dkr.ecr.ap-southeast-1.amazonaws.com/risewinter/data-b-api:master-62-a8cf3d8eec4c211c53800f99464dd89213c15af3

func MD5Hex(val string) string {
	md5Obj := md5.New()
	md5Obj.Write([]byte(val[:]))
	return hex.EncodeToString(md5Obj.Sum(nil))
}

func Sign(val, privateKey string) (string, error) {
	signHex := MD5Hex(val)
	keyByte, _ := pem.Decode([]byte(privateKey))
	if keyByte == nil {
		return "", errors.New("private key error")
	}

	key, errKey := x509.ParsePKCS8PrivateKey(keyByte.Bytes)
	if errKey != nil {
		return "", errKey
	}
	h := crypto.SHA256.New()
	h.Write([]byte(signHex))
	hash := h.Sum(nil)

	signature, errSign := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hash)
	if errSign != nil {
		return "", errSign
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifySign(val, sign, publicKey string) error {
	signHex := MD5Hex(val)

	keyByte, _ := pem.Decode([]byte(publicKey))
	if keyByte == nil {
		return errors.New("public key error")
	}

	key, errKey := x509.ParsePKIXPublicKey(keyByte.Bytes)
	if errKey != nil {
		return errKey
	}

	h := crypto.SHA256.New()
	h.Write([]byte(signHex))
	hash := h.Sum(nil)

	signDecode, errDecode := base64.StdEncoding.DecodeString(sign)
	if errDecode != nil {
		return errDecode
	}

	return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), crypto.SHA256, hash, signDecode)
}

func token() {

	type UserInfo map[string]interface{}

	t := time.Now()
	key := "welcome to XXY's code world"
	userInfo := make(UserInfo)
	var expTime int64 = 1000 * 60 * 10
	var tokenState string
	now := strconv.FormatInt(t.UTC().UnixNano(), 10)
	fmt.Println(now)
	//userInfo["exp"] = "1515482650719371100" //
	userInfo["exp"] = now
	userInfo["iat"] = "0"

	tokenString := jwtN.CreateToken(key, userInfo)
	claims, ok := jwtN.ParseToken(tokenString, key)
	if ok {
		oldT, _ := strconv.ParseInt(claims.(jwt.MapClaims)["exp"].(string), 10, 64)
		ct := t.UTC().UnixNano()
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

func toTest() {
	now := time.Now()
	after := now.Add(time.Second * 30)
	nowS, _ := now.MarshalJSON()
	afterS, _ := after.MarshalJSON()
	log.Info(string(nowS))
	log.Info(string(afterS))
	if after.After(now.Add(time.Second * 31)) {
		log.Info("过期")
	} else {
		log.Info("没有过期")
	}
}

func ToString() {
	liveRate1 := new(LiveRate)
	liveRate1.Payload.Data.GroupId = 1232312312
	liveRate1.Payload.Data.ItemId = 234234234
	liveRate1.Payload.Data.SeriesId = 21123123

	b, _ := json.Marshal(liveRate1)
	s := string(b)
	log.Info(s)
	ss, _ := json.Marshal(s)
	log.Info(string(ss))

	_, err := toPushData(s, nil)
	if err != nil {
		log.Info(err.Error())
	}

}

func toPushData(from string, resultto interface{}) (*LiveRate, error) {
	result := new(LiveRate)
	//ttttttt, err := json.Marshal(from)
	//if err != nil {
	//	return nil, err
	//}
	err := json.Unmarshal([]byte(from), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type LiveRate struct {
	Channel string `json:"channel"`
	Payload struct {
		Source   string `json:"source"`
		PushTime int64  `json:"push_time"`
		Data     struct {
			GroupId    int64  `json:"group_id"`
			ItemId     int64  `json:"item_id"`
			SeriesId   int64  `json:"seriesId"`
			Status     int32  `json:"status"`
			Stage      int32  `json:"stage"`
			Rate       string `json:"rate"`
			From       int32  `json:"from"`
			UpdateTime int64  `json:"update_time"`
		} `json:"data"`
	} `json:"payload"`
}
