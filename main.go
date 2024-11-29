package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://www.florkofcows.com/comic/alabama/"
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	_ = os.MkdirAll("downloads", 0755)

	for {
		img, exists := document.Find("img.size-full").Attr("src")
		if !exists {
			panic("no image")
		}

		parts := strings.Split(url, "/")
		name := parts[len(parts)-2]
		fmt.Println(name)

		go func(img string, name string) {
			res, err := http.Get(strings.ReplaceAll(img, "box5600.temp.domains/~thatsock", "www.florkofcows.com"))
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			f, err := os.Create(filepath.Join("downloads", name+".png"))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = io.Copy(f, res.Body)
			if err != nil {
				panic(err)
			}
		}(img, name)

		url, exists = document.Find(".elementor-post-navigation__next > a").Attr("href")
		if !exists {
			break
		}

		res.Body.Close()

		res, err = http.Get(url)
		if err != nil {
			break
		}

		document, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			panic(err)
		}
	}

	res.Body.Close()

}
