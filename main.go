package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/tlyon3/modexp"
	"math/big"
	"os"
)

var (
	debug = false
	test  = false
	e     = big.NewInt(65537)
	zero  = big.NewInt(0)
	one   = big.NewInt(1)
)

func debugPrint(s string) {
	if debug {
		fmt.Printf("%s", s)
	}
}

func main() {

	if test {
		r := big.NewInt(1)

		fmt.Printf("r == z: %d\n", r.Cmp(big.NewInt(0)))
		p := big.NewInt(7)
		q := big.NewInt(13)
		etest := big.NewInt(5)
		extendedGcd(phiN(p, q), etest)
		fmt.Printf("etest: %s\n", etest.String())
		fmt.Printf("p: %s\n", p.String())
		fmt.Printf("q: %s\n", q.String())
		return
	}

	p, _ := rand.Prime(rand.Reader, 512)
	q, _ := rand.Prime(rand.Reader, 512)
	d := big.NewInt(0)
	//run checks
	valid := false
	//keep generating p and q till we get valid ones
	debugPrint(fmt.Sprintf("Generating p and q...\n"))
	for !valid {
		check := phiN(p, q)
		gcd := new(big.Int)
		ot := new(big.Int)
		if gcd, ot = extendedGcd(check, e); gcd.Cmp(big.NewInt(1)) != 0 {
			valid = false
			p, _ = rand.Prime(rand.Reader, 512)
			q, _ = rand.Prime(rand.Reader, 512)
		} else {
			valid = true

			d.SetString(ot.String(), 10)
		}
	}
	debugPrint(fmt.Sprintf("Generated p and q\n"))
	//ed = 1 (mod phi(n))
	n := big.NewInt(0)
	n.Mul(p, q)
	if d.Cmp(big.NewInt(0)) < 0 {
		d.Add(d, n)
	}
	fmt.Printf("p: %s\n", p.String())
	fmt.Printf("q: %s\n", q.String())
	fmt.Printf("n: %s\n", n.String())
	fmt.Printf("d: %s\n", d.String())
	//verify that for m < n ((m^e % n)^d)%n == m
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the message: ")
	text, _ := reader.ReadString('\n')
	if text == "q" {
		return
	}
	m := new(big.Int)
	m.SetString(text, 10)
	fmt.Printf("m: %s\n", m.String())
	check := big.NewInt(0)
	check = modexp.ModExp(m, e, n)
	debugPrint(fmt.Sprintf("e: %s\n", e.String()))
	check = modexp.ModExp(check, d, n)
	if check.Cmp(m) != 0 {
		panic("Not equal!!!")
	} else {
		fmt.Println("\nValid!")
	}
	//encrypt
	encryptresult := encrypt(m, e, n)
	fmt.Printf("Encrypt result: %s\n", encryptresult.String())
	fmt.Println("Enter the encrypted message: ")
	text, _ = reader.ReadString('\n')
	dm := new(big.Int)
	dm.SetString(text, 10)
	decryptresult := decrypt(dm, d, n)
	fmt.Printf("Decrypt result: %s\n", decryptresult.String())
}

func extendedGcd(a *big.Int, b *big.Int) (gdc *big.Int, ot *big.Int) {
	s := big.NewInt(0)
	t := big.NewInt(1)
	r := new(big.Int)
	r.SetString(b.String(), 10)
	old_s := big.NewInt(1)
	old_t := big.NewInt(0)
	old_r := new(big.Int)
	old_r.SetString(a.String(), 10)
	qp := big.NewInt(0)
	prov := big.NewInt(0)
	quotient := big.NewInt(0)

	for r.Cmp(big.NewInt(0)) != 0 {
		debugPrint("---------\n")
		debugPrint(fmt.Sprintf("s: %s\nt: %s\nr: %s\nold_r: %s\n", s.String(), t.String(), r.String(), old_r.String()))
		quotient.Div(old_r, r)
		debugPrint(fmt.Sprintf("quotient(%s/%s): %s\n", old_r.String(), r.String(), quotient.String()))

		prov.SetString(r.String(), 10)
		r.Sub(old_r, qp.Mul(quotient, prov))
		old_r.SetString(prov.String(), 10)

		prov.SetString(s.String(), 10)
		s.Sub(old_s, qp.Mul(quotient, prov))
		old_s.SetString(prov.String(), 10)

		prov.SetString(t.String(), 10)
		t.Sub(old_t, qp.Mul(quotient, prov))
		old_t.SetString(prov.String(), 10)
	}
	return old_r, old_t
}

func phiN(p *big.Int, q *big.Int) *big.Int {
	result := big.NewInt(0)
	p1 := new(big.Int)
	p.SetString(p.String(), 10)
	p1.Sub(p, big.NewInt(1))
	q1 := new(big.Int)
	q1.SetString(q.String(), 10)
	q1.Sub(q, big.NewInt(1))
	return result.Mul(p1, q1)
}

func encrypt(m *big.Int, e *big.Int, n *big.Int) *big.Int {
	debugPrint(fmt.Sprintf("--IN ENCRYPT--\nm: %s\ne: %s\nn: %s\n", m.String(), e.String(), n.String()))
	debugPrint(fmt.Sprintf("--END ENCRYPT--\n"))
	return modexp.ModExp(m, e, n)
}

func decrypt(m *big.Int, d *big.Int, n *big.Int) *big.Int {
	return modexp.ModExp(m, d, n)
}
