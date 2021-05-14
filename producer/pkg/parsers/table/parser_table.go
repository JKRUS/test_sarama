package table

import (
	"bufio"
	"os"
	"strings"
)

type table struct {
	line    map[int]string
	columns map[int][]string
}

func NewTable() *table {
	return &table{
		line:    make(map[int]string),
		columns: make(map[int][]string),
	}
}

func (t *table) ReadTable(path string) (map[int]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		t.line[i] = scanner.Text()
		w := t.splitString(scanner.Text())
		t.columns[i] = w
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return t.line, nil
}

func (t *table) splitString(s string) []string {
	words := strings.Fields(s)
	var sw = []string{}
	for _, word := range words {
		sw = append(sw, word)
	}
	return sw
}
