package data

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Data struct {
	Obstacles []ObstacleData
	Diagonals []DiagonalObstacleData
	Doors     []DoorData
	NPCs      []NPCData
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
type BehaviorData struct {
	Type    string                 // A string to denote the type of behavior (e.g., "walker", "talker")
	Details map[string]interface{} // Additional details specific to each behavior type
}
type NPCData struct {
	Name         string
	SpriteSheets map[string]string
	FrameCount   int
	X, Y         float64
	Behaviors    []BehaviorData
}

func LoadJsonFile(path string) *Data {
	//Loading json file
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	data := &Data{}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, data)
	return data
}
