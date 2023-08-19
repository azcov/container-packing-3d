package algo

import "github.com/azcov/continer-packing-3d/entities"

type AlgorithmBase interface {
	Run(container entities.Container, items []*entities.Item) *entities.AlgorithmPackingResult
}
