package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Data struct {
	Id            int `DB:"id"`
	First_column  int `DB:"first_column"`
	Second_column int `DB:"second_column"`
}

func ReadTable() (table []Data, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", dataTable)
	rows, err := DB.Queryx(query)
	if err != nil {
		logrus.Fatalf(err.Error())
		return nil, err
	}
	defer rows.Close()

	row := Data{}

	for rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		table = append(table, row)
	}
	return table, nil
}
