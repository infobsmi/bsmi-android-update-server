package main
//
import (
	. "github.com/cnmade/bsmi-android-update-server/app/controller"
	"github.com/cnmade/bsmi-android-update-server/app/orm/model"
	"github.com/cnmade/bsmi-android-update-server/pkg/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"strings"

	"gitee.com/cnmade/pongo2gin"

	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	common.InitApp()

	s := gocron.NewScheduler(time.UTC)

	s.Every(65).Seconds().Do(func(){
		//TODO
		common.Sugar.Info(" I am running now")


		fp := gofeed.NewParser()
		//feed, _ := fp.Parse(fi)
		feed, _ := fp.ParseURL("https://apkcombo.com/latest-updates/feed")


		for _, item := range feed.Items {
			tmpVersion := ""
			tmpA := strings.Split(item.Content, "<br/>")
			if len(tmpA) > 2 {
				rawVersion := tmpA[1]
				if (rawVersion != "") {
					tmpVersion = rawVersion[9:]
				}
			}



			var oldUpdateInfo model.ApkPackage
			common.NewDb.First(&oldUpdateInfo, "guid = ?", item.GUID)
			if oldUpdateInfo.Id  > 0 {
				oldUpdateInfo.Link = item.Link;
				oldUpdateInfo.Ts = item.Published;
				oldUpdateInfo.Version = tmpVersion;
				common.NewDb.Save(oldUpdateInfo)
			} else {
				common.NewDb.Create(&model.ApkPackage{
					Link: item.Link,
					Guid: item.GUID,
					Ts: item.Published,
					Version: tmpVersion,
				})
			}
		}

	})

	s.StartAsync()
	//pstore := persistence.NewInMemoryStore(time.Second)

	r := gin.New()
	r.HTMLRender = pongo2gin.New(pongo2gin.RenderOptions{
		TemplateDir: "views",
		ContentType: "text/html; charset=utf-8",
		AlwaysNoCache: true,
	})

	r.Static("/assets", "./vol/assets")
	store := cookie.NewStore([]byte("gssecret"))
	r.Use(sessions.Sessions("mysession", store))


	a := new(Api)
	api := r.Group("/api")
	{
		api.POST("/",  a.Index)
		api.GET("/all", a.All);
		api.GET("view/:id", a.View)
	}
	log.Info().Msg("Server listen on 127.0.0.1:8086")
	err := r.Run("127.0.0.1:8086")
	if err != nil {
		common.LogError(err)
	}
}
