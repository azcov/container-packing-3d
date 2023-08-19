package entities

type Container struct {
	ID     int
	Length float64
	Width  float64
	Height float64
	Volume float64
}

func NewContainer(id int, length, width, height float64) Container {
	container := Container{
		ID:     id,
		Length: length,
		Width:  width,
		Height: height,
		Volume: length * width * height,
	}
	return container
}
