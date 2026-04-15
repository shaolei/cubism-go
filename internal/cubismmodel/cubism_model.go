package cubismmodel

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
)

// CubismModel is the pure data layer for a Live2D model.
// It wraps the core engine's model pointer and provides parameter read/write,
// canvas info, drawable data, and Save/Load functionality.
// This corresponds to CubismModel in the official SDK.
type CubismModel struct {
	core      core.Core
	moc       moc.Moc
	idManager *id.CubismIdManager

	// Model metadata
	opacity    float32
	version    int
	textures   []string
	sortedIndices []int
	drawables  []DrawableInfo
	drawablesMap map[string]DrawableInfo
	hitAreas   []model.HitArea
	groups     []model.Group

	// JSON config data (needed for lazy initialization of subsystems)
	physics  model.PhysicsJson
	pose     model.PoseJson
	cdi      model.CdiJson
	exps     []model.ExpJson
	userdata model.UserDataJson

	// Parameter save/restore
	savedParameters map[string]float32
}

// DrawableInfo holds cached drawable data for the model.
type DrawableInfo struct {
	Id              string
	Texture         string
	VertexPositions []drawable.Vector2
	VertexUvs       []drawable.Vector2
	VertexIndices   []uint16
	ConstantFlag    drawable.ConstantFlag
	DynamicFlag     drawable.DynamicFlag
	Opacity         float32
	Masks           []int32
}

// NewCubismModel creates a new CubismModel from core components.
func NewCubismModel(c core.Core, m moc.Moc, idMgr *id.CubismIdManager) *CubismModel {
	return &CubismModel{
		core:      c,
		moc:       m,
		idManager: idMgr,
		opacity:   1.0,
	}
}

// ---- Core access ----

// Core returns the core engine instance.
func (m *CubismModel) Core() core.Core {
	return m.core
}

// Moc returns the moc data.
func (m *CubismModel) Moc() moc.Moc {
	return m.moc
}

// ModelPtr returns the raw model pointer for core operations.
func (m *CubismModel) ModelPtr() uintptr {
	return m.moc.ModelPtr
}

// IdManager returns the ID manager for fast parameter/part access.
func (m *CubismModel) IdManager() *id.CubismIdManager {
	return m.idManager
}

// ---- Metadata access ----

// Opacity returns the model opacity.
func (m *CubismModel) Opacity() float32 {
	return m.opacity
}

// SetOpacity sets the model opacity.
func (m *CubismModel) SetOpacity(opacity float32) {
	m.opacity = opacity
}

// Version returns the model version.
func (m *CubismModel) Version() int {
	return m.version
}

// SetVersion sets the model version.
func (m *CubismModel) SetVersion(version int) {
	m.version = version
}

// Textures returns the texture paths.
func (m *CubismModel) Textures() []string {
	return m.textures
}

// SetTextures sets the texture paths.
func (m *CubismModel) SetTextures(textures []string) {
	m.textures = textures
}

// SortedIndices returns the sorted drawing order indices.
func (m *CubismModel) SortedIndices() []int {
	return m.sortedIndices
}

// SetSortedIndices sets the sorted drawing order indices.
func (m *CubismModel) SetSortedIndices(indices []int) {
	m.sortedIndices = indices
}

// Drawables returns the drawable info list.
func (m *CubismModel) Drawables() []DrawableInfo {
	return m.drawables
}

// SetDrawables sets the drawable info list.
func (m *CubismModel) SetDrawables(drawables []DrawableInfo) {
	m.drawables = drawables
}

// DrawablesMap returns the drawable map by ID.
func (m *CubismModel) DrawablesMap() map[string]DrawableInfo {
	return m.drawablesMap
}

// SetDrawablesMap sets the drawable map by ID.
func (m *CubismModel) SetDrawablesMap(dm map[string]DrawableInfo) {
	m.drawablesMap = dm
}

// HitAreas returns the hit areas.
func (m *CubismModel) HitAreas() []model.HitArea {
	return m.hitAreas
}

// SetHitAreas sets the hit areas.
func (m *CubismModel) SetHitAreas(areas []model.HitArea) {
	m.hitAreas = areas
}

// Groups returns the parameter groups.
func (m *CubismModel) Groups() []model.Group {
	return m.groups
}

// SetGroups sets the parameter groups.
func (m *CubismModel) SetGroups(groups []model.Group) {
	m.groups = groups
}

// ---- JSON config data ----

// Physics returns the physics JSON config.
func (m *CubismModel) Physics() model.PhysicsJson {
	return m.physics
}

// SetPhysics sets the physics JSON config.
func (m *CubismModel) SetPhysics(p model.PhysicsJson) {
	m.physics = p
}

// Pose returns the pose JSON config.
func (m *CubismModel) Pose() model.PoseJson {
	return m.pose
}

