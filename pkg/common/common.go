package common

import (
	"database/sql"
	"github.com/cnmade/bsmi-android-update-server/app/orm/model"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/grokify/html-strip-tags-go"
	"github.com/ztrue/tracerr"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)


type Msg struct {
	Msg string
}
type Umsg struct {
	Msg string
	Url string
}

type VBlogItem struct {
	Aid            int
	Title          sql.NullString
	Content        sql.NullString
	Publish_time   sql.NullString
	Publish_status sql.NullInt64
	Views          int
}

/**
 * Logging error
 */
func LogError(err error) {
	if err != nil {
		Sugar.Error(tracerr.Sprint(tracerr.Wrap(err)))
	}
}

/**
 * Logging info
 */
func LogInfo(msg string) {
	if msg != "" {
		Sugar.Info(msg)
	}
}

func LogInfoF(msg string, v interface{}) {
	if msg != "" {
		Sugar.Infof(msg, v)
	}
}



/**
 * close rows defer
 */
func CloseRowsDefer(rows *sql.Rows) {
	_ = rows.Close()
}

/*
* ShowMessage with template
 */
func ShowMessage(c *gin.Context, m *Msg) {

	c.HTML(200, "message.html",
		Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        Config.Site_name,
			"siteDescription": Config.Site_description,
			"message":         m.Msg,
		}))
	return
}

func ShowUMessage(c *gin.Context, m *Umsg) {

	c.HTML(200, "message.html",
		Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        Config.Site_name,
			"siteDescription": Config.Site_description,
			"message":         m.Msg,
			"url":             m.Url,
		}))
	return
}

func GetMinutes() string {
	return time.Now().Format("200601021504")
}

func GetNewDb(config *AppConfig) *gorm.DB {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,         // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(config.Dbdsn), &gorm.Config{

		Logger: newLogger,
	})


	if err != nil {
		panic(err.Error())
	}


	err = db.AutoMigrate(&model.ApkPackage{})

	if err != nil {
		panic(err.Error())
	}


	return db
}

type AppConfig struct {
	Dbdsn            string
	Admin_user       string
	Admin_password   string
	Site_name        string
	Site_description string
	SrvMode string
	ObjectStorage    struct {
		Aws_access_key_id     string
		Aws_secret_access_key string
		Aws_region            string
		Aws_bucket            string
		Cdn_url               string
	}
}

func GetConfig() *AppConfig {
	//_cm := "GetConfig@pkg/common/common"
	//TODO load config from cmd line argument
/*	//f, err := os.Open("./vol/config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config AppConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		Sugar.Infof(_cm + " error: %v", err)
	}*/
	config := AppConfig{
		Dbdsn: "file:./db.sqlite",
	}
	return &config
}



var (
	Config    *AppConfig
	DB        *sql.DB
	NewDb     *gorm.DB
	Logger, _ = zap.NewProduction()
	Sugar *zap.SugaredLogger
)

func InitApp() {
	Config = GetConfig()
//	gin.SetMode(Config.SrvMode)
	gin.SetMode("debug")
	//DB = GetDB(Config)
	NewDb = GetNewDb(Config)
	defer Logger.Sync()
	Sugar = Logger.Sugar()
}

func OutPutHtml( c *gin.Context, s string) {
	c.Header("Content-Type", "text/html;charset=UTF-8")
	c.String(200, "%s", s)
	return
}
func OutPutText( c *gin.Context, s string) {
	c.Header("Content-Type", "text/plain;charset=UTF-8")
	c.String(200, "%s", s)
	return
}
/**
 * ???????????????????????????????????????
 */
func SubCutContent(content string, length int) string {
	if len(content) <= length {
		return content
	}

	content = strip.StripTags(content)
	content = strings.TrimSpace(content)
	content = strings.Replace(content, "<!DOCTYPE html>", "", 1)
	content = strings.Replace(content, "&nbsp;", "", 1)

	tmpContent := []rune(content)

	rawLen := len(tmpContent)

	if length > rawLen {
		return content
	}

	return string(tmpContent[0:length])
}

