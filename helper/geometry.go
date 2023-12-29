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

func (p Point2D) Cross(p2 Point2D) float64 {
	return float64(p.X*p2.Y - p.Y*p2.X)
}

type Point3D struct {
	X, Y, Z int
}

func (p Point3D) Add(p2 Point3D) Point3D {
	return Point3D{X: p.X + p2.X, Y: p.Y + p2.Y, Z: p.Z + p2.Z}
}

func (p Point3D) Sub(p2 Point3D) Point3D {
	return Point3D{X: p.X - p2.X, Y: p.Y - p2.Y, Z: p.Z - p2.Z}
}

func (p Point3D) Neg() Point3D {
	return Point3D{X: -p.X, Y: -p.Y, Z: -p.Z}
}

func (p Point3D) Mul(factor int) Point3D {
	return Point3D{X: p.X * factor, Y: p.Y * factor, Z: p.Z * factor}
}

func (p Point3D) XY() Point2D {
	return Point2D{X: p.X, Y: p.Y}
}
