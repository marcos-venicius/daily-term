package taskmanagement

func (task *Task) Symbol(currentSelectedTaskId int) rune {
	if task.Id == currentSelectedTaskId {
		switch task.State {
		case Todo:
			return 't'
		case InProgress:
			return 'i'
		case Completed:
			return 'c'
		default:
			return '*'
		}
	}

	return ' '
}
