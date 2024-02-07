package shared

type GameState int

const (
	PlayState GameState = iota
	TransitionState
	NewSceneState
	CutSceneState
	TimeStopped
)

type Transition struct {
	Alpha     float64
	FadeSpeed float64
	Timer     int
	Music     bool
}

type KeyPressed struct {
	KeyP bool
	KeyZ bool
	KeyD bool
	KeyV bool
}
