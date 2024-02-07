package ability

import "fmt"

type AbilityType int

const (
	None AbilityType = iota
	GhostMode
	StopTime
	MaxAbility
)

type Ability struct {
	Type      AbilityType
	Activated bool
}

func (a *Ability) CycleAbility() {
	// Increment the current ability type
	a.Type++

	// Check if we've exceeded the bounds of defined abilities, excluding MaxAbility
	if a.Type >= MaxAbility {
		a.Type = None + 1 // Assuming you want to skip 'None'. If not, just reset to None.
	}

	a.Activated = false // Reset the active state when cycling
	fmt.Println(a.Type)
}

func (a *Ability) ActivateAbility() {
	fmt.Println("Activated")
	if a.Type != None {
		a.Activated = true
	}
}

func (a *Ability) DeactivateAbility() {
	a.Activated = false
}
