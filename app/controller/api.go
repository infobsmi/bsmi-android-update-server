package controller

import (
	"database/sql"
	"fmt"
	"github.com/cnmade/bsmi-android-update-server/app/orm/model"
	"github.com/cnmade/bsmi-android-update-server/app/vo"
	"github.com/cnmade/bsmi-android-update-server/pkg/common"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

type Api struct {
}


func (a *Api) Index(c *gin.Context) {

	var reqJson vo.IndexReq


	 updateInfoList := make([]model.ApkPackage, 0)

	err := c.BindJSON(&reqJson)
	if err != nil {
		common.Sugar.Info("parse json error")
		c.JSON(http.StatusOK, updateInfoList)
		return
	}
	if reqJson.AppList == nil {
		common.Sugar.Info("reqJson.AppList nil")
		c.JSON(http.StatusOK, updateInfoList)
		return
	}
	common.NewDb.Limit(1000).
		Where("guid in ? ", reqJson.AppList).
		Find(&updateInfoList)

	c.JSON(http.StatusOK, updateInfoList)

}


func (a *Api) All(c *gin.Context) {

	var updateInfoList []model.ApkPackage;
	common.NewDb.Limit(1000).Find(&updateInfoList);

	c.JSON(http.StatusOK, updateInfoList)
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
