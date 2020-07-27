/**
 * @author: D-S
 * @date: 2020/7/25 6:50 下午
 */

package sHttp

import (
	"fmt"
	"game-test/library/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GET struct {
	URL    string
	Header map[string]string
	Proxy  string
}

func (r *GET) Do() {
	client := new(http.Client)
	if r.Proxy != "" {
		u, _ := url.Parse(r.Proxy)
		proxy := http.ProxyURL(u)
		client.Transport = &http.Transport{
			Proxy: proxy,
		}
		log.Info("==使用代理==")
		fmt.Println(r.Proxy)
	}
	req, _ := http.NewRequest("GET", r.URL, nil)
	for k, v := range r.Header {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Info("==数据==")
	fmt.Println(string(body))
}
