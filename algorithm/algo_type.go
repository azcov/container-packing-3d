package algo

type AlgorithmType int

const (
	EB_AFIT_ID AlgorithmType = iota
)

func GetPackingAlgorithmFromTypeID(algorithmTypeID int) IPackingAlgorithm {
	switch algorithmTypeID {
	case 0: // AlgorithmType.EB_AFIT
		return &EB_AFIT{}
	default:
		panic("Invalid algorithm type.")
	}
}
