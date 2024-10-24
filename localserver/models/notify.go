package models

import (
	"fmt"
)

type Notify struct {
	ID int
	Title string `json:"title"`
	Body string `json:"body"`
}

func SaveNotify(title string, body string){
	stm,_ := Db.Prepare("INSERT INTO no_log(title,body) values(?,?)")
	defer stm.Close()
	stm.Exec(title, body)
}

func GetNotify() []Notify{
	n := []Notify{}

	rows, err := Db.Query("SELECT * FROM no_log ORDER BY ID desc LIMIT 30")
	defer rows.Close()
    if err != nil {
        fmt.Println("Not found")
    }

    for rows.Next() {
        k := Notify{}
        rows.Scan(&k.ID, &k.Title, &k.Body)
        //fmt.Println("info",k)
		n = append(n, k)
    }
	return n
}