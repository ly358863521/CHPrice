package main

import (
	hp "CHprice/hprice"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// Price get price
type Price struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

var proxyIP = []string{
	"https://124.205.11.245:8888",
	"http://117.85.166.97:8118",
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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

func getProxy() string {
	rand.Seed(time.Now().UnixNano())
	return proxyIP[rand.Intn(len(proxyIP))]
}
func getCPrice() {
	CHPrice := hp.New(2019, 3)
	var URLtoCity sync.Map
	var AreatoCity sync.Map
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		colly.Async(true),
	)

	c.SetRequestTimeout(30 * time.Second)
	c.Limit(&colly.LimitRule{Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", randomString()) //通过随机改变 user-agent
		fmt.Println("Visiting", r.URL)
		// if len(r.URL.Path) > 30 {
		// 	err := c.SetProxy(getProxy())
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML(".bigArea", func(e *colly.HTMLElement) {
		url := e.Request.URL.Hostname()
		k := url[:strings.Index(url, ".")]
		if value, ok := URLtoCity.Load(k); ok {
			cityName := value
			e.ForEach("a", func(_ int, e *colly.HTMLElement) {
				// AreatoCity[e.Text] = cityName
				AreatoCity.Store(e.Text, cityName)
				e.Request.Visit(e.Attr("href"))
			})
		}
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
				if _, ok := CHPrice.City[p.Name]; ok {
					CHPrice.City[p.Name].Price = make([]int, len(p.Data))
					copy(CHPrice.City[p.Name].Price, p.Data)
				} else {
					cityName, _ := AreatoCity.Load(p.Name)
					if _, ok := CHPrice.City[cityName.(string)]; !ok {
						fmt.Println(cityName, p.Name)
					} else {
						area := hp.NewArea(p.Name, p.Data)
						CHPrice.City[cityName.(string)].Area = append(CHPrice.City[cityName.(string)].Area, area)
						// CHPrice.City[cityName].Add(p.Name, p.Data)
					}
				}
			}
		}
	})
	CityURL := readJSON("CityURL.json")
	// for _, cityList := range CityURL {
	// 	for _, v := range cityList {
	// 		time.Sleep(time.Second)
	// 		c.Visit(v)
	// 	}
	// }
	for k, v := range CityURL[1] {
		vv := v[8:strings.Index(v, ".")]
		// fmt.Println("url：", k, v)
		// URLtoCity[vv] = k
		URLtoCity.Store(vv, k)
		CHPrice.Add(k)
		c.Visit(v)
		c.Wait()
	}
	// c.Visit("https://anshan.anjuke.com/market/")
	CHPrice.Save("CityPriceB.json")
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
