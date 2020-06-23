package model

import (
	"fmt"
	"testing"
)

func TestValidation(t *testing.T) {
	cases := []struct {
		caseName     string
		message      Message
		errorsLength int
	}{
		{"Check Name is exist", Message{"", "test first", Metadata{"2020-12-12 19:02:33.123", "00eec78c-b4a3-11ea-b3de-0242ac130004"}}, 1},
		{"Check Payload is exist", Message{"first", "", Metadata{"2020-12-12 19:02:33.123", "00eec78c-b4a3-11ea-b3de-0242ac130004"}}, 1},
		{"Check Metadata is exist", Message{"first", "test first", Metadata{}}, 2},
		{"Check timestamp is exist", Message{"first", "test first", Metadata{"", "00eec78c-b4a3-11ea-b3de-0242ac130004"}}, 1},
		{"Check timestamp has correct format", Message{"first", "test first", Metadata{"2020-12-12 19:02", "00eec78c-b4a3-11ea-b3de-0242ac130004"}}, 1},
		{"Check UUID is exist", Message{"first", "test first", Metadata{"2020-12-12 19:02:33.123", ""}}, 1},
		{"Check UUID has correct format", Message{"first", "test first", Metadata{"2020-12-12 19:02:33.123", "00eec78c-b4a3-11ea-b3de-0242ac13002222"}}, 1},
	}

	v := NewValidation()
	for _, tCase := range cases {
		err := v.Validate(tCase.message)
		fmt.Println(err, " ", len(err))
		if len(err) != tCase.errorsLength {
			t.Fatalf("%v is failed", tCase.caseName)
		}
	}
}
