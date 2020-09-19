package main

import "wkv/dns"

func main() {
	dns.Send("127.0.0.1:9901", "www.baidu.com")
}
