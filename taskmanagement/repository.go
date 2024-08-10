package taskmanagement

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/marcos-venicius/daily-term/cycleparser"
)

type Repository struct {
	file *os.File
}

func (r *Repository) SaveBoard(board *Board) error {
	if board.root == nil {
		r.file.Truncate(0)
		r.file.Seek(0, 0)

    return nil
	}

	v, err := cycleparser.ToValue(board.root)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(v)

	if err != nil {
		return err
	}

	r.file.Truncate(0)
	r.file.Seek(0, 0)

	l, err := r.file.Write(bytes)

	if err != nil {
		return err
	}

	if l != len(bytes) {
		return errors.New("Could not save the current board")
	}

	return nil
}

func (r *Repository) LoadBoard(board *Board) {
	stat, err := r.file.Stat()

	if stat.Size() == 0 {
		return
	}

	if err != nil {
		panic(err)
	}

	var bytes []byte = make([]byte, stat.Size())

	readSize, err := r.file.Read(bytes)

	if err != nil {
		panic(err)
	}

	if readSize != int(stat.Size()) {
		panic(fmt.Sprintf("Size file length %d, but read %d", stat.Size(), readSize))
	}

	data := &cycleparser.Value{}

	err = json.Unmarshal(bytes, data)

	if err != nil {
		log.Fatal(err)
	}

	root := &Task{}

	err = cycleparser.FromValue(data, root)

	if err != nil {
		panic(err)
	}

	board.root = root
	board.task = board.root

	current := board.root

	for current != nil {
		board.idCluster.MarkAsUsed(current.Id)

		current = current.Next
	}
}
