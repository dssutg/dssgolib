// A simple 2D float-based vector type and associated operations such as
// creation, arithmetic, geometry, and transforms.
package utils

import (
	"fmt"
	"math"
)

// Vector2f represents a point or direction in 2D space using float64
// coordinates.
type Vector2f struct {
	X float64
	Y float64
}

// NewVector2f returns a new Vector2f with the given X and Y components.
func NewVector2f(x, y float64) Vector2f {
	return Vector2f{x, y}
}

// NewZeroVector2f returns the zero vector (0,0).
func NewZeroVector2f() Vector2f {
	return Vector2f{0, 0}
}

// NewNaNVector2f returns a vector whose components are both NaN.
func NewNaNVector2f() Vector2f {
	return Vector2f{math.NaN(), math.NaN()}
}

// NewVector2fFromRadians constructs a unit vector from an angle in radians. The
// X component is cos(radians), the Y component is sin(radians).
func NewVector2fFromRadians(radians float64) Vector2f {
	s, c := math.Sincos(radians)
	return Vector2f{c, s}
}

// IsNaN reports whether either component of the vector is NaN.
func (a Vector2f) IsNaN() bool {
	return math.IsNaN(a.X) || math.IsNaN((a.Y))
}

// IsZero reports whether both X and Y components are exactly zero.
func (a Vector2f) IsZero() bool {
	return a.X == 0 && a.Y == 0
}

// Copy returns a new Vector2f with the same components as a.
func (a Vector2f) Copy() Vector2f {
	return Vector2f{a.X, a.Y}
}

// ToString returns a human-readable representation of the vector.
func (a Vector2f) ToString() string {
	return fmt.Sprintf("Vector2f{%g, %g}", a.X, a.Y)
}

// Add returns the component-wise sum of a and b.
func (a Vector2f) Add(b *Vector2f) Vector2f {
	return Vector2f{a.X + b.X, a.Y + b.Y}
}

// Sub returns the component-wise difference a - b.
func (a Vector2f) Sub(b *Vector2f) Vector2f {
	return Vector2f{a.X - b.X, a.Y - b.Y}
}

// Mul returns the component-wise product of a and b.
func (a Vector2f) Mul(b *Vector2f) Vector2f {
	return Vector2f{a.X * b.X, a.Y * b.Y}
}

// MulN returns a vector whose components are a scaled by scalar n.
func (a Vector2f) MulN(n float64) Vector2f {
	return Vector2f{a.X * n, a.Y * n}
}

// Neg returns the negation of the vector (-X, -Y).
func (a Vector2f) Neg() Vector2f {
	return Vector2f{-a.X, -a.Y}
}

// Length returns the Euclidean norm (magnitude) of the vector.
func (a Vector2f) Length() float64 {
	return math.Hypot(a.X, a.Y)
}

// Sqrlen returns the squared length X^2 + Y^2.
func (a Vector2f) Sqrlen() float64 {
	return a.X*a.X + a.Y*a.Y
}

// Unit returns a unit vector in the same direction as a.
// If a has zero length, it returns a NaN vector.
func (a Vector2f) Unit() Vector2f {
	length := math.Hypot(a.X, a.Y)
	if length == 0 {
		return NewNaNVector2f()
	}
	return Vector2f{a.X / length, a.Y / length}
}

// Dot returns the dot (inner) product of a and b.
func (a Vector2f) Dot(b *Vector2f) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Cross returns the scalar 2D cross product a and b, which is X1*Y2 - Y1*X2.
func (a Vector2f) Cross(b *Vector2f) float64 {
	return a.X*b.Y - a.Y*b.X
}

// PolarAngle returns the angle in radians from the +X axis to the vector.
func (a Vector2f) PolarAngle() float64 {
	return math.Atan2(a.Y, a.X)
}

// Polar returns the polar coordinates (radius, radians) of the vector.
func (a Vector2f) Polar() (radius, radians float64) {
	radius = math.Hypot(a.Y, a.X)
	radians = math.Atan2(a.Y, a.X)
	return
}

// AngleBetween returns the angle in radians between a and b.
// If either vector has zero length, it returns NaN.
func (a Vector2f) AngleBetween(b *Vector2f) float64 {
	dot := a.X*b.X + a.Y*b.Y
	lengths := math.Hypot(a.X, a.Y) * math.Hypot(b.X, b.Y)
	if lengths == 0 {
		return math.NaN()
	}
	return math.Acos(dot / lengths)
}

// DistanceTo returns the Euclidean distance from a to b.
func (a Vector2f) DistanceTo(b *Vector2f) float64 {
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}

// Lerp returns the linear interpolation between a and b by t in [0,1].
func (a Vector2f) Lerp(b *Vector2f, t float64) Vector2f {
	return Vector2f{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t}
}

// Rotate returns the vector a rotated by the given radians about the origin.
func (a Vector2f) Rotate(radians float64) Vector2f {
	s, c := math.Sincos(radians)
	return Vector2f{a.X*c - a.Y*s, a.X*s + a.Y*c}
}

