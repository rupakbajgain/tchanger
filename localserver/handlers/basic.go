package handlers

import (
	"localserver/models"
	"net/url"
	"io/ioutil"
	"net/http"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type BasicHandler struct{}

func (BasicHandler) name() string {
	return "Default Watcher"
}

func (b BasicHandler) init(j *models.JobsInfo) {
}

func (b BasicHandler) cleanup(j *models.JobsInfo) {
	models.RemoveHash(j.Url)
}

func (b BasicHandler) testHandle(j *models.JobsInfo) bool {
	if !strings.Contains(j.Url, "://"){
		j.Url="http://"+j.Url;
	}
	_, err := url.ParseRequestURI(j.Url)
	if err!=nil{
		return false
	}
	if j.CINFO==""{
		j.CINFO="@every 10m"
	}
	fmt.Println(j)
	b.init(j)
	return true
}

func (b BasicHandler) checkHashes(j *models.JobsInfo,notify bool, hash string){
	if notify{
		previousHash := models.GetHash(j.Url)
		if (previousHash!=hash){
			Notify(b.name(), "Page updated: "+j.Url)
			models.SaveHash(j.Url, hash)
		}
	}else{
		models.SaveHash(j.Url, hash)
	}
}

func (b BasicHandler) exec(j *models.JobsInfo,notify bool){
	//fmt.Println("fetching", j.Url)
	resp, err := http.Get(j.Url)
	if err!=nil{
		fmt.Println("Error getting", j.Url)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("Error reading", j.Url)
		return
	}
	hash := md5.Sum([]byte(body))
	textHash := hex.EncodeToString(hash[:])
	//fmt.Println(textHash)
	b.checkHashes(j, notify, textHash)// of after parse hash use it here too
}