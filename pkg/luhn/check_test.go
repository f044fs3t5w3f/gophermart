package luhn

import "testing"

type testCase struct {
	number string
	valid  bool
}

func TestCheck(t *testing.T) {
	cases := []testCase{
		{
			number: "79927398713",
			valid:  true,
		},
		{
			number: "49927398716",
			valid:  true,
		},
		{
			number: "344323244238",
			valid:  true,
		},
		{
			number: "79927398716",
			valid:  false,
		},
		{
			number: "49927398710",
			valid:  false,
		},
		{
			number: "344323244231",
			valid:  false,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.number, func(t *testing.T) {
			isValid := Check(tCase.number)
			if isValid != tCase.valid {
				t.Errorf("Expected %s to be %t", tCase.number, tCase.valid)
			}
		})
	}
}
