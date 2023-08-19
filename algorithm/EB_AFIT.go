package algo

import (
	"log"
	"math"
	"sort"

	"github.com/azcov/continer-packing-3d/entities"
)

type Layer struct {
	LayerDim  float64
	LayerEval float64
}

type ScrapPad struct {
	CumX float64
	CumZ float64
	Post *ScrapPad
	Pre  *ScrapPad
}

// type AlgorithmPackingResult struct {
// 	AlgorithmID    int
// 	AlgorithmName  string
// 	UnpackedItems  []*Item
// 	PackedItems    []*Item
// 	IsCompletePack bool
// }

type EB_AFIT struct {
	itemsToPack        []*entities.Item
	itemsPackedInOrder []*entities.Item
	layers             []*Layer
	result             *entities.ContainerPackingResult

	scrapfirst *ScrapPad
	smallestZ  *ScrapPad
	trash      *ScrapPad

	evened               bool
	hundredPercentPacked bool
	layerDone            bool
	packing              bool
	packingBest          bool
	quit                 bool

	bboxi           int
	bestIteration   int
	bestVariant     int
	boxi            int
	cboxi           int
	layerListLen    int
	packedItemCount int
	x               int

	itemsToPackCount int
	// bbfx                 decimal.Decimal
	// bbfy                 decimal.Decimal
	// bbfz                 decimal.Decimal
	// bboxx                decimal.Decimal
	// bboxy                decimal.Decimal
	// bboxz                decimal.Decimal
	// bfx                  decimal.Decimal
	// bfy                  decimal.Decimal
	// bfz                  decimal.Decimal
	// boxx                 decimal.Decimal
	// boxy                 decimal.Decimal
	// boxz                 decimal.Decimal
	// cboxx                decimal.Decimal
	// cboxy                decimal.Decimal
	// cboxz                decimal.Decimal
	// layerinlayer         decimal.Decimal
	// layerThickness       decimal.Decimal
	// lilz                 decimal.Decimal
	// packedVolume         decimal.Decimal
	// packedy              decimal.Decimal
	// prelayer             decimal.Decimal
	// prepackedy           decimal.Decimal
	// preremainpy          decimal.Decimal
	// px                   decimal.Decimal
	// py                   decimal.Decimal
	// pz                   decimal.Decimal
	// remainpy             decimal.Decimal
	// remainpz             decimal.Decimal
	// totalItemVolume      decimal.Decimal
	// totalContainerVolume decimal.Decimal

	bbfx                 float64
	bbfy                 float64
	bbfz                 float64
	bboxx                float64
	bboxy                float64
	bboxz                float64
	bfx                  float64
	bfy                  float64
	bfz                  float64
	boxx                 float64
	boxy                 float64
	boxz                 float64
	cboxx                float64
	cboxy                float64
	cboxz                float64
	layerinlayer         float64
	layerThickness       float64
	lilz                 float64
	packedVolume         float64
	packedy              float64
	prelayer             float64
	prepackedy           float64
	preremainpy          float64
	px                   float64
	py                   float64
	pz                   float64
	remainpy             float64
	remainpz             float64
	totalItemVolume      float64
	totalContainerVolume float64
}

