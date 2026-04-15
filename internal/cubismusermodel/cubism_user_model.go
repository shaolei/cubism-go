package cubismusermodel

import (
	"github.com/shaolei/cubism-go/internal/breath"
	"github.com/shaolei/cubism-go/internal/blink"
	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/cubismmodel"
	"github.com/shaolei/cubism-go/internal/expression"
	"github.com/shaolei/cubism-go/internal/look"
	"github.com/shaolei/cubism-go/internal/motion"
	"github.com/shaolei/cubism-go/internal/physics"
	"github.com/shaolei/cubism-go/internal/pose"
	"github.com/shaolei/cubism-go/internal/updater"
)

// CubismUserModel is the composition manager for a Live2D model.
// It holds the pure data CubismModel and manages all subsystem managers
// (motion, expression, blink, breath, look, physics, pose) and the
// update scheduler.
// This corresponds to CubismUserModel in the official SDK.
type CubismUserModel struct {
	// Pure data layer
	model *cubismmodel.CubismModel

	// Motion management
	motionManager *motion.CubismMotionManager
	loopMotions   []int
	motions       map[string][]motion.Motion

	// Subsystem managers
	blinkManager      *blink.BlinkManager
	poseManager       *pose.PoseManager
	expressionManager *expression.CubismExpressionMotionManager
	breathManager     *breath.BreathManager
	lookManager       *look.LookManager
	physicsManager    *physics.PhysicsManager

	// Update scheduler
	scheduler     *updater.UpdateScheduler
	motionUpdated bool

	// Updaters registered with the scheduler (stored for removal)
	blinkUpdater      *updater.FuncUpdater
	expressionUpdater *updater.FuncUpdater
	breathUpdater     *updater.FuncUpdater
	lookUpdater       *updater.FuncUpdater
	physicsUpdater    *updater.FuncUpdater
	poseUpdater       *updater.FuncUpdater
}

// NewCubismUserModel creates a new CubismUserModel with the given data model.
func NewCubismUserModel(m *cubismmodel.CubismModel) *CubismUserModel {
	return &CubismUserModel{
		model:     m,
		scheduler: updater.NewUpdateScheduler(),
	}
}

// GetModel returns the pure data CubismModel.
func (u *CubismUserModel) GetModel() *cubismmodel.CubismModel {
	return u.model
}

// ---- Motion management ----

// Motions returns the motion map.
func (u *CubismUserModel) Motions() map[string][]motion.Motion {
	return u.motions
}

// GetMotionGroups returns the list of motion group names.
func (u *CubismUserModel) GetMotionGroups() []string {
	groups := make([]string, 0, len(u.motions))
	for name := range u.motions {
		groups = append(groups, name)
	}
	return groups
}

// GetMotions returns the motions in a group.
func (u *CubismUserModel) GetMotions(groupName string) []motion.Motion {
	return u.motions[groupName]
}

// SetMotions sets the motion map.
func (u *CubismUserModel) SetMotions(motions map[string][]motion.Motion) {
	u.motions = motions
}

// PlayMotion plays a motion by group name and index.
func (u *CubismUserModel) PlayMotion(groupName string, index int, loop bool) int {
	if u.motionManager == nil {
		u.motionManager = motion.NewCubismMotionManager(u.model.Core(), u.model.ModelPtr(), func(id int) {
			for _, loopId := range u.loopMotions {
				if id == loopId {
					u.motionManager.Reset(id)
					return
				}
			}
			u.motionManager.Close(id)
		})
		u.motionManager.SetIdManager(u.model.IdManager())
	}
	group := u.motions[groupName]
	if len(group) == 0 || index < 0 || index >= len(group) {
		return -1
	}
	id := u.motionManager.StartMotionWithPriority(group[index], loop, motion.PriorityNormal)
	if loop {
		u.loopMotions = append(u.loopMotions, id)
	}
	return id
}

// StopMotion stops a motion by ID.
func (u *CubismUserModel) StopMotion(id int) {
	for i, loopId := range u.loopMotions {
		if id == loopId {
			u.loopMotions = append(u.loopMotions[:i], u.loopMotions[i+1:]...)
			break
		}
	}
	u.motionManager.Close(id)
}

// ---- Blink management ----

// EnableAutoBlink enables automatic blinking.
func (u *CubismUserModel) EnableAutoBlink() {
	for _, group := range u.model.Groups() {
		if group.Name == "EyeBlink" {
			u.blinkManager = blink.NewBlinkManager(u.model.Core(), u.model.ModelPtr(), group.Ids)
			u.blinkManager.SetIdManager(u.model.IdManager())
			u.blinkUpdater = updater.NewFuncUpdater(updater.UpdateOrderEyeBlink, func(delta float64) {
				if !u.motionUpdated {
					u.blinkManager.Update(delta)
				}
			})
			u.scheduler.AddUpdater(updater.UpdateOrderEyeBlink, u.blinkUpdater)
			return
		}
	}
}

