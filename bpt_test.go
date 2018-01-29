package bpt

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"testing"
)

func TestInsertAndRead(t *testing.T) {
	file, err := os.Open("records.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	tree := new(Node)
	tree.Init(Record{"Joan", "Jett"})
	for {

		kv, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record := Record{key: kv[0], value: kv[1]}
		tree, err = Insert(tree, record)
		if err != nil {
			t.Fatal(err)
		}
		result, ok, err := Find(tree, record.key)
		if err != nil {
			t.Fatal("Call to Find() returned an error", err)
		}
		if !ok {
			t.Fatal("Record should have been found, but it wasn't", record.key)
		} else {
			if record != result {
				t.Fatal("Mismatch", record, result)
			}
		}
	}
}