func (e *EB_AFIT) Run(container entities.Container, items []*entities.Item) *entities.AlgorithmPackingResult {
	e.Initialize(container, items)
	e.ExecuteIterations(container)
	e.Report(container)

	result := entities.AlgorithmPackingResult{}
	result.AlgorithmID = int(EB_AFIT_ID)
	result.AlgorithmName = "EB-AFIT"
	result.UnpackedItems = make([]*entities.Item, 0)

	for i := 1; i <= e.itemsToPackCount; i++ {
		e.itemsToPack[i].Quantity = 1

		if !e.itemsToPack[i].IsPacked {
			result.UnpackedItems = append(result.UnpackedItems, e.itemsToPack[i])
		}
	}

	result.PackedItems = e.itemsPackedInOrder

	if len(result.UnpackedItems) == 0 {
		result.IsCompletePack = true
	}

	return &result
}
func (e *EB_AFIT) AnalyzeBox(hmx, hy, hmy, hz, hmz, dim1, dim2, dim3 float64) {
	if dim1 <= hmx && dim2 <= hmy && dim3 <= hmz {
		if dim2 <= hy {
			if hy-dim2 < e.bfy {
				e.boxx = dim1
				e.boxy = dim2
				e.boxz = dim3
				e.bfx = hmx - dim1
				e.bfy = hy - dim2
				e.bfz = math.Abs(hz - dim3)
				e.boxi = e.x
			} else if hy-dim2 == e.bfy && hmx-dim1 < e.bfx {
				e.boxx = dim1
				e.boxy = dim2
				e.boxz = dim3
				e.bfx = hmx - dim1
				e.bfy = hy - dim2
				e.bfz = math.Abs(hz - dim3)
				e.boxi = e.x
			} else if hy-dim2 == e.bfy && hmx-dim1 == e.bfx && math.Abs(hz-dim3) < e.bfz {
				e.boxx = dim1
				e.boxy = dim2
				e.boxz = dim3
				e.bfx = hmx - dim1
				e.bfy = hy - dim2
				e.bfz = math.Abs(hz - dim3)
				e.boxi = e.x
			}
		} else {
			if dim2-hy < e.bbfy {
				e.bboxx = dim1
				e.bboxy = dim2
				e.bboxz = dim3
				e.bbfx = hmx - dim1
				e.bbfy = dim2 - hy
				e.bbfz = math.Abs(hz - dim3)
				e.bboxi = e.x
			} else if dim2-hy == e.bbfy && hmx-dim1 < e.bbfx {
				e.bboxx = dim1
				e.bboxy = dim2
				e.bboxz = dim3
				e.bbfx = hmx - dim1
				e.bbfy = dim2 - hy
				e.bbfz = math.Abs(hz - dim3)
				e.bboxi = e.x
			} else if dim2-hy == e.bbfy && hmx-dim1 == e.bbfx && math.Abs(hz-dim3) < e.bbfz {
				e.bboxx = dim1
				e.bboxy = dim2
				e.bboxz = dim3
				e.bbfx = hmx - dim1
				e.bbfy = dim2 - hy
				e.bbfz = math.Abs(hz - dim3)
				e.bboxi = e.x
			}
		}
	}
}

// / <summary>
// / After finding each box, the candidate boxes and the condition of the layer are examined.
// / </summary>
func (e *EB_AFIT) CheckFound() {
	e.evened = false

	if e.boxi != 0 {
		e.cboxi = e.boxi
		e.cboxx = e.boxx
		e.cboxy = e.boxy
		e.cboxz = e.boxz
	} else {
		if e.bboxi > 0 && (e.layerinlayer != 0 || (e.smallestZ.Pre == nil && e.smallestZ.Post == nil)) {
			if e.layerinlayer == 0 {
				e.prelayer = e.layerThickness
				e.lilz = e.smallestZ.CumZ
			}

			e.cboxi = e.bboxi
			e.cboxx = e.bboxx
			e.cboxy = e.bboxy
			e.cboxz = e.bboxz
			e.layerinlayer = e.layerinlayer + e.bboxy - e.layerThickness
			e.layerThickness = e.bboxy
		} else {
			if e.smallestZ.Pre == nil && e.smallestZ.Post == nil {
				e.layerDone = true
			} else {
				e.evened = true

				if e.smallestZ.Pre == nil {
					e.trash = e.smallestZ.Post
					e.smallestZ.CumX = e.smallestZ.Post.CumX
					e.smallestZ.CumZ = e.smallestZ.Post.CumZ
					e.smallestZ.Post = e.smallestZ.Post.Post
					if e.smallestZ.Post != nil {
						e.smallestZ.Post.Pre = e.smallestZ
					}
				} else if e.smallestZ.Post == nil {
					e.smallestZ.Pre.Post = nil
					e.smallestZ.Pre.CumX = e.smallestZ.CumX
				} else {
					if e.smallestZ.Pre.CumZ == e.smallestZ.Post.CumZ {
						e.smallestZ.Pre.Post = e.smallestZ.Post.Post

						if e.smallestZ.Post.Post != nil {
							e.smallestZ.Post.Post.Pre = e.smallestZ.Pre
						}

						e.smallestZ.Pre.CumX = e.smallestZ.Post.CumX
					} else {
						e.smallestZ.Pre.Post = e.smallestZ.Post
						e.smallestZ.Post.Pre = e.smallestZ.Pre

						if e.smallestZ.Pre.CumZ < e.smallestZ.Post.CumZ {
							e.smallestZ.Pre.CumX = e.smallestZ.CumX
						}
					}
				}
			}
		}
	}
}

