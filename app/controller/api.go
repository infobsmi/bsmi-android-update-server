package controller

import (
	"database/sql"
	"fmt"
	"github.com/cnmade/bsmi-android-update-server/pkg/common"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"net/http"
	"strconv"
	"strings"
)

type Api struct {
}

type updateItem struct {
	Link   string `json:"link"`
	Guid string `json:"guid"`
	Ts  string `json:"ts"`
	Version string `json:"version"`
}

func (a *Api) Index(c *gin.Context) {

	//fi, _ := os.Open("./feed.xml")
//	defer fi.Close()
	fp := gofeed.NewParser()
	//feed, _ := fp.Parse(fi)
	feed, _ := fp.ParseURL("https://apkcombo.com/latest-updates/feed")

	var updateList []updateItem
	for _, item := range feed.Items {
		tmpVersion := ""
		tmpA := strings.Split(item.Content, "<br/>")
		if len(tmpA) > 2 {
			rawVersion := tmpA[1]
			if (rawVersion != "") {
				tmpVersion = rawVersion[9:]
			}
		}
		updateList = append(updateList, updateItem{
			Link: item.Link,
			Guid: item.GUID,
			Ts: item.Published,
			Version: tmpVersion,
		} )
	}

	c.JSON(http.StatusOK, updateList)
}

type apiBlogItem struct {
	Aid     string `form:"aid" json:"aid"  binding:"required"`
	Title   string `form:"title" json:"title"  binding:"required"`
	Content string `form:"content" json:"content"  binding:"required"`
}

func (a *Api) View(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("id"))
	fmt.Println(aid)
	if err != nil {
		common.Sugar.Fatal(err)
	}
	var b apiBlogItem

		rows, err := common.DB.Query("Select aid, title, content from gs_article where aid =  ? limit 1 ", &aid)
		if err != nil {
			common.Sugar.Fatal(err)
		}
		defer common.CloseRowsDefer(rows)
		if rows != nil {
			var (
				aid     sql.NullString
				title   sql.NullString
				content sql.NullString
			)
			for rows.Next() {
				err := rows.Scan(&aid, &title, &content)
				if err != nil {
					fmt.Println(err)
				}
				b.Aid = aid.String
				b.Title = title.String
				b.Content = content.String
			}
			fmt.Println(b)
			err = rows.Err()
			if err != nil {
				common.Sugar.Fatal(err)
			}
		}
	fmt.Println(b)
	c.JSON(http.StatusOK, b)
}
