package gausops

import (
	"github.com/EliCDavis/polyform/generator"
	"github.com/EliCDavis/polyform/refutil"
)

func init() {
	factory := &refutil.TypeFactory{}

	refutil.RegisterType[ColorGradingLutNode](factory)
	refutil.RegisterType[ScaleNode](factory)

	generator.RegisterTypes(factory)
}