// / <summary>
// / Executes the packing algorithm variants.
// / </summary>
func (e *EB_AFIT) ExecuteIterations(container entities.Container) {
	var itelayer int
	var layersIndex int
	bestVolume := 0.0
	quit := false

	for containerOrientationVariant := 1; containerOrientationVariant <= 6 && !quit; containerOrientationVariant++ {
		switch containerOrientationVariant {
		case 1:
			e.px = container.Length
			e.py = container.Height
			e.pz = container.Width
		case 2:
			e.px = container.Width
			e.py = container.Height
			e.pz = container.Length
		case 3:
			e.px = container.Width
			e.py = container.Length
			e.pz = container.Height
		case 4:
			e.px = container.Height
			e.py = container.Length
			e.pz = container.Width
		case 5:
			e.px = container.Length
			e.py = container.Width
			e.pz = container.Height
		case 6:
			e.px = container.Height
			e.py = container.Width
			e.pz = container.Length
		}

		e.layers = append(e.layers, &Layer{LayerEval: -1})
		e.ListCanditLayers()
		sort.Slice(e.layers, func(i, j int) bool {
			return e.layers[i].LayerEval < e.layers[j].LayerEval
		})

		for layersIndex = 1; layersIndex <= e.layerListLen && !quit; layersIndex++ {
			e.packedVolume = 0.0
			e.packedy = 0
			e.packing = true
			e.layerThickness = e.layers[layersIndex].LayerDim
			itelayer = layersIndex
			e.remainpy = e.py
			e.remainpz = e.pz
			e.packedItemCount = 0

			for x := 1; x <= e.itemsToPackCount; x++ {
				e.itemsToPack[x].IsPacked = false
			}

			for {
				e.layerinlayer = 0
				e.layerDone = false
				e.PackLayer()
				e.packedy = e.packedy + e.layerThickness
				e.remainpy = e.py - e.packedy

				if e.layerinlayer != 0 && !quit {
					e.prepackedy = e.packedy
					e.preremainpy = e.remainpy
					e.remainpy = e.layerThickness - e.prelayer
					e.packedy = e.packedy - e.layerThickness + e.prelayer
					e.remainpz = e.lilz
					e.layerThickness = e.layerinlayer
					e.layerDone = false
					e.PackLayer()
					e.packedy = e.prepackedy
					e.remainpy = e.preremainpy
					e.remainpz = e.pz
				}

				e.FindLayer(e.remainpy)
				if !e.packing || quit {
					break
				}
			}

			if e.packedVolume > bestVolume && !quit {
				bestVolume = e.packedVolume
				e.bestVariant = containerOrientationVariant
				e.bestIteration = itelayer
			}

			if e.hundredPercentPacked {
				break
			}
		}

		if e.hundredPercentPacked {
			break
		}

		if container.Length == container.Height && container.Height == container.Width {
			containerOrientationVariant = 6
		}

		e.layers = nil
	}
}

// / <summary>
// / Finds the most proper boxes by looking at all six possible orientations,
// / empty space given, adjacent boxes, and pallet limits.
// / </summary>
func (e *EB_AFIT) FindBox(hmx, hy, hmy, hz, hmz float64) {

	// var y int
	e.bfx = 32767
	e.bfy = 32767
	e.bfz = 32767
	e.bbfx = 32767
	e.bbfy = 32767
	e.bbfz = 32767
	e.boxi = 0
	e.bboxi = 0

	for y := 0; y < e.itemsToPackCount; y += e.itemsToPack[y].Quantity {
		x := y
		for x = y; x < y+e.itemsToPack[y].Quantity-1; x++ {
			if !e.itemsToPack[x].IsPacked {
				break
			}
		}

		if e.itemsToPack[x].IsPacked {
			continue
		}

		if x > e.itemsToPackCount {
			return
		}

		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim1, e.itemsToPack[x].Dim2, e.itemsToPack[x].Dim3)

		if e.itemsToPack[x].Dim1 == e.itemsToPack[x].Dim3 && e.itemsToPack[x].Dim3 == e.itemsToPack[x].Dim2 {
			continue
		}

		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim1, e.itemsToPack[x].Dim3, e.itemsToPack[x].Dim2)
		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim2, e.itemsToPack[x].Dim1, e.itemsToPack[x].Dim3)
		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim2, e.itemsToPack[x].Dim3, e.itemsToPack[x].Dim1)
		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim3, e.itemsToPack[x].Dim1, e.itemsToPack[x].Dim2)
		e.AnalyzeBox(hmx, hy, hmy, hz, hmz, e.itemsToPack[x].Dim3, e.itemsToPack[x].Dim2, e.itemsToPack[x].Dim1)
	}
}

// / <summary>
// / Finds the most proper layer height by looking at the unpacked boxes and the remaining empty space available.
// / </summary>

