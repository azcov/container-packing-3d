package continer_packging_3d

import (
	"math"
	"sync"
	"time"

	algo "github.com/azcov/continer-packing-3d/algorithm"
	"github.com/azcov/continer-packing-3d/entities"
)

func Pack(containers []entities.Container, itemsToPack []*entities.Item, algorithmTypeIDs []int) []entities.ContainerPackingResult {
	var result []entities.ContainerPackingResult
	var sync sync.Mutex

	for _, container := range containers {
		containerPackingResult := entities.ContainerPackingResult{ContainerID: container.ID}

		for _, algorithmTypeID := range algorithmTypeIDs {
			algorithm := algo.GetPackingAlgorithmFromTypeID(algorithmTypeID)

			// Clone the item list to avoid interference with parallel updates
			items := make([]*entities.Item, len(itemsToPack))
			copy(items, itemsToPack)

			start := time.Now()
			algorithmResult := algorithm.Run(container, items)
			end := time.Now()

			algorithmResult.PackTimeInMilliseconds = end.Sub(start).Milliseconds()

			containerVolume := container.Length * container.Width * container.Height
			var itemVolumePacked, itemVolumeUnpacked float64
			for _, item := range algorithmResult.PackedItems {
				itemVolumePacked += item.Dim1 * item.Dim2 * item.Dim3 * float64(item.Quantity)
			}
			for _, item := range algorithmResult.UnpackedItems {
				itemVolumeUnpacked += item.Dim1 * item.Dim2 * item.Dim3 * float64(item.Quantity)
			}

			algorithmResult.PercentContainerVolumePacked = math.Round(itemVolumePacked/containerVolume*100*100) / 100
			algorithmResult.PercentItemVolumePacked = math.Round(itemVolumePacked/(itemVolumePacked+itemVolumeUnpacked)*100*100) / 100

			sync.Lock()
			containerPackingResult.AlgorithmPackingResults = append(containerPackingResult.AlgorithmPackingResults, algorithmResult)
			sync.Unlock()
		}

		// Sort algorithm results by name
		// Note: You need to implement the sorting logic based on AlgorithmName
		//       This part is missing in the provided code.

		sync.Lock()
		result = append(result, containerPackingResult)
		sync.Unlock()
	}

	return result
}
