package cryptography

import (
	"fmt"
	"math/rand"
)

type ECDSA struct {
	curve *EllipticCurve
	CryptosystemInterface
}

var _ECDSAcurve *ECDSA

func ECDSAcurve() *ECDSA {
	if _ECDSAcurve != nil {
		return _ECDSAcurve
	}

	curve := new(EllipticCurve)
	g := Point{
		MIntFromString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
		MIntFromString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"),
		curve}
	curve.a = MIntFromString("0")
	curve.b = MIntFromString("7")
	curve.n = MIntFromString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
	curve.g = g

	_ECDSAcurve = &ECDSA{curve: curve}

	return _ECDSAcurve
}

type ECDSA_account struct {
	publicKey  MInt
	privateKey MInt
}

func (ecdsa *ECDSA) GeneratePrivateKey() MInt {
	nbDigits := len(ecdsa.curve.n.Digits)
	digits := make([]uint8, nbDigits)
	for i := 0; i < nbDigits; i++ {
		digits[i] = uint8(rand.Intn(16))
	}
	return MInt{digits}
}

func (ecdsa *ECDSA) GetPublicKey(privateKey MInt) Point {
	publicKeyPoint := ecdsa.curve.g.Multiply(privateKey)
	return publicKeyPoint
}

func (ecdsa *ECDSA) Sign(message string) string {
	// TODO: implement
	return "ECDSA signature"
}

func (ecdsa *ECDSA) Verify(message string, signature string) bool {
	// TODO: implement
	return true
}

func TestCrypto() {
	g2 := ECDSAcurve().curve.g.Multiply(MIntFromString("2"))
	fmt.Println(g2.ToString())
	pk := ECDSAcurve().GeneratePrivateKey()
	fmt.Println(pk.ToString())
	pbk := ECDSAcurve().GetPublicKey(MIntFromString("70e9df4822bde07ea8e679b27c05cdbe971d9a98d408c158da651a330b8dba2e"))
	fmt.Println(pbk.ToString())
}
