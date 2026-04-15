package cubism

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/cubismmodel"
	"github.com/shaolei/cubism-go/internal/cubismusermodel"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
	"github.com/shaolei/cubism-go/internal/motion"
	"github.com/shaolei/cubism-go/internal/physics"
	"github.com/shaolei/cubism-go/internal/pose"
	"github.com/shaolei/cubism-go/sound"
	"github.com/shaolei/cubism-go/sound/disabled"
)

/*
The main body of cubism-go
*/
type Cubism struct {
	core core.Core
	// A function to load audio files
	LoadSound func(fp string) (s sound.Sound, err error)
}

// Constructor for the [Cubism] struct
func NewCubism(lib string) (c Cubism, err error) {
	c.core, err = core.NewCore(lib)
	return
}

// Load a model from model3.json
func (c *Cubism) LoadModel(path string) (m *Model, err error) {
	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	// Get the directory
	dir := filepath.Dir(absPath)

	// Read model3.json
	buf, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	// Convert to a structure compatible with version 3
	var mj model.ModelJson
	if err = json.Unmarshal(buf, &mj); err != nil {
		return
	}

	// Load the moc3 file
	moc3Path := filepath.Join(dir, mj.FileReferences.Moc)
	mocData, err := c.core.LoadMoc(moc3Path)
	if err != nil {
		return
	}

	// Initialize the ID manager for fast parameter/part access by index
	parameterIds := c.core.GetParameterIds(mocData.ModelPtr)
	partIds := c.core.GetPartIds(mocData.ModelPtr)
	idManager := id.NewCubismIdManager(parameterIds, partIds)

	// Create CubismModel (pure data layer)
	cubismModel := cubismmodel.NewCubismModel(c.core, mocData, idManager)
	cubismModel.SetVersion(mj.Version)
	cubismModel.SetOpacity(1.0)

	// Convert the path of texture image to an absolute path
	textures := mj.FileReferences.Textures
	for i := range textures {
		textures[i] = filepath.Join(dir, textures[i])
	}
	cubismModel.SetTextures(textures)

	cubismModel.SetGroups(mj.Groups)
	cubismModel.SetHitAreas(mj.HitAreas)

	// Get the Drawables from core
	ds := c.core.GetDrawables(mocData.ModelPtr)
	drawables := make([]cubismmodel.DrawableInfo, len(ds))
	drawablesMap := make(map[string]cubismmodel.DrawableInfo, len(ds))
	for i, d := range ds {
		drawables[i] = cubismmodel.DrawableInfo{
			Id:              d.Id,
			Texture:         textures[d.Texture],
			VertexPositions: d.VertexPositions,
			VertexUvs:       d.VertexUvs,
			VertexIndices:   d.VertexIndices,
			ConstantFlag:    d.ConstantFlag,
			DynamicFlag:     d.DynamicFlag,
			Opacity:         d.Opacity,
			Masks:           d.Masks,
		}
		drawablesMap[d.Id] = drawables[i]
	}
	cubismModel.SetDrawables(drawables)
	cubismModel.SetDrawablesMap(drawablesMap)

	// Get the sorted indices
	cubismModel.SetSortedIndices(c.core.GetSortedDrawableIndices(mocData.ModelPtr))

	// Load the physics settings if they exist
	hasPhysics := false
	if mj.FileReferences.Physics != "" {
		physicsPath := filepath.Join(dir, mj.FileReferences.Physics)
		buf, err = os.ReadFile(physicsPath)
		if err != nil {
			return nil, err
		}
		var physicsJson model.PhysicsJson
		if err = json.Unmarshal(buf, &physicsJson); err != nil {
			return nil, err
		}
		cubismModel.SetPhysics(physicsJson)
		hasPhysics = true
	}

	// Load the pose settings if they exist
	hasPose := false
	if mj.FileReferences.Pose != "" {
		posePath := filepath.Join(dir, mj.FileReferences.Pose)
		buf, err = os.ReadFile(posePath)
		if err != nil {
			return nil, err
		}
		var poseJson model.PoseJson
		if err = json.Unmarshal(buf, &poseJson); err != nil {
			return nil, err
		}
		cubismModel.SetPose(poseJson)
		hasPose = true
	}

	// Load the display info settings if they exist
	if mj.FileReferences.DisplayInfo != "" {
		displayInfoPath := filepath.Join(dir, mj.FileReferences.DisplayInfo)
		buf, err = os.ReadFile(displayInfoPath)
		if err != nil {
			return nil, err
		}
		var cdi model.CdiJson
		if err = json.Unmarshal(buf, &cdi); err != nil {
			return nil, err
		}
		cubismModel.SetCdi(cdi)
	}

	// Load the expressions
	var exps []model.ExpJson
	for _, exp := range mj.FileReferences.Expressions {
		expPath := filepath.Join(dir, exp.File)
		buf, err = os.ReadFile(expPath)
		if err != nil {
			return nil, err
		}
		var e model.ExpJson
		if err = json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		e.Name = exp.Name
		exps = append(exps, e)
	}
	cubismModel.SetExpressions(exps)

	// Load the motion settings
	motions := map[string][]motion.Motion{}
	for name, motionList := range mj.FileReferences.Motions {
		motions[name] = []motion.Motion{}
		for _, mtn := range motionList {
			motionPath := filepath.Join(dir, mtn.File)
			buf, err = os.ReadFile(motionPath)
			if err != nil {
				return nil, err
			}
			var mtnJson model.MotionJson
			if err = json.Unmarshal(buf, &mtnJson); err != nil {
				return nil, err
			}
			fp := filepath.Base(mtn.File)
			mtn := mtnJson.ToMotion(fp, mtn.FadeInTime, mtn.FadeOutTime, mtn.Sound)
			if mtn.Sound != "" {
				soundPath := filepath.Join(dir, mtn.Sound)
				// If LoadSound is nil, don't play the sound
				if c.LoadSound == nil {
					mtn.LoadedSound, err = disabled.LoadSound(soundPath)
				} else {
					mtn.LoadedSound, err = c.LoadSound(soundPath)
				}
				if err != nil {
					// Sound loading failed (e.g., unsupported format) — use disabled sound instead of failing the entire model load
					mtn.LoadedSound, _ = disabled.LoadSound(soundPath)
					err = nil
				}
			}
			motions[name] = append(motions[name], mtn)
		}
	}

	// Load user data if it exists
	if mj.FileReferences.UserData != "" {
		userDataPath := filepath.Join(dir, mj.FileReferences.UserData)
		buf, err = os.ReadFile(userDataPath)
		if err != nil {
			return nil, err
		}
		var userdata model.UserDataJson
		if err = json.Unmarshal(buf, &userdata); err != nil {
			return nil, err
		}
		cubismModel.SetUserData(userdata)
	}

	// Create CubismUserModel (composition manager)
	userModel := cubismusermodel.NewCubismUserModel(cubismModel)

	// Set motions on the user model
	userModel.SetMotions(motions)

	// Initialize physics manager if physics data was loaded
	if hasPhysics {
		pm := physics.NewPhysicsManager(c.core, mocData.ModelPtr, cubismModel.Physics())
		if pm != nil {
			pm.SetIdManager(idManager)
			userModel.SetPhysicsManager(pm)
		}
	}

	// Initialize pose manager if pose data was loaded
	if hasPose {
		pm := pose.NewPoseManager(cubismModel.Pose(), c.core, mocData.ModelPtr)
		pm.SetIdManager(idManager)
		userModel.SetPoseManager(pm)
	}

	// Create the public Model (facade)
	m = &Model{
		inner: userModel,
	}

	return
}
