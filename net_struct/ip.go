package net_struct

/**
IP Header
Version:协议版本。v4||v6
HeaderLen:整个ip header的长度（非字节长度）。按32位（4字节）为一个长度单位（层）。因此最小值是5,至少有5层，最多15层
DS:区分服务，基本用不到
TotalLen:总长度。实际上还是由MTU的1500bytes决定,超出时就需要分片传
*/

type IP struct {
	Version   [4]Bits
	HeaderLen [4]Bits
	DS        [8]Bits
	TotalLen  [16]Bits
	Next      *levelTwo
}

/**
Identification:标识.数据报计数器.一个数据报被分片时，每片的该段相同
Flag:标志
*/
type levelTwo struct {
	Identification [16]Bits
	Flag           [3]Bits
}
