package cryptography

func EncodeMsgToMInt(s string) MInt {
	grd := MIntFromString("1")
	nb := MIntFromString("0")
	b := MIntFromInt(256)

	for i := 0; i < len(s); i++ {
		nb = nb.Add(MIntFromInt(int(s[i])).Mult(grd))

		grd = grd.Mult(b)
	}
	return nb
}

func DecodeMIntToMsg(nb MInt) string {
	s := ""
	zero := MIntFromString("0")
	for !nb.Eq(zero) {
		temp := nb.Mod256()
		s = s + string(temp.ToInt())
		nb = nb.Sub(temp).Divide256()
	}
	return s
}
