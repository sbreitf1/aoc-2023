package helper

type Point2D[T Number] struct {
	X, Y T
}

func (p Point2D[T]) Add(p2 Point2D[T]) Point2D[T] {
	return Point2D[T]{X: p.X + p2.X, Y: p.Y + p2.Y}
}

func (p Point2D[T]) Sub(p2 Point2D[T]) Point2D[T] {
	return Point2D[T]{X: p.X - p2.X, Y: p.Y - p2.Y}
}

func (p Point2D[T]) Neg() Point2D[T] {
	return Point2D[T]{X: -p.X, Y: -p.Y}
}

func (p Point2D[T]) Mul(factor T) Point2D[T] {
	return Point2D[T]{X: p.X * factor, Y: p.Y * factor}
}

func (p Point2D[T]) Div(divisor T) Point2D[T] {
	return Point2D[T]{X: p.X / divisor, Y: p.Y / divisor}
}

func (p Point2D[T]) Cross(p2 Point2D[T]) float64 {
	return float64(p.X)*float64(p2.Y) - float64(p.Y)*float64(p2.X)
}

func (p Point2D[T]) InBounds(min, max Point2D[T]) bool {
	return p.X >= min.X && p.Y >= min.Y && p.X <= max.X && p.Y <= max.Y
}

func ConvertPoint2D[FROM, TO Number](from Point2D[FROM]) Point2D[TO] {
	return Point2D[TO]{X: TO(from.X), Y: TO(from.Y)}
}

type Point3D[T Number] struct {
	X, Y, Z T
}

func (p Point3D[T]) Add(p2 Point3D[T]) Point3D[T] {
	return Point3D[T]{X: p.X + p2.X, Y: p.Y + p2.Y, Z: p.Z + p2.Z}
}

func (p Point3D[T]) Sub(p2 Point3D[T]) Point3D[T] {
	return Point3D[T]{X: p.X - p2.X, Y: p.Y - p2.Y, Z: p.Z - p2.Z}
}

func (p Point3D[T]) Neg() Point3D[T] {
	return Point3D[T]{X: -p.X, Y: -p.Y, Z: -p.Z}
}

func (p Point3D[T]) Mul(factor T) Point3D[T] {
	return Point3D[T]{X: p.X * factor, Y: p.Y * factor, Z: p.Z * factor}
}

func (p Point3D[T]) XY() Point2D[T] {
	return Point2D[T]{X: p.X, Y: p.Y}
}
