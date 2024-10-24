package handlers

import (
	"localserver/models"
	"github.com/robfig/cron/v3"
	"github.com/gen2brain/beeep"
	"fmt"
)

var c *cron.Cron

type Handler interface {
	name() string
	testHandle(*models.JobsInfo) bool
	exec(*models.JobsInfo,bool)
	init(*models.JobsInfo)
	cleanup(*models.JobsInfo)
}
/*
const (
	Basic  = 0
)
*/
var Handlers[2] Handler
var HandlerNames[2] string//just caching

func init(){
	c = cron.New()
	c.Start()
	
	Handlers = [2]Handler{
		&BasicHandler{},
		&PythonScriptHandler{},
	}

	for i:=0;i<len(Handlers);i++{
		HandlerNames[i]=Handlers[i].name()
	}
}

func functor(f Handler, k *models.JobsInfo) func() {
	//fmt.Println("Creating functor", &k)
	return func(){
		fmt.Println("Auto Executing:", k.Url)
		f.exec(k,true)
	}
}

func AddToCron(k *models.JobsInfo) bool{
	for i:=len(Handlers)-1;i>=0;i--{//inverse order for compatibility
		v:=Handlers[i].testHandle(k)
		if(v){
			//fmt.Println("Accepted",k)
			k.Pid = i
			k.Handler,_ = c.AddFunc(k.CINFO, functor(Handlers[i], k))
			return true
		}
	}
	fmt.Println("Cannot handle", k)
	return false
}

//we know which handler
func InitToCron(k *models.JobsInfo){
	handler := Handlers[k.Pid]
	handler.init(k)
	k.Handler,_ = c.AddFunc(k.CINFO, functor(handler, k))
	//fmt.Println("Initing", k)
}

func RemoveFromCron(k *models.JobsInfo){
	Handlers[k.Pid].cleanup(k)
	c.Remove(k.Handler)
}

func TestExecute(k *models.JobsInfo){
	Handlers[k.Pid].exec(k, true)
}

func TouchExecute(k *models.JobsInfo){
	Handlers[k.Pid].exec(k, false)
}

func Notify(title string, msg string){
	beeep.Notify(title, msg, "assets/information.png")
	models.SaveNotify(title, msg)
}