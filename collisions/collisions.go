package collisions

import (
	"encoding/json"
	"image"
	"io"
	"log"
	"os"
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
	Music     string
}

// Used for json unmarsharling
type ObstacleData struct {
	X1, Y1 int
	X2, Y2 int
}

type DiagonalObstacleData struct {
	StartX, StartY int
	Width, Height  int
	Count          int
}

type DoorData struct {
	X1, Y1      int
	X2, Y2      int
	NewX, NewY  float64
	Destination string
	Id          string
}

type CollisionData struct {
	Obstacles []ObstacleData
	Diagonals []DiagonalObstacleData
	Doors     []DoorData
	Music     string
}

func New(path string) Collisions {
	data := loadJsonFile(path)
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
	collisions.Music = data.Music
	return collisions
}

func loadJsonFile(path string) *CollisionData {
	//Loading json collsion file
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	data := &CollisionData{}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, data)
	return data
}
