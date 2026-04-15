package core

import (
	"github.com/shaolei/cubism-go/internal/core/base"
	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/ebitengine/purego"
)

// Core implements the Cubism Core API for version 5.x.
type Core struct {
	funcs base.Funcs
}

func NewCore(lib uintptr) (c *Core, err error) {
	c = new(Core)
	base.RegisterCommonFuncs(&c.funcs, lib)
	// v5 uses csmGetDrawableRenderOrders for sorting
	purego.RegisterLibFunc(&c.funcs.CsmGetDrawableSortOrders, lib, "csmGetDrawableRenderOrders")
	return
}

func (c *Core) LoadMoc(path string) (moc.Moc, error) {
	return base.LoadMoc(&c.funcs, path)
}

func (c *Core) GetVersion() string {
	return base.GetVersion(&c.funcs)
}

func (c *Core) GetDynamicFlags(modelPtr uintptr) []drawable.DynamicFlag {
	return base.GetDynamicFlags(&c.funcs, modelPtr)
}

func (c *Core) GetOpacities(modelPtr uintptr) []float32 {
	return base.GetOpacities(&c.funcs, modelPtr)
}

func (c *Core) GetVertexPositions(modelPtr uintptr) [][]drawable.Vector2 {
	return base.GetVertexPositions(&c.funcs, modelPtr)
}

func (c *Core) GetDrawables(modelPtr uintptr) []drawable.Drawable {
	return base.GetDrawables(&c.funcs, modelPtr)
}

func (c *Core) GetParameters(modelPtr uintptr) []parameter.Parameter {
	return base.GetParameters(&c.funcs, modelPtr)
}

func (c *Core) GetParameterIds(modelPtr uintptr) []string {
	return base.GetParameterIds(&c.funcs, modelPtr)
}

func (c *Core) GetParameterValue(modelPtr uintptr, id string) float32 {
	return base.GetParameterValue(&c.funcs, modelPtr, id)
}

func (c *Core) SetParameterValue(modelPtr uintptr, id string, value float32) {
	base.SetParameterValue(&c.funcs, modelPtr, id, value)
}

func (c *Core) GetParameterValueByIndex(modelPtr uintptr, index int) float32 {
	return base.GetParameterValueByIndex(&c.funcs, modelPtr, index)
}

func (c *Core) SetParameterValueByIndex(modelPtr uintptr, index int, value float32) {
	base.SetParameterValueByIndex(&c.funcs, modelPtr, index, value)
}

func (c *Core) GetPartIds(modelPtr uintptr) []string {
	return base.GetPartIds(&c.funcs, modelPtr)
}

func (c *Core) GetPartOpacities(modelPtr uintptr) []float32 {
	return base.GetPartOpacities(&c.funcs, modelPtr)
}

func (c *Core) SetPartOpacity(modelPtr uintptr, id string, value float32) {
	base.SetPartOpacity(&c.funcs, modelPtr, id, value)
}

func (c *Core) SetPartOpacityByIndex(modelPtr uintptr, index int, value float32) {
	base.SetPartOpacityByIndex(&c.funcs, modelPtr, index, value)
}

func (c *Core) GetPartOpacityByIndex(modelPtr uintptr, index int) float32 {
	return base.GetPartOpacityByIndex(&c.funcs, modelPtr, index)
}

func (c *Core) GetParameterCount(modelPtr uintptr) int {
	return base.GetParameterCount(&c.funcs, modelPtr)
}

func (c *Core) GetParameterValues(modelPtr uintptr) []float32 {
	return base.GetParameterValues(&c.funcs, modelPtr)
}

func (c *Core) GetParameterMinimumValues(modelPtr uintptr) []float32 {
	return base.GetParameterMinimumValues(&c.funcs, modelPtr)
}

func (c *Core) GetParameterMaximumValues(modelPtr uintptr) []float32 {
	return base.GetParameterMaximumValues(&c.funcs, modelPtr)
}

func (c *Core) GetParameterDefaultValues(modelPtr uintptr) []float32 {
	return base.GetParameterDefaultValues(&c.funcs, modelPtr)
}

func (c *Core) GetSortedDrawableIndices(modelPtr uintptr) []int {
	return base.GetSortedDrawableIndices(&c.funcs, modelPtr)
}

func (c *Core) GetCanvasInfo(modelPtr uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return base.GetCanvasInfo(&c.funcs, modelPtr)
}

func (c *Core) Update(modelPtr uintptr) {
	base.Update(&c.funcs, modelPtr)
}
