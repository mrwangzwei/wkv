package dns

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
)

//从百度查到的ip地理信息
type IpAddrFromBaidu struct {
	Address string `json:"address"`
	Content struct {
		Address       string `json:"address"`
		AddressDetail struct {
			City         string `json:"city"`
			CityCode     int    `json:"city_code"`
			District     string `json:"district"`
			Province     string `json:"province"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"address_detail"`
		Point struct {
			X string `json:"x"`
			Y string `json:"y"`
		}
	} `json:"content"`
}

const baiduAk = "ZayMXSEG51dNHHVLGjPmFFOditNvVdGb"

func RequestIpAddr(ip string) (ipAddrFromBaiduRes IpAddrFromBaidu, err error) {
	url := "https://api.map.baidu.com/location/ip"
	urlValues := url2.Values{}
	urlValues.Add("ip", ip)
	urlValues.Add("ak", baiduAk)
	urlValues.Add("coor", "bd09ll")
	reqBody := urlValues.Encode()
	resp, err := http.Post(url,
		"application/json",
		strings.NewReader(reqBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &ipAddrFromBaiduRes)
	return
}
