/**
 * @author: D-S
 * @date: 2020/6/4 5:31 下午
 */

package main

import (
	"crypto/rsa"

	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const public_key = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA68z4jwTDQlr/qaN3SoJL
QPWL2rlZ57Yba1YERxpt5QFPKC4XasnyAKrNLtWgnb7AExrdOXflnggrKFI4YjsL
J07q8KBGrze78D+mvPju+26lIchsjMg8/5tgmVW2mYCTE0pgGhhXh7OfmcFTdYIP
iSZUT0WU0RWwmJ/br34O1wCzlvxAXGt4v6jyH/jUtsfu1D/uFfOyCGs9Fqt68Bhv
/jh1lCreAXAJkaUGVdzI2KNXe7MrVE5K0X/iDWbjAGFCIeSlhXs3NO1zyoTbVKMi
4Gib+mE60WL5r1RNyUYr/27K53RTtElaFfenDS0wWUf7w7G+t3yJckMNkqIv9YH4
VvAz3nJAbmeo/K3mvGFIgytBfy5NgUDYxuQwKe0rkgkRCLvlg81RMik1bwc8mGJ4
gc855olJ+uUkMPWO0gI84Im18svfX4sCSgES/hgnkCvboZaJHsAgSV2yboZnkkFl
PHAozXppWb+H79+67X7Eox1d5IN0IaqfVatewwx2vA/zVAn0Y4MZTvY+ppRsCGcd
cW57sbNDh1Oa8qT+bZTkWbACALV+Y0VnfC56xlB/XPo5LD2BTvd1AbyLEcA3x5q/
4h/lk8lWl8Z6vZMlA6aCUmk6dVVNTyM73B7EwIbOxfEjmIhK5r8tVX6aZcxigCOv
rlbVUOcK7z9X0ei9LIt7YPMCAwEAAQ==
-----END PUBLIC KEY-----`

const private_key = `-----BEGIN PRIVATE KEY-----
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
Tw9mvS2FAhEJLuAyAv8Z++dZ7j/iCtV3
-----END PRIVATE KEY-----
`

type PResp struct {
	Result struct {
		Items interface{} `json:"items"`
	} `json:"result"`
}

var args string

func main() {
	flag.StringVar(&args, "args", "", "")
	flag.Parse()
	sortS := sortStr(args)
	fmt.Println("请求排序----" + sortS)
	getSign := getSign(sortS)
	fmt.Println("签名结果----" + getSign)

	errVerify := VerifySign(sortS, getSign, public_key)
	if errVerify == nil {
		fmt.Println("签名成功---")
	}

}

func getTeamAll(gameId int64) bool {
	offset := 0
	for {
		str := "game_id=" + fmt.Sprintf("%v", gameId) + "&limit=100&offset=" + fmt.Sprintf("%v", offset)

		url := "sHttp://api.dawnbyte.com/api/v1/game/teams?" + sortStr(str)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("sign", getSign(sortStr(str)))

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		r := new(PResp)
		json.Unmarshal(body, &r)
		fmt.Println(r.Result.Items)
		if r.Result.Items == nil {
			return true
		}
		offset += 100
	}

}

func sortStr(args string) string {
	params := make(map[string]interface{}, 0)
	argsSlice := strings.Split(args, "&")
	for _, v := range argsSlice {
		param := strings.Split(v, "=")
		if len(param) != 2 {
			return ""
		}
		params[param[0]] = param[1]
	}
	params["request_time"] = time.Now().Unix()
	params["tenant_id"] = 6
	paramsStr := SortMap(params)

	return paramsStr
}

func getSign(sortStr string) string {
	sign, _ := Sign(sortStr, private_key)
	return sign
}

//签名
func Sign(name, privateKey string) (string, error) {
	name = MD5Hex(name)
	keyByte, _ := pem.Decode([]byte(privateKey))
	if keyByte == nil {
		return "", errors.New("private key error")
	}
	key, errKey := x509.ParsePKCS8PrivateKey(keyByte.Bytes)
	if errKey != nil {
		return "", errKey
	}
	h := crypto.SHA256.New()
	h.Write([]byte(name))
	hash := h.Sum(nil)

	signature, errSign := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hash)
	if errSign != nil {
		return "", errSign
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// 校验签名
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

func MD5Hex(val string) string {
	md5Obj := md5.New()
	md5Obj.Write([]byte(val[:]))
	return hex.EncodeToString(md5Obj.Sum(nil))
}
func SortMap(params map[string]interface{}) string {
	temp := make([]string, 0)
	for arg := range params {
		temp = append(temp, arg)
	}
	for i := 0; i < len(temp)-1; i++ {
		for j := 0; j < len(temp)-i-1; j++ {
			if temp[j] > temp[j+1] {
				temp[j], temp[j+1] = temp[j+1], temp[j]
			}
		}
	}
	for index, arg := range temp {
		temp[index] = fmt.Sprintf("%v=%v", arg, params[arg])
	}
	s := strings.Join(temp, "&")
	return s
}
