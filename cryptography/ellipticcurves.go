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
	return point1.x == point2.x && point1.y == point2.y
}

func (point Point) Double() Point {
	lambda := (3 * point.x * point.x + point.curve.a) / (2 * point.y)
	xr := lambda * lambda - 2 * point.x
	return Point{xr, lambda * (point.x - xr) - point.y, point.curve}
}

func (point Point) Add(point2 Point) Point {
	if equals(point, point2) {
		return point.Double()
	}

	lambda := (point2.y - point.y) / (point2.x - point.x)
	xr := lambda * lambda - point.x - point2.x
	yr := lambda * (point.x - xr) - point.y
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