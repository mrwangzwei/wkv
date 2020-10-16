package dotest

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"reflect"
	"testing"
	"unsafe"
)

func TestFunny(t *testing.T) {
	fmt.Println("当前系统位数", 32<<(^uint(0)>>63))
}

func TestStringToBytes(t *testing.T) {
	str := "aaaaa你好a啊你"
	fmt.Println([]byte(str)) //普通转

	//利用底层结构类型强转
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	fmt.Println(*(*[]byte)(unsafe.Pointer(&bh)))
}

func TestBytesToString(t *testing.T) {
	bytes := []byte{97, 97, 97, 97, 97, 228, 189, 160, 229, 165, 189, 97, 229, 149, 138, 228, 189, 160}
	fmt.Println(string(bytes))                      //普通转
	fmt.Println(*(*string)(unsafe.Pointer(&bytes))) //利用底层结构类型强转
}

func TestChannel(t *testing.T) {
	ChannelUse()
}

func TestPolym(t *testing.T) {
	m := NewMaleDuck()
	m.BirthEgg(1)
	m.LayEgg(2)
	f := NewFemaleDuck()
	f.BirthEgg(3)
	f.LayEgg(4)

	fmt.Println()

	LayEgg(m, 1)
	BirthEgg(m, 2)
	LayEgg(f, 3)
	BirthEgg(f, 4)
}

func TestClass(t *testing.T) {
	b := ClassB{}
	b.Aaa()
	b.ClassA.Aaa()
}

func TestStruct(t *testing.T) {
	type Ex struct {
		Aaaa string
	}

	aa := Ex{
		Aaaa: "aaaaaaa",
	}
	bb := aa

	pa := &aa
	pb := &bb
	fmt.Println(aa, bb)
	fmt.Printf("%p %p", &aa, &bb)
	fmt.Println()
	fmt.Println(&pa, &pb)
}

func TestTypePointer(t *testing.T) {
	a := map[string]string{}
	a["aaa"] = "aaa"
	a["bbb"] = "aaa"
	a["ccc"] = "aaa"
	a["ddd"] = "aaa"
	a["eee"] = "aaa"
	b := a
	fmt.Println(a, b)
	fmt.Printf("%p %p", a, b)
	fmt.Println()

	sa := make([]int, 4)
	sa[0] = 1
	sb := sa
	sc := sa
	fmt.Printf("%p %p %p", sa, sb, sc)
	fmt.Println()

	str := "aaaaa"
	fmt.Printf("%p ", &str)
	fmt.Println()

	arrA := [4]int{}
	arrA[0] = 1
	arrB := arrA
	fmt.Printf("%p %p", &arrA, &arrB)
	fmt.Println()
}

func TestArr(t *testing.T) {
	a := [...]int{3: 1, 1, 1: 2}
	fmt.Println(len(a))
	for index, val := range a {
		fmt.Println(index, val)
	}
}

func TestCrawler(t *testing.T) {
	res, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".s-hotsearch-content .hotsearch-item").Each(func(i int, s *goquery.Selection) {
		content := s.Find(".title-content-title").Text()
		fmt.Printf("%d: %s\n", i, content)
		fmt.Println(s.Find(`.title-content`).Attr(`href`))
	})
}

func TestCrawlerColly(t *testing.T) {
	c := colly.NewCollector()
	c.OnHTML(".title-content", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText(`.title-content-title`), e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("http://www.baidu.com")
}
