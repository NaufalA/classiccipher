package main

import (
	"flag"
	"fmt"
	"os"

	crypto "github.com/NaufalA/classiccipher/classic"
)

type cipherKind int

const (
	vigenere cipherKind = iota
	hill
)

func (d cipherKind) String() string {
	return [...]string{"vigenere", "hill"}[d]
}

type mode int

const (
	encrypt mode = iota
	decrypt
)

func vigCipher(m mode, key string, msg string) {
	h, err := crypto.NewVigenere(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var op func(string) (string, error)
	switch m {
	case encrypt:
		op = h.Encrypt
	case decrypt:
		op = h.Decrypt
	}

	res, err := op(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func hillCipher(m mode, key string, msg string) {
	h, err := crypto.NewHill(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var op func(string) (string, error)
	switch m {
	case encrypt:
		op = h.Encrypt
	case decrypt:
		op = h.Decrypt
	}

	res, err := op(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(res)
}

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

var c stringFlag
var key stringFlag
var fk bool
var pt stringFlag
var ct stringFlag

func init() {
	flag.Var(&c, "cipher", "(REQUIRED) Cipher type. [vigenere, hill]")
	flag.Var(&key, "key", "Cipher key string")
	flag.BoolVar(&fk, "findkey", false, "Find hill key string")
	flag.Var(&pt, "plaintext", "String for encryption")
	flag.Var(&ct, "ciphertext", "String for encryption")
}

func main() {
	flag.Parse()

	if !c.set {
		flag.Usage()
		os.Exit(1)
	}

	if c.value == hill.String() && !key.set && fk && pt.set && ct.set {
		hk, err := crypto.FindKeyString(pt.value, ct.value)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(hk)
		os.Exit(0)
	}

	if !key.set {
		flag.Usage()
		os.Exit(1)
	}

	k := key.value

	var m mode
	var text string
	if pt.set {
		m = encrypt
		text = pt.value
	} else if ct.set {
		m = decrypt
		text = ct.value
	} else {
		flag.Usage()
		os.Exit(1)
	}

	switch c.value {
	case hill.String():
		hillCipher(m, k, text)
	case vigenere.String():
		vigCipher(m, k, text)
	default:
		flag.Usage()
		os.Exit(1)
	}

}
