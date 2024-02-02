package cryptography

type EllipticCurve struct {
	a mInt
	b mInt
	n mInt
	g Point
}

type Point struct {
	x mInt
	y mInt
	curve *EllipticCurve
}

func (point Point) ToString() string {
	return "Point(x: " + point.x.ToString() + ", y: " + point.y.ToString() + ")"
}

func equals(point1 Point, point2 Point) bool {
	return point1.x.Eq(point2.x) && point1.y.Eq(point2.y)
}

func (point Point) Double() Point {
	lambda := (point.x.Multi(3).Mult(point.x).Add(point.curve.a)).Divide(point.y.Multi(2))
	xr := lambda.Mult(lambda).SubMod(point.x.Multi(2), point.curve.n)
	return Point{xr, lambda.Mult(point.x.SubMod(xr, point.curve.n)).SubMod(point.y, point.curve.n), point.curve}
}

func (point Point) Add(point2 Point) Point {
	if equals(point, point2) {
		return point.Double()
	}

	lambda := (point2.y.SubMod(point.y, point.curve.n)).Divide(point2.x.SubMod(point.x, point.curve.n))
	xr := lambda.Mult(lambda).SubMod(point.x, point.curve.n).SubMod(point2.x, point.curve.n)
	yr := lambda.Mult(point.x.SubMod(xr, point.curve.n)).SubMod(point.y, point.curve.n)
	return Point{xr, yr, point.curve}
}

func (point Point) Multiply(k mInt) Point {
	if k.Eq(MIntFromString("1")) {
		return point
	}
	if k.Modi(2).Eq(MIntFromString("0")) {
		return point.Double().Multiply(k.Div(2))
	}
	return point.Add(point.Double().Multiply((k.Sub(MIntFromString("1"))).Div(2)))
}