package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
)

var directory = "/root/Photo/"

func main() {
	getPages(10, 0, "http://jandan.net/girl")
}

func getPages(total int, temp int, url string) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("爬取第", temp, "轮次失败")
	}
	fmt.Println("爬取第", temp, "轮次")

	doc, _ := htmlquery.Parse(resp.Body)
	resp.Body.Close()
	//获取图片链接
	list := htmlquery.Find(doc, "//ol[@class='commentlist']/li//img/@src")
	for _, n := range list {
		SaveToLocal(htmlquery.InnerText(n), directory)
	}

	if temp == total {
		return
	} else {
		next_page_item := htmlquery.Find(doc, "//div[@class='comments'][2]/div/a[last()]/@href")
		new_url := "http:" + htmlquery.InnerText(next_page_item[0])
		fmt.Println(new_url)
		temp = temp + 1
		fmt.Println("你好")
		getPages(total, temp, new_url)
	}

}

func SaveToLocal(url string, directory string) {
	url = "https:" + url
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("保存图片失败")
		return
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	path_list := strings.Split(url, "/")
	name := path_list[len(path_list)-1]
	full_name := directory + name
	file, err := os.Create(full_name)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Println("图片下载失败")
	}
	fmt.Println("完成下载")
}
