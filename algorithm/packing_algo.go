package algo

import "github.com/azcov/continer-packing-3d/model"

type IPackingAlgorithm interface {
	Run(container model.Container, items []*model.Item) *model.AlgorithmPackingResult
}
