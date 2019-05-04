package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"os"

	"flag"
	"strconv"
	"time"

	//go语言版本的request
	"github.com/levigross/grequests"
)

var wg sync.WaitGroup

type stDownInfo struct {
	URL  string
	Name string
}

func down(pkt stDownInfo) {
	if len(pkt.URL) <= 0 {
		return
	}
	s := pkt.URL
	fmt.Printf("准备下载: %v\n", s)
	ss2 := strings.Split(string(s), "/")[2:]
	ss2[0] = "downs"
	ss2[len(ss2)-1] = pkt.Name
	dirname := strings.Join(ss2[:len(ss2)-1], "/")
	if _, err := os.Stat(dirname); err != nil {
		fmt.Printf("创建下载文件夹:%s\n", dirname)
		os.MkdirAll(dirname, os.ModePerm)
	}

	filename := strings.Join(ss2, "/")

	_, err := os.Stat(filename)
	if err == nil {
		fmt.Print("文件已存在,跳过")
		return
	}

	res, _ := grequests.Get(s, &grequests.RequestOptions{
		//结构体可以对指定的类型给值，而不一定都赋值
		Headers: map[string]string{
			"Referer":    "http://www.sina.com",
			"User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36; gamegrd "}})
	//条件需要修改，如果没有图片，返回的是盗链，图片4kb
	if res.StatusCode != 200 {
		fmt.Printf("下载失败，:%s\n", s)
		return
	}

	//图片可能是该网站，返回的盗链图片（4kb左右)
	length := res.Header.Get("Content-Length")
	slen, _ := strconv.Atoi(length)
	if slen < 4100 {
		fmt.Printf("下载内容失败:%s\n", s)
		return
	}

	res.DownloadToFile(filename)
	fmt.Printf("成功下载图片到:%s\n", filename)
}
func main() {
	thread := flag.Int("thread", 1, "线程数量")
	size := flag.String("size", "thumbnail", "大小 large, mw1024, mw690, bmiddle, small, thumb180, thumbnail, square ")
	flag.Parse()
	ch := make(chan stDownInfo, *thread)
	now := time.Now()

	for i := 0; i < *thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range ch {
				down(d)
			}
		}()
	}

	f, err := os.Open("res/res.txt")
	if err != nil {
		fmt.Println("res/res.txt 不存在")
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	i := 130
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, "\r", "", -1)
		url := fmt.Sprintf("https://ws1.sinaimg.in/%s/%s.jpg", *size, line)
		dI := stDownInfo{URL: url, Name: fmt.Sprintf("%d.jpg", i)}
		fmt.Printf("down : %s >>>>>>  %s \n", url, dI.Name)
		ch <- dI
		i++
	}

	close(ch)
	wg.Wait()
	fmt.Printf("下载任务完成，耗时:%#v\n", time.Now().Sub(now).Seconds())
}
