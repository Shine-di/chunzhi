/**
 * @author: D-S
 * @date: 2020/5/12 2:12 下午
 */

package auth

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	params := os.Args
	var param string
	var privateKey string
	fmt.Println(len(params))
	if len(params) < 1 {
		fmt.Println("参数不够")
		return
	}
	param = params[1]
	if len(params) >= 3 {
		privateKey = params[2]
	} else {
		privateKey = ""
	}
	param, privateKey = SortParam(param, privateKey)
	sign, err := Sign(param, privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sign)
	fmt.Println("-------------")
	fmt.Println(param)
}

func DoVerifySign(param, privateKey, publicKey string) {
	s, privateKey := SortParam(param, privateKey)
	sign, errSign := Sign(s, privateKey)
	errVerify := VerifySign(s, sign, publicKey)
	fmt.Println(sign)
	fmt.Println(errSign)
	fmt.Println(errVerify)
	fmt.Println("----------------------")
	fmt.Println(s)
}

func SortParam(param, privateKey string) (string, string) {
	var params []string
	if privateKey == "" {
		for _, e := range strings.Split(param, "&") {
			if strings.HasPrefix(e, "tenant_id") {
				keys := strings.Split(e, "=")
				fmt.Println("私钥为空使用商户6")
				keys[1] = "6"
				params = append(params, "tenant_id=6")
			} else {
				params = append(params, e)
			}
		}
		privateKey = private_key6
	} else {
		params = strings.Split(param, "&")
	}
	result := make([]string, 0)
	for _, param := range params {
		if strings.HasPrefix(param, "time_stamp") {
			t := time.Now().Unix()
			s := strconv.Itoa(int(t))
			result = append(result, "time_stamp="+s)
			//result = append(result, param)
		} else {
			result = append(result, param)
		}
	}
	sort.Strings(result)
	return strings.Join(result, "&"), privateKey
}

