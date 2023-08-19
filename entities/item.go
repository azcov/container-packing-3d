package entities

type Item struct {
	ID       int
	IsPacked bool
	Dim1     float64
	Dim2     float64
	Dim3     float64
	CoordX   float64
	CoordY   float64
	CoordZ   float64
	Quantity int
	PackDimX float64
	PackDimY float64
	PackDimZ float64
	Volume   float64
}

func NewItem(id int, dim1, dim2, dim3 float64, quantity int) *Item {
	volume := dim1 * dim2 * dim3
	return &Item{
		ID:       id,
		Dim1:     dim1,
		Dim2:     dim2,
		Dim3:     dim3,
		IsPacked: false,
		CoordX:   0,
		CoordY:   0,
		CoordZ:   0,
		Quantity: quantity,
		PackDimX: 0,
		PackDimY: 0,
		PackDimZ: 0,
		Volume:   volume,
	}
}
