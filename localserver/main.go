package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "localserver/models"
    "localserver/handlers"
    "strconv"
)

var jobs[] models.JobsInfo
/*
func getJobs(c *gin.Context) {
    c.JSON(http.StatusOK, jobs)
}
*/
func removeJob(c *gin.Context) {
    id,_ := strconv.Atoi(c.Param("id"))
    for a, _ := range jobs {
        if jobs[a].ID == id {
            handlers.RemoveFromCron(&jobs[a])
            models.RemoveJob(&jobs[a])
            //--
            copy(jobs[a:], jobs[a+1:])
            jobs = jobs[:len(jobs)-1]
            //--
            c.Redirect(http.StatusFound, "/")
            return
        }
    }
}

func runJob(c *gin.Context) {
    id,_ := strconv.Atoi(c.Param("id"))
    for a, _ := range jobs {
        if jobs[a].ID == id {
            handlers.TestExecute(&jobs[a])
            c.Redirect(http.StatusFound, "/")
            return
        }
    }
}

func addJobHelper(c *gin.Context, url string){
    newJob := models.JobsInfo{Url:url}

    st := handlers.AddToCron(&newJob)
    if st==false{
        c.JSON(http.StatusNotFound, gin.H{"message": "Cannot handle"})
        return
    }
    handlers.TouchExecute(&newJob)

    models.CreateNewJob(&newJob)
    jobs = append(jobs, newJob)

    c.Redirect(http.StatusFound, "/")
    //c.JSON(http.StatusOK, newJob)
}

func addJob(c *gin.Context) {
    url := c.Param("url")
    addJobHelper(c, url)
}

func addJob2(c *gin.Context) {
    url := c.Query("url")
    addJobHelper(c, url)
}

func addJobP(c *gin.Context) {
    url := c.PostForm("url")
    addJobHelper(c, url)
}

func getIndex(c *gin.Context){
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "jobs": jobs,
        "hNames": handlers.HandlerNames,
        "notify": models.GetNotify(),
    })
}

func getModify(c *gin.Context){
    id,_ := strconv.Atoi(c.Param("id"))
    for a, _ := range jobs {
        if jobs[a].ID == id {
            c.HTML(http.StatusOK, "modify.tmpl", gin.H{
                "j": jobs[a],
                "hNames": handlers.HandlerNames,
            })
            return
        }
    }
}

func postModify(c *gin.Context){
    id,_ := strconv.Atoi(c.Param("id"))
    url := c.PostForm("url")
    Props := c.PostForm("props")
    handler,_ := strconv.Atoi(c.PostForm("pid"))
    CINFO := c.PostForm("CINFO")
    restart := false
    for a, _ := range jobs {
        if jobs[a].ID == id {
            if(jobs[a].CINFO!=CINFO || jobs[a].Pid!=handler){
                fmt.Println("Restarting cron")
                restart=true
                handlers.RemoveFromCron(&jobs[a])
            }
            jobs[a].Url = url
            jobs[a].Props = Props
            jobs[a].Pid = handler
            jobs[a].CINFO = CINFO
            models.SaveJob(&jobs[a])
            if(restart){
                handlers.InitToCron(&jobs[a])
            }
            c.Redirect(http.StatusFound, "/")
            return
        }
    }
}

func pNotify(c *gin.Context){
    var k models.Notify
    c.BindJSON(&k)
    fmt.Println(k)
    handlers.Notify(k.Title, k.Body)
    c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func main() {
    //gin.SetMode(gin.ReleaseMode)

    router := gin.Default()
    router.LoadHTMLGlob("templates/*.tmpl")

    //router.GET("/jobs", getJobs)
    router.GET("/remove/:id", removeJob)
    router.GET("/run/:id", runJob)
    router.GET("/addjob/:url", addJob)
    router.GET("/addjob", addJob2)
    router.POST("/addjob", addJobP)
    router.GET("/", getIndex)
    router.GET("/modify/:id", getModify)
    router.POST("/modify/:id", postModify)
    router.POST("/notify", pNotify)

    jobs = models.LoadAllJobs()
    for i,_ := range(jobs){
        v:=&jobs[i]
        handlers.InitToCron(v)
        handlers.TestExecute(v)
    }
    fmt.Println("No of handlers:", len(handlers.Handlers))
    router.Run("localhost:8080")
}