// SetPose sets the pose JSON config.
func (m *CubismModel) SetPose(p model.PoseJson) {
	m.pose = p
}

// Cdi returns the CDI JSON config.
func (m *CubismModel) Cdi() model.CdiJson {
	return m.cdi
}

// SetCdi sets the CDI JSON config.
func (m *CubismModel) SetCdi(c model.CdiJson) {
	m.cdi = c
}

// Expressions returns the expression JSON configs.
func (m *CubismModel) Expressions() []model.ExpJson {
	return m.exps
}

// SetExpressions sets the expression JSON configs.
func (m *CubismModel) SetExpressions(e []model.ExpJson) {
	m.exps = e
}

// UserData returns the user data JSON config.
func (m *CubismModel) UserData() model.UserDataJson {
	return m.userdata
}

// SetUserData sets the user data JSON config.
func (m *CubismModel) SetUserData(u model.UserDataJson) {
	m.userdata = u
}

// ---- Parameter operations ----

// GetParameters returns the list of parameters.
func (m *CubismModel) GetParameters() []parameter.Parameter {
	return m.core.GetParameters(m.moc.ModelPtr)
}

// GetParameterValue returns the value of a parameter by ID.
func (m *CubismModel) GetParameterValue(id string) float32 {
	return m.core.GetParameterValue(m.moc.ModelPtr, id)
}

// SetParameterValue sets the value of a parameter by ID.
func (m *CubismModel) SetParameterValue(id string, value float32) {
	m.core.SetParameterValue(m.moc.ModelPtr, id, value)
}

// AddParameterValue adds a value to the parameter with weight.
func (m *CubismModel) AddParameterValue(id string, value float32, weight float32) {
	handle := m.idManager.GetParameterId(id)
	if !handle.IsValid() {
		return
	}
	current := m.core.GetParameterValueByIndex(m.moc.ModelPtr, int(handle))
	m.core.SetParameterValueByIndex(m.moc.ModelPtr, int(handle), current+value*weight)
}

// MultiplyParameterValue multiplies the parameter value.
func (m *CubismModel) MultiplyParameterValue(id string, value float32, weight float32) {
	handle := m.idManager.GetParameterId(id)
	if !handle.IsValid() {
		return
	}
	current := m.core.GetParameterValueByIndex(m.moc.ModelPtr, int(handle))
	m.core.SetParameterValueByIndex(m.moc.ModelPtr, int(handle), current*(1.0+(value-1.0)*weight))
}

// SetParameterValueWithWeight sets the parameter value with weight interpolation.
func (m *CubismModel) SetParameterValueWithWeight(id string, value float32, weight float32) {
	handle := m.idManager.GetParameterId(id)
	if !handle.IsValid() {
		return
	}
	current := m.core.GetParameterValueByIndex(m.moc.ModelPtr, int(handle))
	m.core.SetParameterValueByIndex(m.moc.ModelPtr, int(handle), current+(value-current)*weight)
}

// SaveParameters saves the current parameter values for later restoration.
func (m *CubismModel) SaveParameters() {
	parameters := m.core.GetParameters(m.moc.ModelPtr)
	if m.savedParameters == nil {
		m.savedParameters = make(map[string]float32, len(parameters))
	}
	for _, parameter := range parameters {
		m.savedParameters[parameter.Id] = parameter.Current
	}
}

// LoadParameters restores the previously saved parameter values.
func (m *CubismModel) LoadParameters() {
	if m.savedParameters == nil {
		return
	}
	for id, value := range m.savedParameters {
		handle := m.idManager.GetParameterId(id)
		if handle.IsValid() {
			m.core.SetParameterValueByIndex(m.moc.ModelPtr, int(handle), value)
		}
	}
}

// ---- Core update operations ----

// Update processes all accumulated parameter changes via the core engine.
func (m *CubismModel) Update() {
	m.core.Update(m.moc.ModelPtr)
}

// GetDynamicFlags returns the dynamic flags for all drawables.
func (m *CubismModel) GetDynamicFlags() []drawable.DynamicFlag {
	return m.core.GetDynamicFlags(m.moc.ModelPtr)
}

// GetSortedDrawableIndices returns the sorted drawable indices.
func (m *CubismModel) GetSortedDrawableIndices() []int {
	return m.core.GetSortedDrawableIndices(m.moc.ModelPtr)
}

// GetOpacities returns the opacities for all drawables.
func (m *CubismModel) GetOpacities() []float32 {
	return m.core.GetOpacities(m.moc.ModelPtr)
}

// GetVertexPositions returns the vertex positions for all drawables.
func (m *CubismModel) GetVertexPositions() [][]drawable.Vector2 {
	return m.core.GetVertexPositions(m.moc.ModelPtr)
}

// Close releases the resources held by the model.
func (m *CubismModel) Close() {
	m.moc.Close()
}
