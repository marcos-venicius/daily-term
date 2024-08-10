package taskmanagement

func (task *Task) Symbol(currentSelectedTaskId int) rune {
	if task.Id == currentSelectedTaskId {
		switch task.State {
		case Todo, InProgress, Completed:
			return '*'
		default:
			return '*'
		}
	}

	return ' '
}
