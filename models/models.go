package models

var ValidType = map[string]int{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type ProperResponse struct {
	Change int
	Thanks string
}
