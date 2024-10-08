# Daily Term

Manage your daily tasks using your terminal with **vim motions**

https://github.com/user-attachments/assets/f5884bdf-b8e4-49df-b335-54aeab1ccd5e

## Install

```bash
go install github.com/marcos-venicius/daily-term@latest
```

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

- `q` `quit` quit
- `nt "<task name>"` `new task "<task name>"` create a new task
- `dt` `delete task` delete current selected task
- `dt <id (int)>` `delete task <id (int)>` delete task by id
- <kbd>Esc</kbd> cancel `COMMAND` mode
