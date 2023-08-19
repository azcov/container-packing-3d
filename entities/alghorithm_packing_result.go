package entities

type AlgorithmPackingResult struct {
	AlgorithmID                  int
	AlgorithmName                string
	IsCompletePack               bool
	PackedItems                  []*Item
	PackTimeInMilliseconds       int64
	PercentContainerVolumePacked float64
	PercentItemVolumePacked      float64
	UnpackedItems                []*Item
}

func NewAlgorithmPackingResult() AlgorithmPackingResult {
	return AlgorithmPackingResult{
		PackedItems:   make([]*Item, 0),
		UnpackedItems: make([]*Item, 0),
	}
}
