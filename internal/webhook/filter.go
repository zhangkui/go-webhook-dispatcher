package webhook

func matches(sub Subscription, event Event) bool {
	if !sub.Active {
		return false
	}
	typeMatch := false
	for _, eventType := range sub.EventTypes {
		if eventType == event.Type || eventType == "*" {
			typeMatch = true
			break
		}
	}
	if !typeMatch {
		return false
	}
	for key, expected := range sub.Filter {
		if event.Attributes[key] != expected {
			return false
		}
	}
	return true
}
