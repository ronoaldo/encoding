package record

import "testing"

type Record struct {
	Seq   int64
	Name  string `record:"6"`
	Code  int64  `record:"6"`
	Label string `record:"6,upper"`
	
    Internal    string `record:"-"`	
	
	Free string
}

func TestMarshal(t *testing.T) {
	encodeTests := []struct{
		src    interface{}
		result string
	} {
		{
			Record{
				Seq: 1,
				Name: "HELLO",
				Code: 0,
				Label: "Hello World",
				Internal: "Unused",
				Free: "world",
			},
			"1 HELLO000000HELLO world",
		},
		{
			Record{
				Seq: 1,
				Name: "ERROR",
				Code: 12345,
				Label: "Stack Overflow",
				Internal: "Unused",
				Free: "Overflow",
			},
			"1 ERROR012345STACK Overflow",
		},
	}
	for _, encTest := range encodeTests {
		b, err := Marshal(encTest.src)
		if err != nil {
			t.Errorf("Unexpected error for Marshal: %v", err)
		}
		if encTest.result != string(b) {
			t.Errorf("Unexpected value for Marshal: `%s`, expected `%s`", string(b), encTest.result)
		}
	}
}