func (e *EB_AFIT) FindLayer(thickness float64) {
	layerThickness := 0.0
	eval := 1000000.0

	for x, item := range e.itemsToPack {
		if item.IsPacked {
			continue
		}

		for y := 1; y <= 3; y++ {
			var exdim, dimen2, dimen3 float64
			switch y {
			case 1:
				exdim = item.Dim1
				dimen2 = item.Dim2
				dimen3 = item.Dim3
			case 2:
				exdim = item.Dim2
				dimen2 = item.Dim1
				dimen3 = item.Dim3
			case 3:
				exdim = item.Dim3
				dimen2 = item.Dim1
				dimen3 = item.Dim2
			}

			layereval := 0.0

			if exdim <= thickness && ((dimen2 <= e.px && dimen3 <= e.pz) || (dimen3 <= e.px && dimen2 <= e.pz)) {
				for z, otherItem := range e.itemsToPack {
					if x != z && !otherItem.IsPacked {
						dimdif := math.Abs(exdim - otherItem.Dim1)

						if math.Abs(exdim-otherItem.Dim2) < dimdif {
							dimdif = math.Abs(exdim - otherItem.Dim2)
						}

						if math.Abs(exdim-otherItem.Dim3) < dimdif {
							dimdif = math.Abs(exdim - otherItem.Dim3)
						}

						layereval += dimdif
					}
				}

				if layereval < eval {
					eval = layereval
					layerThickness = exdim
				}
			}
		}
	}

	if layerThickness == 0 || layerThickness > e.remainpy {
		e.packing = false
	}
}

// / <summary>
// / Finds the first to be packed gap in the layer edge.
// / </summary>
func (e *EB_AFIT) FindSmallestZ() {

	scrapmemb := e.scrapfirst
	smallestZ := scrapmemb

	for scrapmemb.Post != nil {
		if scrapmemb.Post.CumZ < smallestZ.CumZ {
			smallestZ = scrapmemb.Post
		}

		scrapmemb = scrapmemb.Post
	}
}

// / <summary>
// / Initializes everything.
// / </summary>
func (e *EB_AFIT) Initialize(container entities.Container, items []*entities.Item) {
	e.itemsToPack = []*entities.Item{}
	e.itemsPackedInOrder = []*entities.Item{}
	e.result = &entities.ContainerPackingResult{}

	// Add a fake entry at the beginning for 1-based indexing
	e.itemsToPack = append(e.itemsToPack, &entities.Item{})

	e.layers = make([]*Layer, 0)
	e.itemsToPackCount = 0

	for _, item := range items {
		for i := 1; i <= item.Quantity; i++ {
			newItem := &entities.Item{
				ID:       item.ID,
				Dim1:     item.Dim1,
				Dim2:     item.Dim2,
				Dim3:     item.Dim3,
				Quantity: item.Quantity,
				Volume:   item.Volume,
			}
			e.itemsToPack = append(e.itemsToPack, newItem)
		}

		e.itemsToPackCount += item.Quantity
	}

	// Add another fake entry at the end for 1-based indexing
	e.itemsToPack = append(e.itemsToPack, &entities.Item{})

	e.totalContainerVolume = container.Length * container.Height * container.Width
	totalItemVolume := 0.0

	for x := 1; x <= e.itemsToPackCount; x++ {
		totalItemVolume += e.itemsToPack[x].Volume
	}

	scrapfirst := &ScrapPad{}
	scrapfirst.Pre = nil
	scrapfirst.Post = nil
	e.packingBest = false
	e.hundredPercentPacked = false
	e.quit = false

	// Rest of your initialization code...
}

// / <summary>
// / Lists all possible layer heights by giving a weight value to each of them.
// / </summary>
func (e *EB_AFIT) ListCanditLayers() {
	var same bool
	var exdim, dimdif, dimen2, dimen3 float64
	var y, z, k int
	var layereval float64

	layerListLen := 0
	layers := []Layer{}

	for x := 1; x <= e.itemsToPackCount; x++ {
		for y = 1; y <= 3; y++ {
			switch y {
			case 1:
				exdim = e.itemsToPack[x].Dim1
				dimen2 = e.itemsToPack[x].Dim2
				dimen3 = e.itemsToPack[x].Dim3
			case 2:
				exdim = e.itemsToPack[x].Dim2
				dimen2 = e.itemsToPack[x].Dim1
				dimen3 = e.itemsToPack[x].Dim3
			case 3:
				exdim = e.itemsToPack[x].Dim3
				dimen2 = e.itemsToPack[x].Dim1
				dimen3 = e.itemsToPack[x].Dim2
			}

			if exdim > e.py || ((dimen2 > e.px || dimen3 > e.pz) && (dimen3 > e.px || dimen2 > e.pz)) {
				continue
			}

			same = false

			for k = 1; k <= layerListLen; k++ {
				if exdim == layers[k].LayerDim {
					same = true
					continue
				}
			}

			if same {
				continue
			}

			layereval = 0

			for z = 1; z <= e.itemsToPackCount; z++ {
				if x != z {
					dimdif = math.Abs(exdim - e.itemsToPack[z].Dim1)

					if math.Abs(exdim-e.itemsToPack[z].Dim2) < dimdif {
						dimdif = math.Abs(exdim - e.itemsToPack[z].Dim2)
					}
					if math.Abs(exdim-e.itemsToPack[z].Dim3) < dimdif {
						dimdif = math.Abs(exdim - e.itemsToPack[z].Dim3)
					}
					layereval += dimdif
				}
			}

			layerListLen++
			layers = append(layers, Layer{LayerDim: exdim, LayerEval: layereval})
		}
	}

	// Rest of your method...
}

