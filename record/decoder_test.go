package record

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var sample = `JHON DOE            19860719000000000172
JHON DOE            000000000000000001a2`

type Person struct {
	Name     string    `record:"20,upper"`
	Birthday time.Time `record:"8"`
	Height   int       `record:"12"`
}

func TestErrorHandling(t *testing.T) {
	r := bytes.NewBuffer([]byte(sample))
	d := NewDecoder(r)
	p := new(Person)
	if err := d.Decode(p); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Logf("%v", p)

	if err := d.Decode(p); err == nil {
		t.Errorf("Expected error list, got nil (p=%#v)", p)
	} else {
		if errList, ok := err.(ErrorList); !ok {
			t.Errorf("Returned error is not ErrorList: %t", err)
		} else {
			for _, err := range errList.Errors {
				t.Logf("DecodingError: %v (%t)", err, err)
			}
		}
	}
}

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
