package hprice

// Hprice house price
type Hprice struct {
	Name  string
	URL   string
	Price []int
	Area  []*Hprice
}
// Chp China house price
type Chp struct {
	City  []*Hprice
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
	return &Chp{City: []*Hprice{}, Year: year, Month: month}
}