/// <summary>
/// Transforms the found coordinate system to the one entered by the user and writes them
/// to the report file.
/// </summary>

func (e *EB_AFIT) OutputBoxList() {
	var packCoordX, packCoordY, packCoordZ float64
	var packDimX, packDimY, packDimZ float64

	switch e.bestVariant {
	case 1:
		packCoordX = e.itemsToPack[e.cboxi].CoordX
		packCoordY = e.itemsToPack[e.cboxi].CoordY
		packCoordZ = e.itemsToPack[e.cboxi].CoordZ
		packDimX = e.itemsToPack[e.cboxi].PackDimX
		packDimY = e.itemsToPack[e.cboxi].PackDimY
		packDimZ = e.itemsToPack[e.cboxi].PackDimZ
	case 2:
		packCoordX = e.itemsToPack[e.cboxi].CoordZ
		packCoordY = e.itemsToPack[e.cboxi].CoordY
		packCoordZ = e.itemsToPack[e.cboxi].CoordX
		packDimX = e.itemsToPack[e.cboxi].PackDimZ
		packDimY = e.itemsToPack[e.cboxi].PackDimY
		packDimZ = e.itemsToPack[e.cboxi].PackDimX
	case 3:
		packCoordX = e.itemsToPack[e.cboxi].CoordY
		packCoordY = e.itemsToPack[e.cboxi].CoordZ
		packCoordZ = e.itemsToPack[e.cboxi].CoordX
		packDimX = e.itemsToPack[e.cboxi].PackDimY
		packDimY = e.itemsToPack[e.cboxi].PackDimZ
		packDimZ = e.itemsToPack[e.cboxi].PackDimX
	case 4:
		packCoordX = e.itemsToPack[e.cboxi].CoordY
		packCoordY = e.itemsToPack[e.cboxi].CoordX
		packCoordZ = e.itemsToPack[e.cboxi].CoordZ
		packDimX = e.itemsToPack[e.cboxi].PackDimY
		packDimY = e.itemsToPack[e.cboxi].PackDimX
		packDimZ = e.itemsToPack[e.cboxi].PackDimZ
	case 5:
		packCoordX = e.itemsToPack[e.cboxi].CoordX
		packCoordY = e.itemsToPack[e.cboxi].CoordZ
		packCoordZ = e.itemsToPack[e.cboxi].CoordY
		packDimX = e.itemsToPack[e.cboxi].PackDimX
		packDimY = e.itemsToPack[e.cboxi].PackDimZ
		packDimZ = e.itemsToPack[e.cboxi].PackDimY
	case 6:
		packCoordX = e.itemsToPack[e.cboxi].CoordZ
		packCoordY = e.itemsToPack[e.cboxi].CoordX
		packCoordZ = e.itemsToPack[e.cboxi].CoordY
		packDimX = e.itemsToPack[e.cboxi].PackDimZ
		packDimY = e.itemsToPack[e.cboxi].PackDimX
		packDimZ = e.itemsToPack[e.cboxi].PackDimY
	}

	e.itemsToPack[e.cboxi].CoordX = packCoordX
	e.itemsToPack[e.cboxi].CoordY = packCoordY
	e.itemsToPack[e.cboxi].CoordZ = packCoordZ
	e.itemsToPack[e.cboxi].PackDimX = packDimX
	e.itemsToPack[e.cboxi].PackDimY = packDimY
	e.itemsToPack[e.cboxi].PackDimZ = packDimZ

	e.itemsPackedInOrder = append(e.itemsPackedInOrder, e.itemsToPack[e.cboxi])
}

/// <summary>
/// Packs the boxes found and arranges all variables and records properly.
/// </summary>

