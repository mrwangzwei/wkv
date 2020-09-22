package main

import "wkv/dns"

func main() {
	dns.SendJson("127.0.0.1:9991", "www.baidu.com")
}
