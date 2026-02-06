package dungeondice

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrUnknownClass indicates an unsupported class selection.
var ErrUnknownClass = errors.New("unknown class")

// ActionBonus captures class-specific action modifiers.
type ActionBonus struct {
	Attack  int
	Defense int
	Heal    int
}

// Class defines a playable class profile.
type Class struct {
	Name    string
	MaxHP   int
	Defense int
	Bonus   ActionBonus
}

var classList = []Class{
	{
		Name:    "Warrior",
		MaxHP:   14,
		Defense: 2,
		Bonus:   ActionBonus{Attack: 2, Defense: 1, Heal: 0},
	},
	{
		Name:    "Rogue",
		MaxHP:   11,
		Defense: 1,
		Bonus:   ActionBonus{Attack: 3, Defense: 0, Heal: 1},
	},
	{
		Name:    "Mystic",
		MaxHP:   10,
		Defense: 0,
		Bonus:   ActionBonus{Attack: 1, Defense: 0, Heal: 3},
	},
}

var classIndex = func() map[string]Class {
	index := make(map[string]Class, len(classList))
	for _, class := range classList {
		index[strings.ToLower(class.Name)] = class
	}
	return index
}()

// ClassByName returns the matching class definition.
func ClassByName(name string) (Class, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return Class{}, fmt.Errorf("%w: empty", ErrUnknownClass)
	}
	key := strings.ToLower(trimmed)
	class, ok := classIndex[key]
	if !ok {
		return Class{}, fmt.Errorf("%w: %s", ErrUnknownClass, trimmed)
	}
	return class, nil
}

// ClassNames returns a sorted list of available class names.
func ClassNames() []string {
	names := make([]string, 0, len(classIndex))
	for name := range classIndex {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