func (e *EB_AFIT) PackLayer() {
	var lenx, lenz, lpz float64

	if e.layerThickness == 0 {
		e.packing = false
		return
	}

	e.scrapfirst.CumX = e.px
	e.scrapfirst.CumZ = 0

	for !e.quit {
		e.FindSmallestZ()

		if e.smallestZ.Pre == nil && e.smallestZ.Post == nil {
			// *** SITUATION-1: NO BOXES ON THE RIGHT AND LEFT SIDES ***

			lenx = e.smallestZ.CumX
			lpz = float64(e.remainpz - e.smallestZ.CumZ)
			e.FindBox(lenx, e.layerThickness, e.remainpy, lpz, lpz)
			e.CheckFound()

			if e.layerDone {
				break
			}
			if e.evened {
				continue
			}

			e.itemsToPack[e.cboxi].CoordX = 0
			e.itemsToPack[e.cboxi].CoordY = e.packedy
			e.itemsToPack[e.cboxi].CoordZ = e.smallestZ.CumZ
			if e.cboxx == e.smallestZ.CumX {
				e.smallestZ.CumZ += e.cboxz
			} else {
				e.smallestZ.Post = &ScrapPad{}

				e.smallestZ.Post.Post = nil
				e.smallestZ.Post.Pre = e.smallestZ
				e.smallestZ.Post.CumX = e.smallestZ.CumX
				e.smallestZ.Post.CumZ = e.smallestZ.CumZ
				e.smallestZ.CumX = e.cboxx
				e.smallestZ.CumZ += e.cboxz
			}
		} else if e.smallestZ.Pre == nil {
			// *** SITUATION-2: NO BOXES ON THE LEFT SIDE ***

			lenx = e.smallestZ.CumX
			lenz = e.smallestZ.Post.CumZ - e.smallestZ.CumZ
			lpz = float64(e.remainpz - e.smallestZ.CumZ)
			e.FindBox(lenx, e.layerThickness, e.remainpy, lenz, lpz)
			e.CheckFound()

			if e.layerDone {
				break
			}
			if e.evened {
				continue
			}

			e.itemsToPack[e.cboxi].CoordY = e.packedy
			e.itemsToPack[e.cboxi].CoordZ = e.smallestZ.CumZ
			if e.cboxx == e.smallestZ.CumX {
				e.itemsToPack[e.cboxi].CoordX = 0

				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Post.CumZ {
					e.smallestZ.CumZ = e.smallestZ.Post.CumZ
					e.smallestZ.CumX = e.smallestZ.Post.CumX
					e.trash = e.smallestZ.Post
					e.smallestZ.Post = e.smallestZ.Post.Post

					if e.smallestZ.Post != nil {
						e.smallestZ.Post.Pre = e.smallestZ
					}
				} else {
					e.smallestZ.CumZ += e.cboxz
				}
			} else {
				e.itemsToPack[e.cboxi].CoordX = e.smallestZ.CumX - e.cboxx

				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Post.CumZ {
					e.smallestZ.CumX = e.smallestZ.CumX - e.cboxx
				} else {
					e.smallestZ.Post.Pre = &ScrapPad{}

					e.smallestZ.Post.Pre.Post = e.smallestZ.Post
					e.smallestZ.Post.Pre.Pre = e.smallestZ
					e.smallestZ.Post = e.smallestZ.Post.Pre
					e.smallestZ.Post.CumX = e.smallestZ.CumX
					e.smallestZ.CumX -= e.cboxx
					e.smallestZ.Post.CumZ += e.cboxz
				}
			}
		} else if e.smallestZ.Post == nil {
			// *** SITUATION-3: NO BOXES ON THE RIGHT SIDE ***

			lenx = e.smallestZ.CumX - e.smallestZ.Pre.CumX
			lenz = e.smallestZ.Pre.CumZ - e.smallestZ.CumZ
			lpz = float64(e.remainpz - e.smallestZ.CumZ)
			e.FindBox(lenx, e.layerThickness, e.remainpy, lenz, lpz)
			e.CheckFound()

			if e.layerDone {
				break
			}
			if e.evened {
				continue
			}

			e.itemsToPack[e.cboxi].CoordY = e.packedy
			e.itemsToPack[e.cboxi].CoordZ = e.smallestZ.CumZ
			e.itemsToPack[e.cboxi].CoordX = e.smallestZ.Pre.CumX

			if e.cboxx == e.smallestZ.CumX-e.smallestZ.Pre.CumX {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.Pre.CumX = e.smallestZ.CumX
					e.smallestZ.Pre.Post = nil
				} else {
					e.smallestZ.CumZ += e.cboxz
				}
			} else {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.Pre.CumX += e.cboxx
				} else {
					e.smallestZ.Pre.Post = &ScrapPad{}

					e.smallestZ.Pre.Post.Pre = e.smallestZ.Pre
					e.smallestZ.Pre.Post.Post = e.smallestZ
					e.smallestZ.Pre = e.smallestZ.Pre.Post
					e.smallestZ.Pre.CumX = e.smallestZ.Pre.Pre.CumX + e.cboxx
					e.smallestZ.Pre.CumZ += e.cboxz
				}
			}
		} else if e.smallestZ.Pre.CumZ == e.smallestZ.Post.CumZ {
			// *** SITUATION-4: THERE ARE BOXES ON BOTH SIDES ***

			// *** SUBSITUATION-4A: SIDES ARE EQUAL TO EACH OTHER ***

			lenx = e.smallestZ.CumX - e.smallestZ.Pre.CumX
			lenz = e.smallestZ.Pre.CumZ - e.smallestZ.CumZ
			lpz = float64(e.remainpz - e.smallestZ.CumZ)
			e.FindBox(lenx, e.layerThickness, e.remainpy, lenz, lpz)
			e.CheckFound()

			if e.layerDone {
				break
			}
			if e.evened {
				continue
			}

			e.itemsToPack[e.cboxi].CoordY = e.packedy
			e.itemsToPack[e.cboxi].CoordZ = e.smallestZ.CumZ

			if e.cboxx == e.smallestZ.CumX-e.smallestZ.Pre.CumX {
				e.itemsToPack[e.cboxi].CoordX = e.smallestZ.Pre.CumX

				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Post.CumZ {
					e.smallestZ.Pre.CumX = e.smallestZ.Post.CumX

					if e.smallestZ.Post.Post != nil {
						e.smallestZ.Pre.Post = e.smallestZ.Post.Post
						e.smallestZ.Post.Post.Pre = e.smallestZ.Pre
					} else {
						e.smallestZ.Pre.Post = nil
					}
				} else {
					e.smallestZ.CumZ += e.cboxz
				}
			} else if e.smallestZ.Pre.CumX < e.px-e.smallestZ.CumX {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.CumX -= e.cboxx
					e.itemsToPack[e.cboxi].CoordX = e.smallestZ.CumX
				} else {
					e.itemsToPack[e.cboxi].CoordX = e.smallestZ.Pre.CumX
					e.smallestZ.Pre.Post = &ScrapPad{}

					e.smallestZ.Pre.Post.Pre = e.smallestZ.Pre
					e.smallestZ.Pre.Post.Post = e.smallestZ
					e.smallestZ.Pre = e.smallestZ.Pre.Post
					e.smallestZ.Pre.CumX = e.smallestZ.Pre.Pre.CumX + e.cboxx
					e.smallestZ.Pre.CumZ += e.cboxz
				}
			} else {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.Pre.CumX += e.cboxx
					e.itemsToPack[e.cboxi].CoordX = e.smallestZ.Pre.CumX
				} else {
					e.itemsToPack[e.cboxi].CoordX = e.smallestZ.CumX - e.cboxx
					e.smallestZ.Post.Pre = &ScrapPad{}

					e.smallestZ.Post.Pre.Post = e.smallestZ.Post
					e.smallestZ.Post.Pre.Pre = e.smallestZ
					e.smallestZ.Post = e.smallestZ.Post.Pre
					e.smallestZ.Post.CumX = e.smallestZ.CumX
					e.smallestZ.Post.CumZ += e.cboxz
					e.smallestZ.CumX -= e.cboxx
				}
			}
		} else {
			// *** SUBSITUATION-4B: SIDES ARE NOT EQUAL TO EACH OTHER ***

			lenx = e.smallestZ.CumX - e.smallestZ.Pre.CumX
			lenz = e.smallestZ.Pre.CumZ - e.smallestZ.CumZ
			lpz = float64(e.remainpz - e.smallestZ.CumZ)
			e.FindBox(lenx, e.layerThickness, e.remainpy, lenz, lpz)
			e.CheckFound()

			if e.layerDone {
				break
			}
			if e.evened {
				continue
			}

			e.itemsToPack[e.cboxi].CoordY = e.packedy
			e.itemsToPack[e.cboxi].CoordZ = e.smallestZ.CumZ
			e.itemsToPack[e.cboxi].CoordX = e.smallestZ.Pre.CumX

			if e.cboxx == e.smallestZ.CumX-e.smallestZ.Pre.CumX {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.Pre.CumX = e.smallestZ.CumX
					e.smallestZ.Pre.Post = e.smallestZ.Post
					e.smallestZ.Post.Pre = e.smallestZ.Pre
				} else {
					e.smallestZ.CumZ += e.cboxz
				}
			} else {
				if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Pre.CumZ {
					e.smallestZ.Pre.CumX += e.cboxx
				} else if e.smallestZ.CumZ+e.cboxz == e.smallestZ.Post.CumZ {
					e.itemsToPack[e.cboxi].CoordX = e.smallestZ.CumX - e.cboxx
					e.smallestZ.CumX -= e.cboxx
				} else {
					e.smallestZ.Pre.Post = &ScrapPad{}

					e.smallestZ.Pre.Post.Pre = e.smallestZ.Pre
					e.smallestZ.Pre.Post.Post = e.smallestZ
					e.smallestZ.Pre = e.smallestZ.Pre.Post
					e.smallestZ.Pre.CumX = e.smallestZ.Pre.Pre.CumX + e.cboxx
					e.smallestZ.Pre.CumZ += e.cboxz
				}
			}
		}

		e.VolumeCheck()
	}
}

