package cryptography

import "fmt"


type ECDSA struct {
	curve EllipticCurve
	CryptosystemInterface
}

var ECDSAcurve = EllipticCurve{
	a: mIntFromString("0"),
	b: mIntFromString("7"),
	n: mIntFromString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
	g: Point{
		mIntFromString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"), 
		mIntFromString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"), 
		*ECDSAcurve}}

type ECDSA_account struct {
	publicKey int
	privateKey int
}

func (ecdsa ECDSA) Sign(message string) string {
	// TODO: implement
	return "ECDSA signature"
}

func (ecdsa ECDSA) Verify(message string, signature string) bool {
	// TODO: implement
	return true
}

func TestCrypto() {
	g2 := ECDSAcurve.g.Multiply(2)
	fmt.Println(g2.x, g2.y)
}