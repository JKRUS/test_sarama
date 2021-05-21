package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type WriteMSG struct {
	db *sqlx.DB
}

func NewMsgPostgres(db *sqlx.DB) *WriteMSG {
	return &WriteMSG{db: db}
}

func (r *WriteMSG) WriteMessage(msg []byte) (int, error) {
	var id int
	tabRes := map[string][]int{
		"0": make([]int, 0),
		"1": make([]int, 0),
	}
	err := json.Unmarshal(msg, &tabRes)
	if err != nil {
		return 0, err
	}
	fmt.Println(tabRes)
	var nowVar, ind int
	var query string
	var row = &sqlx.Row{}
	var columns []int
	for rows := range tabRes {
		sind := strconv.Itoa(ind)
		fmt.Println(tabRes[rows])
		nowVar = 0
		for _ = range tabRes[sind] {
			nowVar++
		}
		columns = tabRes[sind]
		ind++
		if len(columns) == 0 {
			continue
		}
		switch nowVar {
		case 1:
			query = fmt.Sprintf("INSERT INTO %s (first_column) values($1) RETURNING id", dataTable)
			row = r.db.QueryRowx(query, columns[0])
		case 2:
			query = fmt.Sprintf("INSERT INTO %s (first_column, second_column) values($1, $2) RETURNING id", dataTable)
			row = r.db.QueryRowx(query, columns[0], columns[1])
		}
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
