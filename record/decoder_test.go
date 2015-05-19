package record

import (
	"fmt"
	"time"
)

func ExampleUnmarshal() {
	record := "JHON DOE            19860719000000000172\n"
	p := struct {
		Name     string    `record:"20,upper"`
		Birthday time.Time `record:"8"`
		Height   int       `record:"12"`
	}{}
	if err := Unmarshal([]byte(record), &p); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("%s, birthday %s, %dcm tall", p.Name, p.Birthday.Format("January, 02"), p.Height)
	//Output: JHON DOE, birthday July, 19, 172cm tall
}
