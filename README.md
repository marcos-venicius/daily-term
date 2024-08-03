# Daily Term

Manage your daily tasks using your terminal with **vim motions**

https://github.com/user-attachments/assets/c6dc62ec-f643-42dc-aab2-dbfe32f8f7c6

## Modes

- `NORMAL`
- `COMMAND`
- `DELETE`

## NORMAL mode keybindings

- <kbd>q</kbd> quit
- <kbd>:</kbd> enter `COMMAND` mode
- <kbd>k</kbd> previous task
- <kbd>j</kbd> next task
- <kbd>d</kbd> enter `DELETE` mode
- <kbd>t</kbd> move task to state `Todo`
- <kbd>i</kbd> move task to state `In Progress`
- <kbd>c</kbd> move task to state `Completed`
- <kbd>Esc</kbd> clear error

## DELETE mode keybindings

- <kbd>d</kbd> delete current selected task
- <kbd>Esc</kbd> cancel `DELETE` mode

## COMMAND mode commands

- `quit` quit
- `new task "<task name>"` create a new task
- `delete task` delete current selected task
- `delete task <id (int)>` delete task by id
- <kbd>Esc</kbd> cancel `COMMAND` mode
