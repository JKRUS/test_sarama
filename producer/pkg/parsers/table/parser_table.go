package table

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type table struct {
	line    map[int]string
	columns map[int][]int
}

func NewTable() *table {
	return &table{
		line:    make(map[int]string),
		columns: make(map[int][]int),
	}
}

func (t *table) ReadTable(path string) (map[int][]int, error) {
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
	return t.columns, nil
}

func (t *table) splitString(s string) []int {
	words := strings.Fields(s)
	var sw []int
	for _, word := range words {
		i, _ := strconv.Atoi(word)
		sw = append(sw, i)
	}
	return sw
}