func SortParamMap(param map[string]string, privateKey string) (string, string) {
	var params []string
	for k, v := range param {
		item := k + "=" + v
		params = append(params, item)
	}
	t := time.Now().Unix()
	s := strconv.Itoa(int(t))
	params = append(params, "time_stamp="+s)
	sort.Strings(params)
	return strings.Join(params, "&"), privateKey
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
	val, _ = SortParam(val, "")
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

////签名
//func Sign(name, privateKey string) (string, error) {
//	name = MD5Hex(name)
//	keyByte, _ := pem.Decode([]byte(privateKey))
//	if keyByte == nil {
//		return "", errors.New("private key error")
//	}
//	key, errKey := x509.ParsePKCS8PrivateKey(keyByte.Bytes)
//	if errKey != nil {
//		return "", errKey
//	}
//	h := crypto.SHA256.New()
//	h.Write([]byte(name))
//	hash := h.Sum(nil)
//
//	signature, errSign := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hash)
//	if errSign != nil {
//		return "", errSign
//	}
//	return base64.StdEncoding.EncodeToString(signature), nil
//}
//
//// 校验签名
//func VerifySign(val, sign, publicKey string) error {
//	signHex := MD5Hex(val)
//
//	keyByte, _ := pem.Decode([]byte(publicKey))
//	if keyByte == nil {
//		return errors.New("public key error")
//	}
//
//	key, errKey := x509.ParsePKIXPublicKey(keyByte.Bytes)
//	if errKey != nil {
//		return errKey
//	}
//
//	h := crypto.SHA256.New()
//	h.Write([]byte(signHex))
//	hash := h.Sum(nil)
//
//	signDecode, errDecode := base64.StdEncoding.DecodeString(sign)
//	if errDecode != nil {
//		return errDecode
//	}
//
//	return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), crypto.SHA256, hash, signDecode)
//}
//
//func MD5Hex(val string) string {
//	md5Obj := md5.New()
//	md5Obj.Write([]byte(val[:]))
//	return hex.EncodeToString(md5Obj.Sum(nil))
//}

//RSA公钥私钥产生
func GetRsaKey(bits int, tenantId string) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//生成文件
	file, err := os.OpenFile("rsa-data.go", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	private := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	privateB := pem.EncodeToMemory(private)
	data := fmt.Sprintf("const Private_key%v = `%v`\n", tenantId, string(privateB))
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	public := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicB := pem.EncodeToMemory(public)
	_, err = file.WriteString(fmt.Sprintf("const Public_key%v = `%v`\n", tenantId, string(publicB)))
	if err != nil {
		return err
	}
	return nil
}

const Private_keyd = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDop2sdYJQ4yCEg
Qu5Lp32Ou82zBTs6ty62qgMy4oF21hyZQq7iKmuZZY4gBsz5GShgncTqpLD0TxiZ
qOYpARyLGibsQPZ3wIJnPPlbtCuFEDZdEhfLo08aGv01TgnhIcoE7OwxR3dV9jeJ
N77qw8ZJzSR2UlAZ0y3QvNxPWqJBoWap5d8UAU2jQrx2d+uj7EVUTVtRsCE4QPeU
z385e0c1VOu8qDYz/j/rDrS2BlsmeDXq9ZFobmvoO+pZ2lHjFTRCeaLCVKnpX8d4
+mZK3bTpXm9Fdu8Nx9jZrC3vykirqUPnnZ7TB+pBvRt3Vz9SEpk1yXD/8lNyBN62
4IJLSXdDAgMBAAECggEAfJJgLUuwMbMe4ZpU49dbyFhQrMFpVGgPMClaKx3S+mFs
0Lc+0sSp9mnFLurVR6+rygfQD199jGLppiUkj+ITeXvYSXoDPl2qtUKVtf+Dqezj
XvQ4H4Zi7XR0Dd2qNoyUEg0V7tD4WePLGsLpi+SlwJCCLISodRt5FaJ6SFccOA0B
ZKz9JX4hlEOVkjKgFIoovM9G9F+rvGheoDdEpyRRXBVI9G5HhRW5cqB0sNK6A9gX
i/Uv3jeno4W4QO5FboR+S0snyoPOMmQwkhYpJAWblIldEUfmuagnmp0a6sXQd2t5
Wfr7A4P8JtS+u507r6Biquikob+XJApidIPKniUiKQKBgQD7xc05fWcbXv5z77aM
4Te+ZOviQhNGuwz0ngArSBcOVUiA4CgfBQrle6KnZyomUmhGAq6TPjU2iQTtDlZx
UY8dUvZ2ENpuRhNmS8i7Rh+hSrNUVCSK1v1n4+asiSfYYyJws6bCgw7pa0YcszWu
GNvef1gctchLxDCs+rHne5stdwKBgQDsj3BJam1Z/J5D1G1a8aQ6MiPdm4wD3DlS
g+jN2hzCyG+nxaKXKZnx/FE3rV4z62nt6HYuItCd1u+GsLRlTuzlhqVuyo/MzlSj
X0r9Cdf/81h1oEVmAS2puCSNY1TqtE7g/r/0wWDegAI+fPVPyvWJRwa9n7q2KBKc
Abge2YRHlQKBgQD1DLnJqdewGU5aM0efaRmjc4DvMFaosihS8nHBrqHaLoGqBgKm
5naLk0Fl5BBvSif5dGTMJXEPil9EB391Peeop/YARjkDuarqFvrh48enahiPDHKg
u83az0PWTIx+nUaJISI/EeZypBmSl464y7M8pP9yui+gJu0lf7+mSXVo0wKBgF+k
itCUAAxO76oa+++2HSEOXqPdnNl+s4piHMEFu3UhVsttQ5R8VGqbCjdJl/nD53sx
7n4uw0vdt9AsJ3OCWpNeQgquST+T+HJpN8dgsH0iZRSBrS1VsqGY+uZTT+To669a
MEAD42dyN/YNzZzqQSW0mswWBYZaY1PB+jA2352VAoGARurIZOma2sXn1bJHmELH
EbTUJFDPpcR+ZHPwgi+dfzJrGssLSN9fJP3h1AZcKtSS6R9NqkMzYSZBRpnoaJHM
1Sl++/mYa+w+4PsULz7c5RRb4KF3bqouwRAlkcnuAwyt4MaFI+SAZNH8rF2EFxto
ZxKO0EJ8PhrqcHHYjUgHaqw=
-----END PRIVATE KEY-----`

const private_key6 = `-----BEGIN PRIVATE KEY-----
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
Tw9mvS2FAhEJLuAyAv8Z++dZ7j/iCtV3
-----END PRIVATE KEY-----
`
