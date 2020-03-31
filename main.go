/**
* @Author: D-S
* @Date: 2020/3/20 10:31 上午
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	//CSGO.WS()

GetToken()

}

func GetToken()  {
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
