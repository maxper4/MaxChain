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

func equals(point1 Point, point2 Point) bool {
	return point1.x.Eq(point2.x) && point1.y.Eq(point2.y)
}

func (point Point) Double() Point {
	lambda := (point.x.Multi(3).Mult(point.x).Add(point.curve.a)).Divide(point.y.Multi(2))
	xr := lambda.Mult(lambda).Sub(point.x.Multi(2))
	return Point{xr, lambda.Mult(point.x.Sub(xr)).Sub(point.y), point.curve}
}

func (point Point) Add(point2 Point) Point {
	if equals(point, point2) {
		return point.Double()
	}

	lambda := (point2.y.Sub(point.y)).Divide(point2.x.Sub(point.x))
	xr := lambda.Mult(lambda).Sub(point.x).Sub(point2.x)
	yr := lambda.Mult(point.x.Sub(xr)).Sub(point.y)
	return Point{xr, yr, point.curve}
}

func (point Point) Multiply(k int) Point {
	if k == 1 {
		return point
	}
	if k % 2 == 0 {
		return point.Double().Multiply(k / 2)
	}
	return point.Add(point.Double().Multiply((k - 1) / 2))
}