package record

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
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
				t.Logf("DecodingError: %v", err)
			}
		}
	}
}

func TestTimeLayout(t *testing.T) {
	r := bytes.NewBuffer([]byte(`19/07/1986`))
	d := NewDecoder(r).TimeLayout("02/01/2006")
	ts := struct {
		Time time.Time `record:"10"`
	}{}
	if err := d.Decode(&ts); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Logf("%v", ts)
}

func TestOptionalTag(t *testing.T) {
	var err error
	r := bytes.NewBuffer([]byte(`     `))
	s := struct {
		N int `record:"5"`
	}{}
	if err = NewDecoder(r).Decode(&s); err == nil {
		t.Errorf("Nil error dec. empty number; expected invalid syntax")
	} else if !strings.Contains(err.Error(), "invalid syntax") {
		t.Errorf("Unexpected error dec. empty number: %v", err)
	}
	t.Logf("Decode error when empty, non-optional: %v", err)

	s1 := struct {
		N int `record:"5,optional"`
	}{}
	if err = NewDecoder(r).Decode(&s1); err != nil {
		t.Errorf("Unexpected error decoding optional number: %v, expected nil", err)
	}
	t.Logf("Decode error when empty, optional: %v", err)
}

func Example_fileParsing() {
	data := `0HEADER____
1DATA     1
9TRAILER___`

	sc := bufio.NewScanner(bytes.NewBuffer([]byte(data)))
	for sc.Scan() {
		line := sc.Text()

		switch []rune(line)[0] {
		case '0':
			header := struct {
				Type int    `record:"1"`
				Text string `record:"10"`
			}{}
			if err := NewDecoder(strings.NewReader(line)).Decode(&header); err != nil {
				fmt.Printf("Unexpected decoding error: %v", err)
			}
			fmt.Printf("Decoded header: %v\n", header)

		case '1':
			content := struct {
				Type int    `record:"1"`
				Text string `record:"9"`
				ID   int64  `record:"1"`
			}{}
			if err := NewDecoder(strings.NewReader(line)).Decode(&content); err != nil {
				fmt.Printf("Unexpected decoding error: %v", err)
			}
			fmt.Printf("Decoded content: %v\n", content)

		default:
			fmt.Printf("Skipping line %s", line)
		}
	}
	//Output: Decoded header: {0 HEADER____}
	//Decoded content: {1 DATA 1}
	//Skipping line 9TRAILER___
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