// DisableAutoBlink disables automatic blinking.
func (u *CubismUserModel) DisableAutoBlink() {
	if u.blinkUpdater != nil {
		u.scheduler.RemoveUpdater(u.blinkUpdater)
		u.blinkUpdater = nil
	}
	u.blinkManager = nil
}

// ---- Breath management ----

// EnableBreath enables the breathing effect with default parameters.
func (u *CubismUserModel) EnableBreath() {
	u.breathManager = breath.NewBreathManager(u.model.Core(), u.model.ModelPtr())
	u.breathManager.SetIdManager(u.model.IdManager())
	u.breathUpdater = updater.NewFuncUpdater(updater.UpdateOrderBreath, func(delta float64) {
		u.breathManager.Update(delta)
	})
	u.scheduler.AddUpdater(updater.UpdateOrderBreath, u.breathUpdater)
}

// EnableBreathWithParameters enables the breathing effect with custom parameters.
func (u *CubismUserModel) EnableBreathWithParameters(params []breath.BreathParameterData) {
	u.breathManager = breath.NewBreathManager(u.model.Core(), u.model.ModelPtr())
	u.breathManager.SetParameters(params)
	u.breathManager.SetIdManager(u.model.IdManager())
	if u.breathUpdater == nil {
		u.breathUpdater = updater.NewFuncUpdater(updater.UpdateOrderBreath, func(delta float64) {
			u.breathManager.Update(delta)
		})
		u.scheduler.AddUpdater(updater.UpdateOrderBreath, u.breathUpdater)
	}
}

// DisableBreath disables the breathing effect.
func (u *CubismUserModel) DisableBreath() {
	if u.breathUpdater != nil {
		u.scheduler.RemoveUpdater(u.breathUpdater)
		u.breathUpdater = nil
	}
	u.breathManager = nil
}

// ---- Look management ----

// EnableLook enables the look/eye-tracking effect.
func (u *CubismUserModel) EnableLook(params []look.LookParameterData) {
	u.lookManager = look.NewLookManager(u.model.Core(), u.model.ModelPtr())
	u.lookManager.SetParameters(params)
	u.lookManager.SetIdManager(u.model.IdManager())
	u.lookUpdater = updater.NewFuncUpdater(updater.UpdateOrderLook, func(delta float64) {
		u.lookManager.Update(delta)
	})
	u.scheduler.AddUpdater(updater.UpdateOrderLook, u.lookUpdater)
}

// DisableLook disables the look/eye-tracking effect.
func (u *CubismUserModel) DisableLook() {
	if u.lookUpdater != nil {
		u.scheduler.RemoveUpdater(u.lookUpdater)
		u.lookUpdater = nil
	}
	u.lookManager = nil
}

// SetLookTarget sets the target point for look/eye-tracking.
func (u *CubismUserModel) SetLookTarget(dragX, dragY float32) {
	if u.lookManager != nil {
		u.lookManager.SetTarget(dragX, dragY)
	}
}

// GetLookTarget returns the current look target point.
func (u *CubismUserModel) GetLookTarget() (float32, float32) {
	if u.lookManager == nil {
		return 0, 0
	}
	return u.lookManager.GetTarget()
}

// ---- Physics management ----

// SetPhysicsOptions sets the gravity and wind options.
func (u *CubismUserModel) SetPhysicsOptions(gravityX, gravityY, windX, windY float32) {
	if u.physicsManager != nil {
		u.physicsManager.SetOptions(physics.Options{
			Gravity: physics.Vector2{X: gravityX, Y: gravityY},
			Wind:    physics.Vector2{X: windX, Y: windY},
		})
	}
}

// GetPhysicsOptions returns the current physics settings.
func (u *CubismUserModel) GetPhysicsOptions() (gravityX, gravityY, windX, windY float32) {
	if u.physicsManager == nil {
		return 0, 0, 0, 0
	}
	opts := u.physicsManager.GetOptions()
	return opts.Gravity.X, opts.Gravity.Y, opts.Wind.X, opts.Wind.Y
}

// DisablePhysics disables the physics simulation.
func (u *CubismUserModel) DisablePhysics() {
	if u.physicsUpdater != nil {
		u.scheduler.RemoveUpdater(u.physicsUpdater)
		u.physicsUpdater = nil
	}
	u.physicsManager = nil
}

// SetPhysicsManager sets the physics manager (called during model loading).
func (u *CubismUserModel) SetPhysicsManager(pm *physics.PhysicsManager) {
	u.physicsManager = pm
	if pm != nil {
		u.physicsUpdater = updater.NewFuncUpdater(updater.UpdateOrderPhysics, func(delta float64) {
			u.physicsManager.Evaluate(delta)
		})
		u.scheduler.AddUpdater(updater.UpdateOrderPhysics, u.physicsUpdater)
	}
}