// RotateAround rotates a by radians about the given point.
func (a Vector2f) RotateAround(point *Vector2f, radians float64) Vector2f {
	dx := a.X - point.X
	dy := a.Y - point.Y

	s, c := math.Sincos(radians)

	rx := (dx*c - dy*s) + point.X
	ry := (dx*s + dy*c) + point.Y

	return Vector2f{rx, ry}
}

// ClampLength returns a vector in the same direction as a with length clamped
// to maxLength. If |a| <= maxLength, a is returned unchanged.
func (a Vector2f) ClampLength(maxLength float64) Vector2f {
	length := math.Hypot(a.X, a.Y)

	if length != 0 && length > maxLength {
		ratio := maxLength / length
		return Vector2f{a.X * ratio, a.Y * ratio}
	}

	return Vector2f{a.X, a.Y}
}

// Project returns the projection of a onto the vector 'onto' If 'onto' has.
// zero length, a NaN vector is returned .
func (a Vector2f) Project(onto *Vector2f) Vector2f {
	ontoDot := onto.X*onto.X + onto.Y*onto.Y

	if ontoDot == 0 {
		return NewNaNVector2f()
	}

	dotRatio := (a.X*onto.X + a.Y*onto.Y) / ontoDot

	return Vector2f{onto.X * dotRatio, onto.Y * dotRatio}
}

// Reflect returns the reflection of a around the given normal vector The.
// normal is assumed to be normalized (unit length) .
func (a Vector2f) Reflect(normal Vector2f) Vector2f {
	dot := a.X*normal.X + a.Y*normal.Y
	return Vector2f{a.X - normal.X*2*dot, a.Y - normal.Y*2*dot}
}

// Normal returns a vector perpendicular to a, rotated +90 degrees (-Y, X).
func (a Vector2f) Normal() Vector2f {
	return Vector2f{-a.Y, a.X}
}

// Abs returns a vector whose components are the absolute values of a's components.
func (a Vector2f) Abs() Vector2f {
	return Vector2f{math.Abs(a.X), math.Abs(a.Y)}
}

// Sign returns a vector of the signs of a's components: -1 for negative, +1 for
// positive, and 0 for zero.
func (a Vector2f) Sign() Vector2f {
	return Vector2f{Sign(a.X), Sign(a.Y)}
}

// MinComponent returns the smaller of the X and Y components.
func (a Vector2f) MinComponent() float64 {
	return math.Min(a.X, a.Y)
}

// MaxComponent returns the larger of the X and Y components.
func (a Vector2f) MaxComponent() float64 {
	return math.Max(a.X, a.Y)
}

// Aspect returns the aspect ratio X / Y. If Y is zero, it returns NaN.
func (a Vector2f) Aspect() float64 {
	if a.Y == 0 {
		return math.NaN()
	}
	return a.X / a.Y
}

// Trunc returns a vector whose components are truncated toward zero.
func (a Vector2f) Trunc() Vector2f {
	return Vector2f{math.Trunc(a.X), math.Trunc(a.Y)}
}

// Floor returns a vector whose components are rounded down to the nearest
// integers.
func (a Vector2f) Floor() Vector2f {
	return Vector2f{math.Floor(a.X), math.Floor(a.Y)}
}

// Round returns a vector whose components are rounded to the nearest integers.
func (a Vector2f) Round() Vector2f {
	return Vector2f{math.Round(a.X), math.Round(a.Y)}
}

// Ceil returns a vector whose components are rounded up to the nearest integers.
func (a Vector2f) Ceil() Vector2f {
	return Vector2f{math.Ceil(a.X), math.Ceil(a.Y)}
}

// SnappedRound rounds each component of a to the nearest multiple of the
// corresponding component in by. If either component of by is zero, it returns
// a NaN vector.
func (a Vector2f) SnappedRound(by *Vector2f) Vector2f {
	if by.X == 0 || by.Y == 0 {
		return NewNaNVector2f()
	}
	return Vector2f{
		math.Round(a.X/by.X) * by.X,
		math.Round(a.Y/by.Y) * by.Y,
	}
}

// SnappedFloor floors each component of a to the nearest multiple of the
// corresponding component in by. If either component of by is zero, it returns
// a NaN vector.
func (a Vector2f) SnappedFloor(by *Vector2f) Vector2f {
	if by.X == 0 || by.Y == 0 {
		return NewNaNVector2f()
	}
	return Vector2f{
		math.Floor(a.X/by.X) * by.X,
		math.Floor(a.Y/by.Y) * by.Y,
	}
}

// SnappedCeil ceils each component of a to the nearest multiple of the
// corresponding component in by. If either component of by is zero, it returns
// a NaN vector.
func (a Vector2f) SnappedCeil(by *Vector2f) Vector2f {
	if by.X == 0 || by.Y == 0 {
		return NewNaNVector2f()
	}
	return Vector2f{
		math.Ceil(a.X/by.X) * by.X,
		math.Ceil(a.Y/by.Y) * by.Y,
	}
}
