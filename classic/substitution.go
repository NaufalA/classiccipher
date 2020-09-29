package classic

import "fmt"

// Substitution is a mono-alphabetic substitution cipher.
type Substitution struct {
	key map[byte]byte
}

// NewSubstitution creates lorem
func NewSubstitution(key string) (*Substitution, error) {
	bs, err := upperOnly(key)
	if err != nil {
		return nil, fmt.Errorf("effective key is empty")
	}
	if len(bs) != 26 {
		return nil, fmt.Errorf("key \"%s\" is not 26 characters long", string(bs))
	}

	k := make(map[byte]byte, 26)
	for i, b := range bs {
		k[b] = 'A' + byte(i)
	}
	return &Substitution{k}, nil
}

// Encrypt encrypts plaintext pt by substitution with the key.
func (s *Substitution) Encrypt(pt string) string {
	bs, _ := upperOnly(pt)
	for i, b := range bs {
		bs[i] = s.key[b]
	}
	return string(bs)
}

// Decrypt decrypts ciphertext ct by substitution with the key.
func (s *Substitution) Decrypt(ct string) string {
	bs, _ := upperOnly(ct)
	r := reverseMap(s.key)
	for i, b := range bs {
		bs[i] = r[b]
	}
	return string(bs)
}
