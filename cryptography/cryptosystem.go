package cryptography

type CryptosystemInterface interface {
	Sign(message string) string
	Verify(message string, signature string) bool
}
