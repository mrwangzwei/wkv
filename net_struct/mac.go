package net_struct

/**
MAC Header
数据报文最外层 最小64bytes 其中规定数据部分最大1500bytes 即整包最大1564bytes
各机构制定的帧格式有异，此处进以ethernet II为例
DAddressMac:目的mac		6bytes
SAddressMac:源mac	6bytes
ProtocolType:协议类型	2bytes
Data:数据	MTU 46-1500bytes
FCS:校验和。CRC算法校验	4bytes
*/
type Mac struct {
	DAddressMac  Bytes
	SAddressMac  Bytes
	ProtocolType Bytes
	Data         Bytes
	FCS          Bytes
}
