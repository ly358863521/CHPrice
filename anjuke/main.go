package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Price get price
type Price struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

func getCUrl() {
	var CityURL []map[string]string
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	c.OnHTML(".letter_city", func(e *colly.HTMLElement) {
		e.ForEach(".city_list", func(_ int, e *colly.HTMLElement) {
			letterCity := make(map[string]string)
			e.ForEach("a", func(_ int, e *colly.HTMLElement) {
				letterCity[e.Text] = e.Attr("href") + "/market"
			})
			CityURL = append(CityURL, letterCity)
		})
	})
	c.Visit("https://www.anjuke.com/sy-city.html")
	f, _ := os.Create("./CityURL.json")
	json.NewEncoder(f).Encode(CityURL)
	f.Close()
}

func getCPrice() {
	CHPrice := make(map[string][]int)
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{Parallelism: 5})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	c.OnHTML(".bigArea", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			e.Request.Visit(e.Attr("href"))
		})
	})
	c.OnHTML("script", func(e *colly.HTMLElement) {
		index := strings.Index(e.Text, "'regionChart")
		if index > 0 {
			s := e.Text[index:]
			s = s[strings.Index(s, "ydata")+7 : strings.Index(s, ";")-4]
			s = s[:strings.Index(s, "}")+1]
			var p Price
			fmt.Println(s)
			if err := json.Unmarshal([]byte(s), &p); err == nil {
				CHPrice[p.Name] = make([]int, len(p.Data))
				copy(CHPrice[p.Name], p.Data)
			}
		}
	})
	CityURL := readJSON("CityURL.json")
	for _, cityList := range CityURL {
		for _, v := range cityList {
			time.Sleep(time.Second)
			c.Visit(v)
		}
	}
	// c.Visit("https://anshan.anjuke.com/market/")
	c.Wait()
	f, _ := os.Create("CityPrice.json")
	json.NewEncoder(f).Encode(CHPrice)
	f.Close()
}

func readJSON(filename string) []map[string]string {
	var CityURL []map[string]string
	jsonfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonfile)
	json.Unmarshal(byteValue, &CityURL)
	return CityURL

}
func main() {
	// M := chp.New(2019, 3)
	// fmt.Println(M)
	// getCUrl()
	getCPrice()
	// a := readJSON("CityURL.json")
	// fmt.Println(a)

}
