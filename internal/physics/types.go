package physics

// Vector2 represents a 2D vector, matching the CubismMath::CubismVector2 type
// used throughout the physics engine.
type Vector2 struct {
	X float32
	Y float32
}

// Add returns the sum of two vectors.
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of two vectors.
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale returns the vector scaled by a scalar.
func (v Vector2) Scale(s float32) Vector2 {
	return Vector2{X: v.X * s, Y: v.Y * s}
}

// Normalize returns the normalized vector. If the vector is zero-length, returns zero.
func (v Vector2) Normalize() Vector2 {
	lenSq := v.X*v.X + v.Y*v.Y
	if lenSq == 0 {
		return Vector2{}
	}
	invLen := 1.0 / sqrtFloat32(lenSq)
	return Vector2{X: v.X * invLen, Y: v.Y * invLen}
}

// Length returns the length of the vector.
func (v Vector2) Length() float32 {
	return sqrtFloat32(v.X*v.X + v.Y*v.Y)
}

// Normalization represents the normalization range for physics parameters.
// Matches CubismPhysicsNormalization in the official SDK.
type Normalization struct {
	Minimum float32
	Maximum float32
	Default float32
}

// PhysicsSource represents the type of physics input/output.
// Matches CubismPhysicsSource in the official SDK.
type PhysicsSource int

const (
	PhysicsSourceX     PhysicsSource = iota // X-axis position
	PhysicsSourceY                          // Y-axis position
	PhysicsSourceAngle                      // Angle
)

// Particle represents a physics particle in the pendulum simulation.
// Matches CubismPhysicsParticle in the official SDK.
type Particle struct {
	InitialPosition Vector2
	Mobility        float32
	Delay           float32
	Acceleration    float32
	Radius          float32
	Position        Vector2
	LastPosition    Vector2
	LastGravity     Vector2
	Force           Vector2
	Velocity        Vector2
}

// Input represents a physics input binding from a model parameter.
// Matches CubismPhysicsInput in the official SDK.
type Input struct {
	SourceId              string
	SourceParameterIndex  int
	Weight                float32
	Type                  PhysicsSource
	Reflect               bool
	NormalizationPosition Normalization
	NormalizationAngle    Normalization
}

// Output represents a physics output binding to a model parameter.
// Matches CubismPhysicsOutput in the official SDK.
type Output struct {
	DestinationId         string
	DestinationParameterIndex int
	VertexIndex           int
	TranslationScale     Vector2
	AngleScale           float32
	Weight               float32
	Type                 PhysicsSource
	Reflect              bool
	ValueBelowMinimum    float32
	ValueExceededMaximum float32
}

// SubRig represents a single physics sub-rig (strand of particles).
// Matches CubismPhysicsSubRig in the official SDK.
// Each sub-rig has its own set of inputs, outputs, and particles,
// referenced by base indices into the flat arrays.
type SubRig struct {
	InputCount        int
	OutputCount       int
	ParticleCount     int
	BaseInputIndex    int
	BaseOutputIndex   int
	BaseParticleIndex int
	NormalizationPosition Normalization
	NormalizationAngle    Normalization
}

// Rig represents the complete physics rig data.
// Matches CubismPhysicsRig in the official SDK.
type Rig struct {
	SubRigCount int
	Settings    []SubRig
	Inputs      []Input
	Outputs     []Output
	Particles   []Particle
	Gravity     Vector2
	Wind        Vector2
	Fps         float32
}

// Options represents externally configurable physics options.
// Matches CubismPhysics::Options in the official SDK.
type Options struct {
	Gravity Vector2
	Wind    Vector2
}

// rigOutput stores the computed output values for a single sub-rig,
// used for frame interpolation.
type rigOutput struct {
	outputs []float32
}