// GetPhysicsManager returns the physics manager.
func (u *CubismUserModel) GetPhysicsManager() *physics.PhysicsManager {
	return u.physicsManager
}

// SetPoseManager sets the pose manager (called during model loading).
func (u *CubismUserModel) SetPoseManager(pm *pose.PoseManager) {
	u.poseManager = pm
	if pm != nil {
		u.poseUpdater = updater.NewFuncUpdater(updater.UpdateOrderPose, func(delta float64) {
			u.poseManager.Update(u.model.Core(), u.model.ModelPtr(), delta)
		})
		u.scheduler.AddUpdater(updater.UpdateOrderPose, u.poseUpdater)
	}
}

// GetPoseManager returns the pose manager.
func (u *CubismUserModel) GetPoseManager() *pose.PoseManager {
	return u.poseManager
}

// ---- Expression management ----

// PlayExpression plays an expression by name.
func (u *CubismUserModel) PlayExpression(name string) {
	if u.expressionManager == nil {
		u.expressionManager = expression.NewCubismExpressionMotionManager(u.model.Expressions())
		u.expressionManager.InitWithCore(u.model.Core(), u.model.ModelPtr())
		u.expressionManager.SetIdManager(u.model.IdManager())
		u.expressionUpdater = updater.NewFuncUpdater(updater.UpdateOrderExpression, func(delta float64) {
			u.expressionManager.Update(u.model.Core(), u.model.ModelPtr(), delta)
		})
		u.scheduler.AddUpdater(updater.UpdateOrderExpression, u.expressionUpdater)
	}
	u.expressionManager.PlayExpression(name)
}

// StopExpression stops the current expression.
func (u *CubismUserModel) StopExpression() {
	if u.expressionManager != nil {
		u.expressionManager.StopExpression()
	}
}

// GetCurrentExpression returns the name of the currently playing expression.
func (u *CubismUserModel) GetCurrentExpression() string {
	if u.expressionManager == nil {
		return ""
	}
	return u.expressionManager.GetCurrentExpression()
}

// GetExpressionNames returns the list of available expression names.
func (u *CubismUserModel) GetExpressionNames() []string {
	if u.expressionManager == nil {
		u.expressionManager = expression.NewCubismExpressionMotionManager(u.model.Expressions())
		u.expressionManager.InitWithCore(u.model.Core(), u.model.ModelPtr())
		u.expressionManager.SetIdManager(u.model.IdManager())
	}
	return u.expressionManager.GetExpressionNames()
}

// ---- Update loop ----

// Update performs the full model update cycle.
// The update flow matches the official SDK's CubismUserModel:
// 1. Motion update (always first, before the scheduler)
// 2. Scheduler-driven late updates: blink → expression → look → breath → physics → pose
// 3. Core update
// The blink updater respects motionUpdated — it only applies when motion did NOT update.
func (u *CubismUserModel) Update(delta float64) {
	// 1. Motion update: applies motion curves to parameters (always first)
	u.motionUpdated = false
	if u.motionManager != nil {
		u.motionManager.Update(delta)
		u.motionUpdated = !u.motionManager.IsFinished()
	}

	// 2. Late updates via scheduler: subsystems run in priority order
	u.scheduler.OnLateUpdate(delta)

	// 3. Core update: processes all accumulated parameter changes
	u.model.Update()

	// 4. Update drawable caches from core
	dfs := u.model.GetDynamicFlags()

	drawOrderDidChange := false
	renderOrderDidChange := false
	opacityDidChange := false
	vertexPositionsDidChange := false

	drawables := u.model.Drawables()
	for i := range drawables {
		drawables[i].DynamicFlag = dfs[i]
		if dfs[i].DrawOrderDidChange {
			drawOrderDidChange = true
		}
		if dfs[i].RenderOrderDidChange {
			renderOrderDidChange = true
		}
		if dfs[i].OpacityDidChange {
			opacityDidChange = true
		}
		if dfs[i].VertexPositionsDidChange {
			vertexPositionsDidChange = true
		}
	}

	if drawOrderDidChange || renderOrderDidChange {
		u.model.SetSortedIndices(u.model.GetSortedDrawableIndices())
	}

	var opacities []float32
	if opacityDidChange {
		opacities = u.model.GetOpacities()
	}

	var vertexPositions [][]drawable.Vector2
	if vertexPositionsDidChange {
		vertexPositions = u.model.GetVertexPositions()
	}

	for i := range drawables {
		if opacityDidChange {
			drawables[i].Opacity = opacities[i]
		}
		if vertexPositionsDidChange {
			drawables[i].VertexPositions = vertexPositions[i]
		}
	}
}
