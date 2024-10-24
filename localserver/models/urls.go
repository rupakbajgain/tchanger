package models

import (	
	_ "fmt"
)
/*
// to handle hash
type HashInfo struct {
	ID int
	Url string
	Hash string
}
*/
func GetHash(url string) string {
	stm,_ := Db.Prepare("SELECT hash FROM hash WHERE url =  ?")
	defer stm.Close()
	//print("Get hash", url)
	var hash string
	err := stm.QueryRow(url).Scan(&hash)
	if err!=nil{
		stm,_ = Db.Prepare("INSERT INTO hash(url) values(?)")
		stm.Exec(url)
	}
	return hash
}

func SaveHash(url string, hash string){
	stm,_ := Db.Prepare("UPDATE hash SET hash = ? WHERE url =  ?")
	defer stm.Close()
	stm.Exec(hash, url)
}

func InsertHash(url string, hash string){
	stm,_ := Db.Prepare("INSERT INTO hash(url,hash) values(?,?)")
	defer stm.Close()
	stm.Exec(url, hash)
}

func RemoveHash(url string){
	stm,_ := Db.Prepare("DELETE FROM hash WHERE url = ?")
	defer stm.Close()
	stm.Exec(url)
}
