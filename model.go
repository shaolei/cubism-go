package cubism

import (
	"fmt"

	"github.com/shaolei/cubism-go/internal/breath"
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/shaolei/cubism-go/internal/cubismmodel"
	"github.com/shaolei/cubism-go/internal/cubismusermodel"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/look"
	"github.com/shaolei/cubism-go/internal/model"
	"github.com/shaolei/cubism-go/internal/motion"
	"github.com/shaolei/cubism-go/internal/physics"
	"github.com/shaolei/cubism-go/internal/pose"
)

// Model delegates to CubismUserModel (composition manager) and CubismModel (pure data layer),
// matching the official SDK architecture.
type Model struct {
	inner *cubismusermodel.CubismUserModel
}

// Get the version of the model
func (m *Model) GetVersion() int {
	return m.inner.GetModel().Version()
}

// Get the core
func (m *Model) GetCore() core.Core {
	return m.inner.GetModel().Core()
}

// Get the moc
func (m *Model) GetMoc() moc.Moc {
	return m.inner.GetModel().Moc()
}

// Get the opacity of the model
func (m *Model) GetOpacity() float32 {
	return m.inner.GetModel().Opacity()
}

// Get the path of a texture image
func (m *Model) GetTextures() []string {
	return m.inner.GetModel().Textures()
}

// Get the sorted drawing order indices
func (m *Model) GetSortedIndices() []int {
	return m.inner.GetModel().SortedIndices()
}

// Get the drawables
func (m *Model) GetDrawables() []Drawable {
	infos := m.inner.GetModel().Drawables()
	result := make([]Drawable, len(infos))
	for i, info := range infos {
		result[i] = Drawable{
			Id:              info.Id,
			Texture:         info.Texture,
			VertexPositions: info.VertexPositions,
			VertexUvs:       info.VertexUvs,
			VertexIndices:   info.VertexIndices,
			ConstantFlag:    info.ConstantFlag,
			DynamicFlag:     info.DynamicFlag,
			Opacity:         info.Opacity,
			Masks:           info.Masks,
		}
	}
	return result
}

// Get the Drawable with the specified ID
func (m *Model) GetDrawable(id string) (d Drawable, err error) {
	dm := m.inner.GetModel().DrawablesMap()
	if info, ok := dm[id]; ok {
		d = Drawable{
			Id:              info.Id,
			Texture:         info.Texture,
			VertexPositions: info.VertexPositions,
			VertexUvs:       info.VertexUvs,
			VertexIndices:   info.VertexIndices,
			ConstantFlag:    info.ConstantFlag,
			DynamicFlag:     info.DynamicFlag,
			Opacity:         info.Opacity,
			Masks:           info.Masks,
		}
		return
	}
	err = fmt.Errorf("Drawable not found: %s", id)
	return
}

// Get the list of hit areas
func (m *Model) GetHitAreas() []model.HitArea {
	return m.inner.GetModel().HitAreas()
}

// Close releases the resources held by the Model.
// After calling Close, the Model must not be used anymore.
func (m *Model) Close() {
	m.inner.GetModel().Close()
}

// Get the list of parameters
func (m *Model) GetParameters() []parameter.Parameter {
	return m.inner.GetModel().GetParameters()
}

// Get the value of the parameter
func (m *Model) GetParameterValue(id string) float32 {
	return m.inner.GetModel().GetParameterValue(id)
}

// Set the value of the parameter
func (m *Model) SetParameterValue(id string, value float32) {
	m.inner.GetModel().SetParameterValue(id, value)
}

// AddParameterValue adds a value to the parameter with the given weight.
// Equivalent to: SetParameterValue(id, GetParameterValue(id) + value * weight)
// Matches the official SDK's AddParameterValue behavior.
func (m *Model) AddParameterValue(id string, value float32, weight float32) {
	m.inner.GetModel().AddParameterValue(id, value, weight)
}

// MultiplyParameterValue multiplies the parameter value by (1 + (value-1)*weight).
// Equivalent to: SetParameterValue(id, GetParameterValue(id) * (1 + (value-1)*weight))
// Matches the official SDK's MultiplyParameterValue behavior.
func (m *Model) MultiplyParameterValue(id string, value float32, weight float32) {
	m.inner.GetModel().MultiplyParameterValue(id, value, weight)
}

// SetParameterValueWithWeight sets the parameter value with weight interpolation.
// Equivalent to: SetParameterValue(id, GetParameterValue(id) + (value - GetParameterValue(id)) * weight)
// Matches the official SDK's SetParameterValue with weight behavior.
func (m *Model) SetParameterValueWithWeight(id string, value float32, weight float32) {
	m.inner.GetModel().SetParameterValueWithWeight(id, value, weight)
}

// SaveParameters saves the current parameter values for later restoration.
// This is used by the motion system to save the pre-motion state before
// applying motion curves, allowing multiple motions to blend correctly.
func (m *Model) SaveParameters() {
	m.inner.GetModel().SaveParameters()
}

// LoadParameters restores the previously saved parameter values.
// Used by the motion system to reset parameters before applying the
// current frame's motion updates.
func (m *Model) LoadParameters() {
	m.inner.GetModel().LoadParameters()
}

// Get the list of motion group names
func (m *Model) GetMotionGroupNames() (names []string) {
	for k := range m.inner.Motions() {
		names = append(names, k)
	}
	return
}

// Get the list of motion group names
func (m *Model) GetMotionGroups() []string {
	return m.inner.GetMotionGroups()
}

