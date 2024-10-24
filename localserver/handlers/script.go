package handlers

import (
	"localserver/models"
	"strings"
	//"fmt"
	"os"
	"os/exec"
)

type PythonScriptHandler struct{}

func (PythonScriptHandler) name() string {
	return "Script Runner"
}

func (b PythonScriptHandler) init(j *models.JobsInfo) {
}

func (b PythonScriptHandler) cleanup(j *models.JobsInfo) {
}

func (b PythonScriptHandler) testHandle(j *models.JobsInfo) bool {
	if !strings.HasPrefix(j.Url, "@python:"){
		return false
	}
	if j.CINFO==""{
		j.CINFO="@every 10m"
	}
	b.init(j)
	return true
}

func (b PythonScriptHandler) exec(j *models.JobsInfo,notify bool){
	cmd:=exec.Command("/usr/bin/python", j.Url[8:])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}