package capsule

import "fmt"

var validTransitions = map[string][]string{
	"":          {"preparing"},
	"preparing": {"running", "error"},
	"running":   {"pausing", "error"},
	"pausing":   {"paused", "error"},
	"paused":    {"resuming", "deleting", "error"},
	"resuming":  {"running", "error"},
	"error":     {"paused", "deleting", "running"},
	"deleting":  {"", "error"},
}

func ValidTransition(from, to string) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == to {
			return true
		}
	}
	return false
}

func (s *State) StatusPrevious() string {
	switch s.Status {
	case "running", "pausing":
		return "running"
	case "paused", "resuming":
		return "paused"
	case "error":
		return "error"
	case "deleting":
		return "deleting"
	default:
		return ""
	}
}

func ValidateTransition(from, to string) error {
	if !ValidTransition(from, to) {
		return fmt.Errorf("invalid state transition: %s -> %s", from, to)
	}
	return nil
}
