package models

import (
	"github.com/robfig/cron/v3"
	"fmt"
)

// Basic structure to pass between programs
type JobsInfo struct {
	ID int `json:"ID"`
	Url string `json:"url"`
	Props string `json:"props"` //parameter to program
	Pid int `json:"PID"`//handler
	CINFO string `json:"CINFO"`

	//extras
	Handler cron.EntryID `json:"-"`//* cron hadle
}

func LoadAllJobs() []JobsInfo {
	jobs := []JobsInfo{}
	//load contents
	rows, err := Db.Query("SELECT * FROM jobs")
	defer rows.Close()
    if err != nil {
        fmt.Println("Not found")
    }

    for rows.Next() {
        k := JobsInfo{}
        rows.Scan(&k.ID, &k.Url, &k.Props, &k.Pid, &k.CINFO)
        //fmt.Println("info",k)
		jobs = append(jobs, k)
    }
	return jobs
}

func CreateNewJob(j *JobsInfo) {
	stm,_ := Db.Prepare("INSERT INTO jobs(url,CINFO,props,handler) values(?,?,?,?)")
	defer stm.Close()
	res,_ := stm.Exec(j.Url,j.CINFO,j.Props,j.Pid)
	id,_ := res.LastInsertId()
	j.ID = int(id)
	//fmt.Println(j)
}


func SaveJob(j* JobsInfo){
	stm,_ := Db.Prepare("UPDATE jobs SET url=?, CINFO=?, props=?, handler=? WHERE ID =  ?")
	defer stm.Close()
	stm.Exec(j.Url,j.CINFO,j.Props,j.Pid,j.ID)
}

func RemoveJob(j* JobsInfo){
	stm,_ := Db.Prepare("DELETE FROM jobs WHERE ID = ?")
	defer stm.Close()
	stm.Exec(j.ID)
}