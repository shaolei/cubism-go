package physics

import "math"

func sqrtFloat32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func cosFloat32(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func sinFloat32(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func atan2Float32(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func absFloat32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// degreesToRadian converts degrees to radians.
func degreesToRadian(degrees float32) float32 {
	return degrees / 180.0 * float32(math.Pi)
}

// directionToRadian calculates the radian angle from one direction to another.
// Matches CubismMath::DirectionToRadian in the official SDK.
func directionToRadian(from, to Vector2) float32 {
	q1 := atan2Float32(to.Y, to.X)
	q2 := atan2Float32(from.Y, from.X)
	ret := q1 - q2

	for ret < -float32(math.Pi) {
		ret += float32(math.Pi) * 2.0
	}
	for ret > float32(math.Pi) {
		ret -= float32(math.Pi) * 2.0
	}

	return ret
}

// radianToDirection converts a radian angle to a direction vector.
// Matches CubismMath::RadianToDirection in the official SDK.
func radianToDirection(totalAngle float32) Vector2 {
	return Vector2{
		X: sinFloat32(totalAngle),
		Y: cosFloat32(totalAngle),
	}
}
