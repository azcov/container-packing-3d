package model

type Orientation int

const (
	OrientationLHW Orientation = iota
	OrientationWHL
	OrientationWLH
	OrientationHLW
	OrientationLWH
	OrientationHWL
)

func DimensionByOrientation(n Orientation, px, py, pz float64) (float64, float64, float64) {
	switch n {
	case OrientationLHW:
		return px, py, pz
	case OrientationWHL:
		return pz, py, px
	case OrientationWLH:
		return pz, px, py
	case OrientationHLW:
		return py, px, pz
	case OrientationLWH:
		return px, pz, py
	case OrientationHWL:
		return py, pz, px
	}
	return px, py, pz
}
