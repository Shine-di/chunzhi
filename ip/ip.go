/**
 * @author: D-S
 * @date: 2020/4/20 4:36 下午
 */

package ip

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetExternalIp() string {
	client := http.Client{}

	resp, err := client.Get("sHttp://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return string(rb)
}

func GetInternalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Error(err)
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func IsBelong(ip, cidr string) bool {
	ipAddr := strings.Split(ip, `.`)
	if len(ipAddr) < 4 {
		return false
	}
	cidrArr := strings.Split(cidr, `/`)
	if len(cidrArr) < 2 {
		return false
	}
	var tmp = make([]string, 0)
	for key, value := range strings.Split(`255.255.255.0`, `.`) {
		iint, _ := strconv.Atoi(value)

		iint2, _ := strconv.Atoi(ipAddr[key])

		tmp = append(tmp, strconv.Itoa(iint&iint2))
	}
	return strings.Join(tmp, `.`) == cidrArr[0]
}

// 代理服务器(产品官网 www.16yun.cn)
const ProxyServer = "t.16yun.cn:31111"

type ProxyAuth struct {
	Username string
	Password string
}

func (p ProxyAuth) ProxyClient() http.Client {

	var proxyURL *url.URL
	if p.Username != "" && p.Password != "" {
		proxyURL, _ = url.Parse("sHttp://" + p.Username + ":" + p.Password + "@" + ProxyServer)
	} else {
		proxyURL, _ = url.Parse("sHttp://" + ProxyServer)
	}
	return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
}

func main() {

	targetURI := "https://httpbin.org/ip"

	// 初始化 proxy sHttp client
	client := ProxyAuth{"username", "password"}.ProxyClient()

	request, _ := http.NewRequest("GET", targetURI, bytes.NewBuffer([]byte(``)))

	// 设置Proxy-Tunnel
	// rand.Seed(time.Now().UnixNano())
	// tunnel := rand.Intn(10000)
	// request.Header.Set("Proxy-Tunnel", strconv.Itoa(tunnel) )

	response, err := client.Do(request)

	if err != nil {
		panic("failed to connect: " + err.Error())
	} else {
		bodyByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("读取 Body 时出错", err)
			return
		}
		response.Body.Close()

		body := string(bodyByte)

		fmt.Println("Response Status:", response.Status)
		fmt.Println("Response Header:", response.Header)
		fmt.Println("Response Body:\n", body)
	}
}
