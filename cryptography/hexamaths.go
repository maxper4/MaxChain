package cryptography

type mInt struct {
	digits []uint8
}

func (m mInt) getDigits() []uint8 {
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

func mIntFromString(s string) mInt {
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