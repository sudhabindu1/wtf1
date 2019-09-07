package main

import (
	"encoding/json"
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/sudhabindu1/wtf1/models"
	"github.com/sudhabindu1/wtf1/modules"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init()  {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := gin.Default()
	r.HTMLRender = ginview.Default()
	r.GET("/", getRandomMessage)
	r.GET("/index", index)
	r.GET("/message/:uid", getMessageWithId)
	r.POST("/insert", insertMessage)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080
}

func index(c *gin.Context)  {
	m, err := modules.FindMessage()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.HTML(http.StatusOK, "page.html", gin.H{
		"message": m.Message,
		"link": m.Link,
		"color": m.Color,
	})
}

func insertMessage(c *gin.Context)  {
	if c.Request.Body != nil {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		m := models.RadioMessage{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			c.String(http.StatusBadRequest, "uid should not be sent")
		}
		if m.Uid != "" {
			c.String(http.StatusBadRequest, "uid should not be sent")
		}
		m.Uid = RandStringRunes(10)
		insertId, err := modules.InsertMessage(&m)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, fmt.Sprintf("created: %v", insertId))
	}
}


func getMessageWithId(c *gin.Context)  {
	uid := c.Param("uid")
	m, err := modules.FindMessageWithId(uid)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, *m)
	return

}

func getRandomMessage(c *gin.Context)  {
	m, err := modules.FindMessage()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, *m)
	return

}
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}