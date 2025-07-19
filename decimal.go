package pltt

import (
	"slices"
	"strings"
)

type decimal int64

const d1 = 1_000_000_000
const dd1 decimal = 1_000_000_000

func (d decimal) Add(u decimal) decimal {
	return d + u
}
func (d decimal) Sub(u decimal) decimal {
	return d - u
}

func (d decimal) Truncate(m decimal) decimal {
	return d - d%m
}

func (d decimal) String() string {
	if d == 0 {
		return "0"
	}
	neg := d < 0
	if neg {
		d = -d
	}
	chars := make([]byte, 0)
	writeDigit := func() {
		chars = append(chars, '0'+byte(d%10))
		d /= 10
	}
	skipped := 0
	for d%10 == 0 && skipped < 9 {
		d /= 10
		skipped++
	}
	if skipped < 9 {
		for range 9 - skipped {
			writeDigit()
		}
		chars = append(chars, '.')
	}
	writeDigit()
	for d != 0 {
		writeDigit()
	}
	if neg {
		chars = append(chars, '-')
	}
	slices.Reverse(chars)
	return string(chars)
}

func (d decimal) Format(string) string {
	return d.String()
}

func parseDecimal(s string) decimal {
	neg := s[0] == '-'
	firstDigitI := 0
	if neg {
		firstDigitI = 1
	}
	ans := decimal(0)
	i := strings.Index(s, ".")
	if i == -1 {
		i = len(s)
	}
	dec := dd1 / 10
	for j := i + 1; j < len(s); j++ {
		ans += decimal(s[j]-'0') * dec
		dec /= 10
	}
	dec = dd1
	for j := i - 1; j >= firstDigitI && j < len(s); j-- {
		ans += decimal(s[j]-'0') * dec
		dec *= 10
	}
	if neg {
		ans = -ans
	}
	return ans
}
