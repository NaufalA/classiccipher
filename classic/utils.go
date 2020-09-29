package classic

import "fmt"

// upperOnly extracts letters A-Z and a-z from a string and
// returns them all upper case in a byte slice.
// returns error if resulting slice is empty.
func upperOnly(str string) ([]byte, error) {
	upper := make([]byte, 0, len(str))
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'A' && c <= 'Z' {
			upper = append(upper, c)
		} else if c >= 'a' && c <= 'z' {
			upper = append(upper, c-32)
		}
	}
	if len(upper) <= 0 {
		return nil, fmt.Errorf("string doesn't contain A-Z nor a-z")
	}
	return upper, nil
}

func reverseMap(m map[byte]byte) map[byte]byte {
	n := make(map[byte]byte)
	for k, v := range m {
		n[v] = k
	}
	return n
}
