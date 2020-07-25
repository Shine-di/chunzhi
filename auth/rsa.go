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
	"strings"
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
				fmt.Println("私钥为空使用商户1")
				keys[1] = "1"
				params = append(params, "tenant_id=1")
			} else {
				params = append(params, e)
			}
		}
		privateKey = Private_keyd
	} else {
		params = strings.Split(param, "&")
	}
	result := make([]string, 0)
	for _, param := range params {
		if strings.HasPrefix(param, "time_stamp") {
			//t := time.Now().Unix()
			//s := strconv.Itoa(int(t))
			//result = append(result, "time_stamp="+s)
			result = append(result, param)
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
