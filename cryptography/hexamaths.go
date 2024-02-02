package cryptography

import (
	"strconv"
)

type mInt struct {
	digits []uint8
}

func (m mInt) getDigits() []uint8 {	// reverse digits: less significant first
	digits := make([]uint8, len(m.digits))
	for i := len(m.digits) - 1; i >= 0; i-- {
		digits[len(m.digits) - 1 - i] = m.digits[i]
	}
	return digits
}

func (m mInt) Add(m2 mInt) mInt {
	var m3 mInt
	digitsNb := len(m.digits)
	if len(m2.digits) < digitsNb {
		digitsNb = len(m2.digits)
	}
	m3.digits = make([]uint8, digitsNb)
	mdigits := m.getDigits()
	m2digits := m2.getDigits()

	for i := 0; i < digitsNb; i++ {
		carry := (m3.digits[i] + mdigits[i] + m2digits[i]) / 16
		m3.digits[i] = (m3.digits[i] + mdigits[i] + m2digits[i]) % 16

		j := i + 1
		for carry > 0 {
			if j == len(m3.digits) {
				m3.digits = append(m3.digits, 0)
			}
			m3.digits[j] += carry
			carry = m3.digits[j] / 16
			m3.digits[j] = m3.digits[j] % 16
			j++
		}
	}
	if len(mdigits) < len(m2digits) {
		for i := len(mdigits); i < len(m2digits); i++ {
			if i >= len(m3.digits) {
				m3.digits = append(m3.digits, m2digits[i])
			} else {
				m3.digits[i] += m2digits[i]
			}
		}
	} else if len(m2digits) < len(mdigits) {
		for i := len(m2digits); i < len(mdigits); i++ {
			if i >= len(m3.digits) {
				m3.digits = append(m3.digits, mdigits[i])
			} else {
				m3.digits[i] += mdigits[i]
			}
		}
	}

	for i, j := 0, len(m3.digits)-1; i < j; i, j = i+1, j-1 {
		m3.digits[i], m3.digits[j] = m3.digits[j], m3.digits[i]
	}
	return m3
}

func (m mInt) Sub(m2 mInt) mInt {
	if !m.GreaterEq(m2) {
		panic("Substraction of a smaller number from a bigger one: " + m.ToString() + " - " + m2.ToString() + " is not possible.")
	}

	var m3 mInt
	mdigits := m.getDigits()
	m2digits := m2.getDigits()
	digitsNb := len(m2digits)

	m3.digits = make([]uint8, len(mdigits))

	for i := 0; i < digitsNb; i++ {
		if mdigits[i] < m2digits[i] {
			mdigits[i] += 16
			mdigits[i+1]--
		}
		m3.digits[i] = mdigits[i] - m2digits[i]
	}

	for i, j := 0, len(m3.digits)-1; i < j; i, j = i+1, j-1 {
		m3.digits[i], m3.digits[j] = m3.digits[j], m3.digits[i]
	}

	for i := 0; i < len(m3.digits); i++ {	// remove leading 0s
		if m3.digits[i] != 0 {
			m3.digits = m3.digits[i:]
			break
		}
	}

	return m3
}

func (m mInt) SubMod(m2 mInt, mod mInt) mInt {
	if !m.GreaterEq(m2) {
		return m.Add(mod).SubMod(m2, mod)
	}
	return m.Sub(m2).Mod(mod)
}

func (m mInt) Multi(i int) mInt {
	return m.Mult(MIntFromString(strconv.Itoa(i)))
}

func (m mInt) Mult(m2 mInt) mInt {
	var m3 mInt
	var multdigitsnb int
	var m2digits = m2.getDigits()
	var mdigits = m.getDigits()

	if len(m.digits) < len(m2.digits) {
		multdigitsnb = len(m.digits)
	} else {
		multdigitsnb = len(m2.digits)
		m2digits, mdigits = mdigits, m2digits
	}

	for i := 0; i < multdigitsnb; i++ {
		var m4 mInt
		m4.digits = make([]uint8, i)
		for j := 0; j < len(m.digits); j++ {
			if j >= len(m4.digits) {
				m4.digits = append(m4.digits, 0)
			}
			a := (m4.digits[j] + mdigits[i] * m2digits[j])
			carry := a / 16
			m4.digits[j] = a % 16
			c := j + 1
			for carry > 0 {
				if c == len(m4.digits) {
					m4.digits = append(m4.digits, 0)
				}
				m4.digits[c] += carry
				carry = m4.digits[c] / 16
				m4.digits[c] = m4.digits[c] % 16
				c++
			}
		}
		for j := 0; j < i; j++ {
			m4.digits = append(m4.digits, 0)
		}
		m4.digits = m4.getDigits()
		for j := 0; j+i < len(m4.digits); j++ {
			m4.digits[j] = m4.digits[j+i]
		}
		for j := len(m4.digits) - 1; j >= len(m4.digits) - i; j-- {
			m4.digits[j] = 0
		}

		m3.digits = m3.getDigits()
		m3 = m3.Add(m4)
		m3.digits = m3.getDigits()
	}
	
	m3.digits = m3.getDigits()
	return m3
}

func (m mInt) Div(i int) mInt {
	return m.Divide(MIntFromString(strconv.Itoa(i)))
}

func (m mInt) Divide(m2 mInt) mInt {
	q := 0
	for m.GreaterEq(m2) {
		m = m.Sub(m2)
		q++
	}
	return MIntFromString(strconv.Itoa(q))
}

func (m mInt) Modi(i int) mInt {
	return m.Mod(MIntFromString(strconv.Itoa(i)))
}

func (m mInt) Mod(m2 mInt) mInt {
	for m.GreaterEq(m2) {
		m = m.Sub(m2)
	}
	return m
}

func (m mInt) Eq(m2 mInt) bool {
	if len(m.digits) != len(m2.digits) {
		return false
	}
	for i := 0; i < len(m.digits); i++ {
		if m.digits[i] != m2.digits[i] {
			return false
		}
	}
	return true
}

func (m mInt) realLen() int {
	for i := 0; i < len(m.digits); i++ {
		if m.digits[i] != 0 {
			return len(m.digits) - i
		}
	}
	return 0
}

func (m mInt) GreaterEq(m2 mInt) bool {
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
		if m.digits[i] > m2.digits[i] {
			return true
		} else if m.digits[i] < m2.digits[i] {
			return false
		}
	}

	return true
}

func MIntFromString(s string) mInt {
	var m mInt
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			m.digits = append(m.digits, s[i] - '0')
		} else if s[i] >= 'a' && s[i] <= 'f'{
			m.digits = append(m.digits, s[i] - 'a' + 10)
		} else {
			panic("Invalid digit in string:" + s)
		}
	}
	return m
}

func (m mInt) ToString() string {
	var s string
	for i := 0; i < len(m.digits); i++ {
		if m.digits[i] < 10 {
			s += string(m.digits[i] + '0')
		} else {
			s += string(m.digits[i] - 10 + 'a')
		}
	}
	return s
}