func (e *EB_AFIT) Report(container entities.Container) {
	e.quit = false
	var px, py, pz float64
	log.Printf("%.f", px)
	switch e.bestVariant {
	case 1:
		px, py, pz = container.Length, container.Height, container.Width
	case 2:
		px, py, pz = container.Width, container.Height, container.Length
	case 3:
		px, py, pz = container.Width, container.Length, container.Height
	case 4:
		px, py, pz = container.Height, container.Length, container.Width
	case 5:
		px, py, pz = container.Length, container.Width, container.Height
	case 6:
		px, py, pz = container.Height, container.Width, container.Length
	}

	e.packingBest = true

	e.layers = []*Layer{&Layer{LayerEval: -1}}
	e.ListCanditLayers()
	sort.Slice(e.layers, func(i, j int) bool {
		return e.layers[i].LayerEval < e.layers[j].LayerEval
	})
	e.packedVolume = 0
	e.packedy = 0
	e.packing = true
	e.layerThickness = e.layers[e.bestIteration].LayerDim
	e.remainpy = py
	e.remainpz = pz

	for x := 1; x <= e.itemsToPackCount; x++ {
		e.itemsToPack[x].IsPacked = false
	}

	for e.packing && !e.quit {
		layerinlayer := 0.0
		e.layerDone = false
		e.PackLayer()
		e.packedy += e.layerThickness
		e.remainpy = py - e.packedy

		if math.Abs(layerinlayer-0.0001) > 0 {
			prepackedy := e.packedy
			preremainpy := e.remainpy
			e.remainpy = e.layerThickness - e.prelayer
			e.packedy = e.packedy - e.layerThickness + e.prelayer
			e.remainpz = e.lilz
			e.layerThickness = layerinlayer
			e.layerDone = false
			e.PackLayer()
			e.packedy = prepackedy
			e.remainpy = preremainpy
			e.remainpz = pz
		}

		if !e.quit {
			e.FindLayer(e.remainpy)
		}
	}
}

