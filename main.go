package main
//
import (
	. "github.com/cnmade/bsmi-android-update-server/app/controller"
	"github.com/cnmade/bsmi-android-update-server/pkg/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"

	"gitee.com/cnmade/pongo2gin"
)

func main() {
	common.InitApp()
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
		api.GET("/", a.Index)
		api.GET("view/:id", a.View)
	}
	log.Info().Msg("Server listen on 127.0.0.1:8086")
	err := r.Run("127.0.0.1:8086")
	if err != nil {
		common.LogError(err)
	}
}
