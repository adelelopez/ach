// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"testing"
)

func TestValidators__checkDigit(t *testing.T) {
	cases := map[string]int{
		// invalid
		"":         -1,
		"123456":   -1,
		"1a8ab":    -1,
		"0730002a": -1,
		"0730A002": -1,
		// valid
		"07300022": 8, // Wells Fargo - Iowa
		"10200007": 6, // Wells Fargo - Colorado
	}

	v := validator{}
	for rtn, check := range cases {
		answer := v.CalculateCheckDigit(rtn)
		if check != answer {
			t.Errorf("input=%s answer=%d expected=%d", rtn, answer, check)
		}
		if err := CheckRoutingNumber(fmt.Sprintf("%s%d", rtn, check)); err != nil && check >= 0 {
			t.Errorf("input=%s answer=%d expected=%d: %v", rtn, answer, check, err)
		}
	}
}

func TestValidators__isCreditCardYear(t *testing.T) {
	cases := map[string]bool{
		// invalid (or out of range)
		"10": false,
		"00": false,
		"51": false,
		"17": false,
		// valid
		"20": true,
		"19": true,
	}
	v := validator{}
	for yy, valid := range cases {
		err := v.isCreditCardYear(yy)
		if valid && err != nil {
			t.Errorf("yy=%s failed: %v", yy, err)
		}
		if !valid && err == nil {
			t.Errorf("yy=%s should have failed", yy)
		}
	}
}

func TestValidators__validateSimpleDate(t *testing.T) {
	cases := map[string]string{
		// invalid
		"":       "",
		"01":     "",
		"001520": "", // no 15th month
		"001240": "", // no 40th Day
		"190001": "", // no 0th month
		"190100": "", // no 0th day
		// valid
		"190101": "190101", // Jan 1st
		"201231": "201231", // Dec 31st
		"220731": "220731", // July 31st
		"350430": "350430", // April 30th
		"500229": "500229", // Feb 29th
	}

	v := validator{}
	for input, expected := range cases {
		answer := v.validateSimpleDate(input)
		if expected != answer {
			t.Errorf("input=%q got=%q expected=%q", input, answer, expected)
		}
	}
}

func TestValidators__validateSimpleTime(t *testing.T) {
	cases := map[string]string{
		// invalid
		"":       "",
		"01":     "",
		"012":    "",
		"123142": "",
		// valid
		"0000": "0000",
		"0100": "0100",
		"2359": "2359",
		"1201": "1201",
		"1238": "1238",
	}
	v := validator{}
	for input, expected := range cases {
		answer := v.validateSimpleTime(input)
		if expected != answer {
			t.Errorf("input=%q got=%q expected=%q", input, answer, expected)
		}
	}
}
