package models

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

// ConstantGaussianProcessCovarianceKernel
type ConstantGaussianProcessCovarianceKernel struct {
	covMatrix  *mat.SymDense
	stateWidth int
}

func (c *ConstantGaussianProcessCovarianceKernel) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	c.stateWidth = settings.StateWidths[partitionIndex]
	c.covMatrix = mat.NewSymDense(c.stateWidth, nil)
	c.SetParams(settings.OtherParams[partitionIndex])
}

func (c *ConstantGaussianProcessCovarianceKernel) SetParams(
	params *simulator.OtherParams,
) {
	row := 0
	col := 0
	for _, param := range params.FloatParams["upper_triangle_covariance_matrix"] {
		c.covMatrix.SetSym(row, col, param)
		col += 1
		if col == c.stateWidth {
			row += 1
			col = row
		}
	}
}

func (c *ConstantGaussianProcessCovarianceKernel) GetCovariance(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) *mat.SymDense {
	return c.covMatrix
}
