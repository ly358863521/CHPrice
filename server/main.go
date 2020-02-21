package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// // Set client options
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// // Connect to MongoDB
	// client, err := mongo.Connect(context.TODO(), clientOptions)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Check the connection
	// err = client.Ping(context.TODO(), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connected to MongoDB!")

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
	r.Run(":12345")
}