// Get the list of motions in the group
func (m *Model) GetMotions(groupName string) []motion.Motion {
	return m.inner.GetMotions(groupName)
}

// Play a motion
func (m *Model) PlayMotion(groupName string, index int, loop bool) (id int) {
	return m.inner.PlayMotion(groupName, index, loop)
}

// Stop a motion
func (m *Model) StopMotion(id int) {
	m.inner.StopMotion(id)
}

// Enable Auto Blink
func (m *Model) EnableAutoBlink() {
	m.inner.EnableAutoBlink()
}

// Disable Auto Blink
func (m *Model) DisableAutoBlink() {
	m.inner.DisableAutoBlink()
}

// EnableBreath enables the breathing effect with default parameters.
// The breathing effect uses sine waves to create subtle body movement.
func (m *Model) EnableBreath() {
	m.inner.EnableBreath()
}

// EnableBreathWithParameters enables the breathing effect with custom parameters.
func (m *Model) EnableBreathWithParameters(params []breath.BreathParameterData) {
	m.inner.EnableBreathWithParameters(params)
}

// DisableBreath disables the breathing effect.
func (m *Model) DisableBreath() {
	m.inner.DisableBreath()
}

// EnableLook enables the look/eye-tracking effect with the given parameters.
// dragX and dragY should be in the range [-1.0, 1.0].
func (m *Model) EnableLook(params []look.LookParameterData) {
	m.inner.EnableLook(params)
}

// DisableLook disables the look/eye-tracking effect.
func (m *Model) DisableLook() {
	m.inner.DisableLook()
}

// SetLookTarget sets the target point for the look/eye-tracking effect.
// dragX and dragY should be in the range [-1.0, 1.0], where (0, 0) is center.
func (m *Model) SetLookTarget(dragX, dragY float32) {
	m.inner.SetLookTarget(dragX, dragY)
}

// GetLookTarget returns the current look target point.
func (m *Model) GetLookTarget() (float32, float32) {
	return m.inner.GetLookTarget()
}

// SetPhysicsOptions sets the gravity and wind options for the physics engine.
func (m *Model) SetPhysicsOptions(gravityX, gravityY, windX, windY float32) {
	m.inner.SetPhysicsOptions(gravityX, gravityY, windX, windY)
}

// GetPhysicsOptions returns the current physics gravity and wind settings.
func (m *Model) GetPhysicsOptions() (gravityX, gravityY, windX, windY float32) {
	return m.inner.GetPhysicsOptions()
}

// DisablePhysics disables the physics simulation.
func (m *Model) DisablePhysics() {
	m.inner.DisablePhysics()
}

// Play an expression by name
func (m *Model) PlayExpression(name string) {
	m.inner.PlayExpression(name)
}

// Stop the current expression
func (m *Model) StopExpression() {
	m.inner.StopExpression()
}

// Get the name of the currently playing expression
func (m *Model) GetCurrentExpression() string {
	return m.inner.GetCurrentExpression()
}

// Get the list of available expression names
func (m *Model) GetExpressionNames() []string {
	return m.inner.GetExpressionNames()
}

// Update the model
// The update flow matches the official SDK's CubismUserModel:
// 1. Motion update (always first, before the scheduler)
// 2. Scheduler-driven late updates: blink → expression → look → breath → physics → pose
// 3. Core update
// The blink updater respects motionUpdated — it only applies when motion did NOT update.
func (m *Model) Update(delta float64) {
	m.inner.Update(delta)
}

// ---- Internal access for model loading and renderer ----

// GetCubismUserModel returns the internal CubismUserModel for use by loading code.
func (m *Model) GetCubismUserModel() *cubismusermodel.CubismUserModel {
	return m.inner
}

// GetCoreModelPtr returns the raw model pointer.
func (m *Model) GetCoreModelPtr() uintptr {
	return m.inner.GetModel().ModelPtr()
}

// GetIdManager returns the ID manager.
func (m *Model) GetIdManager() *id.CubismIdManager {
	return m.inner.GetModel().IdManager()
}

// GetPhysicsManager returns the physics manager.
func (m *Model) GetPhysicsManager() *physics.PhysicsManager {
	return m.inner.GetPhysicsManager()
}

// GetPoseManager returns the pose manager.
func (m *Model) GetPoseManager() *pose.PoseManager {
	return m.inner.GetPoseManager()
}

// GetCubismModel returns the internal CubismModel data layer.
func (m *Model) GetCubismModel() *cubismmodel.CubismModel {
	return m.inner.GetModel()
}

// DrawableFromInfo converts a DrawableInfo to a public Drawable.
func DrawableFromInfo(info cubismmodel.DrawableInfo) Drawable {
	return Drawable{
		Id:              info.Id,
		Texture:         info.Texture,
		VertexPositions: info.VertexPositions,
		VertexUvs:       info.VertexUvs,
		VertexIndices:   info.VertexIndices,
		ConstantFlag:    info.ConstantFlag,
		DynamicFlag:     info.DynamicFlag,
		Opacity:         info.Opacity,
		Masks:           info.Masks,
	}
}

// convertDrawableInfos converts internal DrawableInfo slice to public Drawable slice.
func convertDrawableInfos(infos []cubismmodel.DrawableInfo) []Drawable {
	result := make([]Drawable, len(infos))
	for i, info := range infos {
		result[i] = DrawableFromInfo(info)
	}
	return result
}
