package core

import (
	"fmt"
	"strconv"
	"strings"

	core_5_0_0 "github.com/shaolei/cubism-go/internal/core/core_5_0_0"
	core_6_0_1 "github.com/shaolei/cubism-go/internal/core/core_6_0_1"
	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/minimum"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
)

type Core interface {
	LoadMoc(path string) (moc.Moc, error)
	GetVersion() string
	GetDynamicFlags(uintptr) []drawable.DynamicFlag
	GetOpacities(uintptr) []float32
	GetVertexPositions(uintptr) [][]drawable.Vector2
	GetDrawables(uintptr) []drawable.Drawable
	GetParameters(uintptr) []parameter.Parameter
	GetParameterValue(uintptr, string) float32
	SetParameterValue(uintptr, string, float32)
	GetPartIds(uintptr) []string
	SetPartOpacity(uintptr, string, float32)
	GetSortedDrawableIndices(uintptr) []int
	GetCanvasInfo(uintptr) (drawable.Vector2, drawable.Vector2, float32)
	Update(uintptr)
}

func NewCore(lib string) (c Core, err error) {
	l, err := openLibrary(lib)
	if err != nil {
		return
	}
	mc, err := minimum.NewCore(l)
	if err != nil {
		return
	}
	version := mc.GetVersion()

	// Parse major version to determine which core implementation to use
	major, err := parseMajorVersion(version)
	if err != nil {
		err = fmt.Errorf("failed to parse version %s: %w", version, err)
		return
	}

	switch major {
	case 5:
		c, err = core_5_0_0.NewCore(l)
	case 6:
		c, err = core_6_0_1.NewCore(l)
	default:
		err = fmt.Errorf("unsupported version: %s (major: %d)", version, major)
	}
	return
}

// parseMajorVersion extracts the major version number from a version string like "5.0.0" or "6.0.1"
func parseMajorVersion(version string) (int, error) {
	parts := strings.SplitN(version, ".", 2)
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format")
	}
	return strconv.Atoi(parts[0])
}
