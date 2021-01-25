package models

import (
	"container/list"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	_SQL_DB *sql.DB
	_ERR error
)

type RuleItem struct {
	Name	string
	Fn      string
	Interval  int
	Alert     string
	Expr 	 string
	For      string
	Labels   map[string]string
	Annotations map[string]string
}


func InitDB(db_url string){
	if nil != _SQL_DB {
		return
	}

	_SQL_DB, _ERR = sql.Open("mysql", db_url)
	if _ERR != nil {
		fmt.Errorf("Open mysql database error: %s\n", _ERR)
		return
	}
	_ERR = _SQL_DB.Ping()
	if _ERR != nil {
		fmt.Errorf("%s",_ERR)
		return
	}
}

func QueryRulesFromMysql() (*list.List, error) {
	var (
		rule_labels, rule_annotations  string
	)

	l := list.New()

	queryString := "select rule_name, rule_fn, rule_interval, rule_alert, rule_expr, rule_for, rule_labels, rule_annotations from rules;"
	rows, err := _SQL_DB.Query(queryString)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer  rows.Close()

	for rows.Next(){
		var item RuleItem
		item.Labels = make(map[string]string)
		item.Annotations = make(map[string]string)

		err := rows.Scan(&item.Name, &item.Fn, &item.Interval, &item.Alert, &item.Expr, &item.For, &rule_labels, &rule_annotations)
		if err != nil{
			fmt.Errorf("rows.Scan failed : %s",err)
			continue
		}

		labels := strings.Split(rule_labels, ",")
		lablen := len(labels)
		for i:=0; i<lablen; i++{
			parts := strings.Split(labels[i], "=")
			plen := len(parts)

			for j:=0; j<plen; j+=2 {
				item.Labels[parts[j]] = parts[j+1]
			}
		}

		annotations := strings.Split(rule_annotations, ",")
		annlen := len(annotations)
		for k:=0; k<annlen; k++{
			parts := strings.Split(annotations[k], "=")
			plen := len(parts)
			for j:=0; j<plen; j+=2 {
				item.Annotations[parts[j]] = parts[j+1]
			}
		}

		l.PushBack(item)
	}
	return l,err
}
