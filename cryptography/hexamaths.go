package cryptography

import (
	"math/rand"
	"strconv"
)

type MInt struct {
	Digits []uint8
}

func (m MInt) getDigits() []uint8 { // reverse Digits: less significant first
	Digits := make([]uint8, len(m.Digits))
	for i := len(m.Digits) - 1; i >= 0; i-- {
		Digits[len(m.Digits)-1-i] = m.Digits[i]
	}
	return Digits
}

func (m MInt) Add(m2 MInt) MInt {
	var m3 MInt
	DigitsNb := len(m.Digits)
	if len(m2.Digits) < DigitsNb {
		DigitsNb = len(m2.Digits)
	}
	m3.Digits = make([]uint8, DigitsNb)
	mDigits := m.getDigits()
	m2Digits := m2.getDigits()

	for i := 0; i < DigitsNb; i++ {
		carry := (m3.Digits[i] + mDigits[i] + m2Digits[i]) / 16
		m3.Digits[i] = (m3.Digits[i] + mDigits[i] + m2Digits[i]) % 16

		j := i + 1
		for carry > 0 {
			if j == len(m3.Digits) {
				m3.Digits = append(m3.Digits, 0)
			}
			m3.Digits[j] += carry
			carry = m3.Digits[j] / 16
			m3.Digits[j] = m3.Digits[j] % 16
			j++
		}
	}
	if len(mDigits) < len(m2Digits) {
		for i := len(mDigits); i < len(m2Digits); i++ {
			if i >= len(m3.Digits) {
				m3.Digits = append(m3.Digits, m2Digits[i])
			} else {
				m3.Digits[i] += m2Digits[i]
			}
		}
	} else if len(m2Digits) < len(mDigits) {
		for i := len(m2Digits); i < len(mDigits); i++ {
			if i >= len(m3.Digits) {
				m3.Digits = append(m3.Digits, mDigits[i])
			} else {
				m3.Digits[i] += mDigits[i]
			}
		}
	}

	for i, j := 0, len(m3.Digits)-1; i < j; i, j = i+1, j-1 {
		m3.Digits[i], m3.Digits[j] = m3.Digits[j], m3.Digits[i]
	}
	return m3
}

func (m MInt) Sub(m2 MInt) MInt {
	if !m.GreaterEq(m2) {
		panic("Substraction of a smaller number from a bigger one: " + m.ToString() + " - " + m2.ToString() + " is not possible.")
	}

	var m3 MInt
	mDigits := m.getDigits()
	m2Digits := m2.getDigits()
	DigitsNb := len(m2Digits)

	m3.Digits = make([]uint8, len(mDigits))

	for i := 0; i < DigitsNb; i++ {
		if mDigits[i] < m2Digits[i] {
			mDigits[i] += 16
			mDigits[i+1]--
		}
		m3.Digits[i] = mDigits[i] - m2Digits[i]
	}

	for i, j := 0, len(m3.Digits)-1; i < j; i, j = i+1, j-1 {
		m3.Digits[i], m3.Digits[j] = m3.Digits[j], m3.Digits[i]
	}

	for i := 0; i < len(m3.Digits); i++ { // remove leading 0s
		if m3.Digits[i] != 0 {
			m3.Digits = m3.Digits[i:]
			break
		}
	}

	return m3
}

func (m MInt) SubMod(m2 MInt, mod MInt) MInt {
	if !m.GreaterEq(m2) {
		return m.Add(mod).SubMod(m2, mod)
	}
	return m.Sub(m2).Mod(mod)
}

func (m MInt) Multi(i int) MInt {
	return m.Mult(MIntFromString(strconv.Itoa(i)))
}