// / <summary>
// / After packing of each item, the 100% packing condition is checked.
// / </summary>
func (e *EB_AFIT) VolumeCheck() {
	e.itemsToPack[e.cboxi].IsPacked = true
	e.itemsToPack[e.cboxi].PackDimX = e.cboxx
	e.itemsToPack[e.cboxi].PackDimY = e.cboxy
	e.itemsToPack[e.cboxi].PackDimZ = e.cboxz
	e.packedVolume += e.itemsToPack[e.cboxi].Volume
	e.packedItemCount++

	if e.packingBest {
		e.OutputBoxList()
	} else if e.packedVolume == e.totalContainerVolume || e.packedVolume == e.totalItemVolume {
		e.packing = false
		e.hundredPercentPacked = true
	}
}

// #endregion Private Methods

// #region Private Classes

/// <summary>
/// A list that stores all the different lengths of all item dimensions.
/// From the master's thesis:
/// "Each Layerdim value in this array represents a different layer thickness
/// value with which each iteration can start packing. Before starting iterations,
/// all different lengths of all box dimensions along with evaluation values are
/// stored in this array" (p. 3-6).
/// </summary>

/// <summary>
/// From the master's thesis:
/// "The double linked list we use keeps the topology of the edge of the
/// current layer under construction. We keep the x and z coordinates of
/// each gap's right corner. The program looks at those gaps and tries to
/// fill them with boxes one at a time while trying to keep the edge of the
/// layer even" (p. 3-7).
/// </summary>
