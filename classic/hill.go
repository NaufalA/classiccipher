package classic

import (
	"fmt"
	"math"
	"strings"
)

const mod = 26

// Hill is an instance of the Hill Cipher on an specific alphabet
type Hill struct {
	mKey *Matrix
}

// NewHill initializes new cipher ready for given alphabet
func NewHill(key string) (*Hill, error) {
	bs, err := upperOnly(key)
	if err != nil {
		return nil, fmt.Errorf("effective key must contain at least 2 symbols, got \"%s\"", string(bs))
	}
	mk, err := newKey(string(bs))
	if err != nil {
		return nil, err
	}
	return &Hill{mk}, nil
}

// verifyText makes sure message are usable in the current cipher.
// Returns effective message if valid.
func (h *Hill) verifyText(str string) (string, error) {
	bs, err := upperOnly(str)
	if err != nil {
		return "", fmt.Errorf("effective message is empty")
	}
	if len(bs)%h.mKey.order != 0 {
		return "", fmt.Errorf("message length is not multiple of key's length, consider adding padding")
	}
	return string(bs), nil
}

// performOperations apply cipher matrix operations on the given key and text. It assume all
// key and message validations were applied before. Returns the resulting string.
func (h *Hill) performOperations(key *Matrix, str string) string {
	// Use builder for optimum string creation
	var result strings.Builder

	for i := 0; i < len(str); i += key.order {
		vector := make([]int, key.order)
		for j, r := range str[i : i+key.order] {
			vector[j] = int(r) - 'A'
		}
		prodVector, _ := key.VectorProductMod(mod, vector...) // Neglect error because size is exact
		for _, ri := range prodVector {
			r := rune(ri) + 'A'
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Encrypt plain text using given key. Returns an error if message length
// is not multiple of key's order (matrix order).
func (h *Hill) Encrypt(raw string) (string, error) {
	pt, err := h.verifyText(raw)
	if err != nil {
		return "", err
	}
	return h.performOperations(h.mKey, pt), nil
}

// Decrypt cipher text using given key. Returns an error if cipher text length
// is not multiple of key's order (matrix order).
func (h *Hill) Decrypt(raw string) (string, error) {
	ct, err := h.verifyText(raw)
	if err != nil {
		return "", err
	}
	invertedKey, _ := h.mKey.InverseMod(mod) // Neglect error since it's checked by key-text verification
	return h.performOperations(invertedKey, ct), nil
}

// MKeyString returns matrix key in formatted string
func (h *Hill) MKeyString() string {
	return Matrix(*h.mKey).String()
}

// StoM initializes a Hill Cipher square matrix from string
func StoM(raw string) (*Matrix, error) {
	bs, err := upperOnly(raw)
	if err != nil {
		return nil, fmt.Errorf("effective string is empty")
	}
	sqr := math.Sqrt(float64(len(bs)))
	if sqr-math.Floor(sqr) != 0 {
		return nil, fmt.Errorf("effective string size must be a square number, got %d", len(bs))
	}
	if int(sqr) < 2 {
		return nil, fmt.Errorf("cannot create string of order %d < 2", int(sqr))
	}
	is := make([]int, len(bs))
	for i, b := range bs {
		is[i] = int(b) - 'A'
	}
	m, _ := NewMatrix(int(sqr), is) // Error is neglected since order is square
	if !m.IsInvertibleMod(mod) {
		return nil, fmt.Errorf("string is not invertible modulo %d", mod)
	}
	return m, nil
}

// newKey initializes a Hill Cipher matrix key from string
func newKey(rawK string) (*Matrix, error) {
	bs, err := upperOnly(rawK)
	if err != nil {
		return nil, fmt.Errorf("effective key is empty")
	}
	sqr := math.Sqrt(float64(len(bs)))
	if sqr-math.Floor(sqr) != 0 {
		return nil, fmt.Errorf("effective key size must be a square number, got %d", len(bs))
	}
	if int(sqr) < 2 {
		return nil, fmt.Errorf("cannot create key of order %d < 2", int(sqr))
	}
	is := make([]int, len(bs))
	for i, b := range bs {
		is[i] = int(b) - 'A'
	}
	key, _ := NewMatrix(int(sqr), is) // Error is neglected since order is square
	if !key.IsInvertibleMod(mod) {
		return nil, fmt.Errorf("key is not invertible modulo %d", mod)
	}
	return key, nil
}

// FindKey returns
func FindKey(rawP, rawC string) (*Matrix, error) {
	pm, err := StoM(rawP)
	if err != nil {
		return nil, err
	}
	cm, err := StoM(rawC)
	if err != nil {
		return nil, err
	}
	pi, err := pm.InverseMod(mod)
	k, err := pi.CrossProductMod(mod, cm)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// FindKeyString returns
func FindKeyString(rawP, rawC string) (string, error) {
	km, err := FindKey(rawP, rawC)
	if err != nil {
		return "", err
	}
	var r []rune
	for i := 0; i < km.order; i++ {
		for j := 0; j < km.order; j++ {
			r = append(r, rune(km.data[i][j]+'A'))
		}
	}
	return string(r), nil
}
