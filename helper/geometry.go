package helper

type Point2D struct {
	X, Y int
}

func (p Point2D) Add(p2 Point2D) Point2D {
	return Point2D{X: p.X + p2.X, Y: p.Y + p2.Y}
}

func (p Point2D) Sub(p2 Point2D) Point2D {
	return Point2D{X: p.X - p2.X, Y: p.Y - p2.Y}
}

func (p Point2D) Neg() Point2D {
	return Point2D{X: -p.X, Y: -p.Y}
}

func (p Point2D) Mul(factor int) Point2D {
	return Point2D{X: p.X * factor, Y: p.Y * factor}
}
