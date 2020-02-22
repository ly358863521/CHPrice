package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type province struct {
	name  string
	citys []string
	price int
}

func provinceinit(c *mongo.Collection) map[string]int {
	res := make(map[string]int)
	cur, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		name := result["name"].(string)
		if name == "不存在" || name == "无归属" {
			continue
		}
		name = strings.TrimSuffix(name, "省")
		name = strings.TrimSuffix(name, "市")
		name = strings.TrimSuffix(name, "自治区")
		name = strings.TrimSuffix(name, "回族")
		name = strings.TrimSuffix(name, "壮族")
		name = strings.TrimSuffix(name, "维吾尔")
		res[name] = int(result["price"].(int32))
	}
	// fmt.Println(res)
	return res
}

func getProvince(p, c, a mongo.Collection, pName string) gin.H {
	var province bson.M
	filter := bson.D{primitive.E{Key: "name", Value: primitive.Regex{Pattern: pName, Options: ""}}}
	err := p.FindOne(context.Background(), filter).Decode(&province)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(province)
	citys := province["citys"].(primitive.A)
	if len(citys) == 0 {
		return nil
	}

	firstCity := citys[0]
	var cityslist []string
	var cityPrice []int
	var thePrice []int
	var lastPrice []int
	var arealist []string
	var areathePrice []int
	var arealastPrice []int

	for _, i := range citys {
		var price bson.M
		err := c.FindOne(context.Background(), bson.D{{Key: "Name", Value: i}}).Decode(&price)
		// fmt.Println(price)
		if err != nil {
			log.Fatal(err)
		}
		if price == nil {
			continue
		}
		pricelist := price["Price"].(primitive.A)
		if len(pricelist) > 2 && int(pricelist[len(pricelist)-1].(int32)) != 0 {
			thePrice = append(thePrice, int(pricelist[len(pricelist)-1].(int32)))
			lastPrice = append(lastPrice, int(pricelist[len(pricelist)-2].(int32)))
			cityslist = append(cityslist, i.(string))
		} else {
			continue
		}

		if firstCity == i {
			// fmt.Println(price)
			if price["Area"] != nil {
				arealisttmp := price["Area"].(primitive.A)
				arealist = make([]string, 0)
				for _, v := range arealisttmp {
					var areabson bson.M
					err := a.FindOne(context.Background(), bson.D{{Key: "Name", Value: v.(string)}, {Key: "City", Value: firstCity}}).Decode(&areabson)
					if err != nil {
						log.Fatal(err)
					}
					areaPrice := areabson["Price"].(primitive.A)
					if len(areaPrice) > 2 && int(areaPrice[len(areaPrice)-1].(int32)) != 0 {
						areathePrice = append(areathePrice, int(areaPrice[len(areaPrice)-1].(int32)))
						arealastPrice = append(arealastPrice, int(areaPrice[len(areaPrice)-2].(int32)))
						arealist = append(arealist, v.(string))
					}
				}
			}
			cityPrice = make([]int, len(pricelist))
			for k, v := range pricelist {
				cityPrice[k] = int(v.(int32))
			}
			// copy(cityPrice, pricelist)
		}
	}
	return gin.H{
		"cityName":      firstCity,
		"cityPrice":     cityPrice,
		"citys":         citys,
		"thePrice":      thePrice,
		"lastPrice":     lastPrice,
		"arealist":      arealist,
		"areathePrice":  areathePrice,
		"arealastPrice": arealastPrice,
	}
}
func getAreas(c, a mongo.Collection, cName string) gin.H {
	var price bson.M
	err := c.FindOne(context.Background(), bson.D{{Key: "Name", Value: cName}}).Decode(&price)
	if err != nil {
		log.Fatal(err)
	}
	pricelisttmp := price["Price"].(primitive.A)
	if len(pricelisttmp) == 0 {
		return nil
	}
	pricelist := make([]int, len(pricelisttmp))
	for k, v := range pricelisttmp {
		pricelist[k] = int(v.(int32))
	}
	var arealastPrice []int
	var areathePrice []int
	arealist := make([]string, 0)
	if price["Area"] != nil {
		arealisttmp := price["Area"].(primitive.A)
		for _, v := range arealisttmp {
			var areabson bson.M
			err := a.FindOne(context.Background(), bson.D{{Key: "Name", Value: v.(string)}, {Key: "City", Value: cName}}).Decode(&areabson)
			if err != nil {
				log.Fatal(err)
			}
			areaPrice := areabson["Price"].(primitive.A)
			if len(areaPrice) > 2 && int(areaPrice[len(areaPrice)-1].(int32)) != 0 {
				areathePrice = append(areathePrice, int(areaPrice[len(areaPrice)-1].(int32)))
				arealastPrice = append(arealastPrice, int(areaPrice[len(areaPrice)-2].(int32)))
				arealist = append(arealist, v.(string))
			}

		}
	}
	return gin.H{
		"pricelist":     pricelist,
		"arealist":      arealist,
		"areathePrice":  areathePrice,
		"arealastPrice": arealastPrice,
	}
}

func getArea(a mongo.Collection, pName, aName string) gin.H {
	var areabson bson.M
	err := a.FindOne(context.Background(), bson.D{{Key: "Name", Value: aName}, {Key: "City", Value: pName}}).Decode(&areabson)
	if err != nil {
		log.Fatal(err)
	}
	areaPricetmp := areabson["Price"].(primitive.A)
	areaPrice := make([]int, len(areaPricetmp))
	for i, v := range areaPricetmp {
		areaPrice[i] = int(v.(int32))
	}
	return gin.H{
		"areaprice": areaPrice,
	}
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	provinceSET := client.Database("CHPrice").Collection("province")
	provinceinit(provinceSET)
	citySET := client.Database("CHPrice").Collection("Citys")
	areaSET := client.Database("CHPrice").Collection("Areas")

	r := gin.Default()

	r.StaticFS("/static", http.Dir("../echarts/static"))

	// r.LoadHTMLGlob("../echarts/*")
	r.LoadHTMLFiles("../echarts/hello.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "hello.html", nil)
	})

	r.GET("/getdata", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "liyan",
		})
	})

	r.GET("/getProvince", func(c *gin.Context) {
		provinceMap := provinceinit(provinceSET)
		c.JSON(200, provinceMap)
	})

	r.GET("/getCitys/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println(name)
		c.JSON(200, getProvince(*provinceSET, *citySET, *areaSET, name))
	})

	r.GET("/getAreas/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println(name)
		c.JSON(200, getAreas(*citySET, *areaSET, name))
	})

	r.GET("/getArea/:pname/:aname", func(c *gin.Context) {
		pname := c.Param("pname")
		aname := c.Param("aname")
		c.JSON(200, getArea(*areaSET, pname, aname))
	})

	r.Run(":12345")
}
