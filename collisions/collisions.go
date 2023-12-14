package collisions

import (
	"image"
	"rpg_demo/data"
)

type Door struct {
	Rect        *image.Rectangle
	Id          string
	Destination string
	NewX, NewY  float64
}

type Collisions struct {
	Obstacles []*image.Rectangle
	Doors     []*Door
}

func New(data *data.Data) Collisions {

	collisions := Collisions{}
	for _, obs := range data.Obstacles {
		i := image.Rect(obs.X1, obs.Y1, obs.X2, obs.Y2)
		collisions.Obstacles = append(collisions.Obstacles, &i)
	}
	for _, d := range data.Diagonals {
		for i := 0; i < d.Count; i++ {
			x1 := d.StartX + (d.Width * i)
			y1 := d.StartY + (d.Height * i)
			x2 := x1 + d.Width
			y2 := y1 + d.Height
			i := image.Rect(x1, y1, x2, y2)
			collisions.Obstacles = append(collisions.Obstacles, &i)
		}
	}
	for _, d := range data.Doors {
		i := image.Rect(d.X1, d.Y1, d.X2, d.Y2)
		d := &Door{
			Rect:        &i,
			Id:          d.Id,
			Destination: d.Destination,
			NewX:        d.NewX,
			NewY:        d.NewY,
		}
		collisions.Doors = append(collisions.Doors, d)
	}
	return collisions
}
