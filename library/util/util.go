package util

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/micro/go-micro/metadata"
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

func CopyContext(ctx context.Context) context.Context {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// copy the metadata to prevent race
	md = metadata.Copy(md)

	ctxx := metadata.NewContext(context.Background(), md)
	return ctxx
}

func HttpGet(queryUrl string) ([]byte, error) {
	r, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("Error returning status code :" + fmt.Sprintf("%d", r.StatusCode))
	}
	result, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func HttpHeader(queryUrl string, headers map[string]interface{}) ([]byte, error) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	//增加header选项
	for k, v := range headers {
		reqest.Header.Add(k, fmt.Sprintf("%v", v))
	}
	//处理返回结果
	r, _ := client.Do(reqest)
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("Error returning status code :" + fmt.Sprintf("%d", r.StatusCode))
	}
	result, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

//根据时间计算游戏场次
func WorkTimeIndex(n time.Time, ns []time.Time) int {
	var is []int
	for _, v := range ns {
		is = append(is, int(v.Unix()))
	}
	i := int(n.Unix())
	sort.Ints(is)
	return sort.SearchInts(is, i) + 1
}

//获取奖金
func GetBonus(str string) int {
	bonusStr := ""
	for _, v := range str {
		_, err := strconv.Atoi(string(v))
		if err == nil {
			bonusStr += string(v)
		}
	}
	bonusInt, err := strconv.Atoi(bonusStr)
	if err != nil {
		return 0
	} else {
		return bonusInt
	}
}

func InSliceInt(list []int64, target int64) bool {
	for _, v := range list {
		if target == v {
			return true
		}
	}
	return false
}
func InSliceStr(list []string, target string) bool {
	for _, v := range list {
		if target == v {
			return true
		}
	}
	return false
}

func ReturnFileName(fileName string) string {
	return strings.ToLower(strings.Replace(fileName, " ", "", -1))
}

func ParseShanghaiLocation(timeStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
}
