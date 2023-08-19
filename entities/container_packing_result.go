package entities

type ContainerPackingResult struct {
	ContainerID             int
	AlgorithmPackingResults []*AlgorithmPackingResult
}

func NewContainerPackingResult() ContainerPackingResult {
	return ContainerPackingResult{
		AlgorithmPackingResults: make([]*AlgorithmPackingResult, 0),
	}
}
