package hprice

import (
	"encoding/json"
	"os"
)

// Hprice house price
type Hprice struct {
	Name  string
	Price []int
	Area  []*Area
}

// Area Area house price
type Area struct {
	Name  string
	Price []int
}

// Chp China house price
type Chp struct {
	City  map[string]*Hprice
	Year  [12]int
	Month [12]int
}

// New new chp
func New(startY, startM int) *Chp {
	var year [12]int
	var month [12]int
	for i := 0; i < 12; i++ {
		month[i] = (startM + i) % 12
		if month[i] == 0 {
			month[i] = 12
		}
		year[i] = (startM+i)/13 + startY
	}
	return &Chp{City: map[string]*Hprice{}, Year: year, Month: month}
}

// NewArea new area
func NewArea(name string, price []int) *Area {
	tmpprice := make([]int, len(price))
	copy(tmpprice, price)
	return &Area{Name: name, Price: tmpprice}
}

// Add a city
func (c *Chp) Add(name string) {
	if _, ok := c.City[name]; !ok {
		c.City[name] = &Hprice{Name: name, Price: []int{}, Area: []*Area{}}
	}
}

// Save to json
func (c *Chp) Save(filename string) {
	f, _ := os.Create("../data/" + filename)
	json.NewEncoder(f).Encode(c)
	f.Close()
}

// Add a Area
func (c *Hprice) Add(name string, price []int) {
	var area *Area
	area.Name = name
	area.Price = make([]int, len(price))
	copy(area.Price, price)
	c.Area = append(c.Area, area)
}