func (m MInt) Mult(m2 MInt) MInt {
	var m3 MInt
	var multDigitsnb int
	var m2Digits = m2.getDigits()
	var mDigits = m.getDigits()

	if len(m.Digits) < len(m2.Digits) {
		multDigitsnb = len(m.Digits)
	} else {
		multDigitsnb = len(m2.Digits)
		m2Digits, mDigits = mDigits, m2Digits
	}

	for i := 0; i < multDigitsnb; i++ {
		var m4 MInt
		m4.Digits = make([]uint8, i)
		for j := 0; j < len(m.Digits); j++ {
			if j >= len(m4.Digits) {
				m4.Digits = append(m4.Digits, 0)
			}
			a := (m4.Digits[j] + mDigits[i]*m2Digits[j])
			carry := a / 16
			m4.Digits[j] = a % 16
			c := j + 1
			for carry > 0 {
				if c == len(m4.Digits) {
					m4.Digits = append(m4.Digits, 0)
				}
				m4.Digits[c] += carry
				carry = m4.Digits[c] / 16
				m4.Digits[c] = m4.Digits[c] % 16
				c++
			}
		}
		for j := 0; j < i; j++ {
			m4.Digits = append(m4.Digits, 0)
		}
		m4.Digits = m4.getDigits()
		for j := 0; j+i < len(m4.Digits); j++ {
			m4.Digits[j] = m4.Digits[j+i]
		}
		for j := len(m4.Digits) - 1; j >= len(m4.Digits)-i; j-- {
			m4.Digits[j] = 0
		}

		m3.Digits = m3.getDigits()
		m3 = m3.Add(m4)
		m3.Digits = m3.getDigits()
	}

	m3.Digits = m3.getDigits()
	return m3
}

func (m MInt) Div(i int) MInt {
	return m.Divide(MIntFromString(strconv.Itoa(i)))
}

func (m MInt) Divide(m2 MInt) MInt {
	q := 0
	for m.GreaterEq(m2) {
		m = m.Sub(m2)
		q++
	}
	return MIntFromString(strconv.Itoa(q))
}

func (m MInt) Modi(i int) MInt {
	return m.Mod(MIntFromString(strconv.Itoa(i)))
}

func (m MInt) Mod(m2 MInt) MInt {
	for m.GreaterEq(m2) {
		m = m.Sub(m2)
	}
	return m
}

func (m MInt) Eq(m2 MInt) bool {
	if len(m.Digits) != len(m2.Digits) {
		return false
	}
	for i := 0; i < len(m.Digits); i++ {
		if m.Digits[i] != m2.Digits[i] {
			return false
		}
	}
	return true
}

func (m MInt) realLen() int {
	for i := 0; i < len(m.Digits); i++ {
		if m.Digits[i] != 0 {
			return len(m.Digits) - i
		}
	}
	return 0
}

func (m MInt) GreaterEq(m2 MInt) bool {
	if m.Eq(m2) {
		return true
	}

	realLenm := m.realLen()
	realLenm2 := m2.realLen()

	if realLenm > realLenm2 {
		return true
	} else if realLenm < realLenm2 {
		return false
	}

	for i := 0; i < realLenm; i++ {
		if m.Digits[i] > m2.Digits[i] {
			return true
		} else if m.Digits[i] < m2.Digits[i] {
			return false
		}
	}

	return true
}

func MIntFromString(s string) MInt {
	var m MInt
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			m.Digits = append(m.Digits, s[i]-'0')
		} else if s[i] >= 'a' && s[i] <= 'f' {
			m.Digits = append(m.Digits, s[i]-'a'+10)
		} else {
			panic("Invalid digit in string:" + s)
		}
	}
	return m
}

func (m MInt) ToString() string {
	var s string
	for i := 0; i < len(m.Digits); i++ {
		if m.Digits[i] < 10 {
			s += string(m.Digits[i] + '0')
		} else {
			s += string(m.Digits[i] - 10 + 'a')
		}
	}
	return s
}

func Rand() MInt {
	var m MInt
	nbDigits := rand.Intn(100)
	m.Digits = make([]uint8, nbDigits)
	for i := 0; i < nbDigits; i++ {
		m.Digits[i] = uint8(rand.Intn(16))
	}
	return m
}
