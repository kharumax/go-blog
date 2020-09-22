package main

import (
     "log"
     "os"

     "go-tech-blog/handler"
     "go-tech-blog/repository"

     "github.com/labstack/echo/v4"
     _ "github.com/go-sql-driver/mysql" // Using MySQL driver
     "github.com/jmoiron/sqlx"
     "github.com/labstack/echo/v4/middleware"
     "gopkg.in/go-playground/validator.v9"
)



var db *sqlx.DB //変数定義（グローバル）
var e = createMux()

func main() {
     db = connectDB() //ここでグローバル変数にDBを入れる
     repository.SetDB(db)

     e.GET("/",handler.ArticleIndex)

     e.GET("/articles",handler.ArticleIndex)
     e.GET("/articles/new",handler.ArticleNew)
     e.GET("/articles/:articleID",handler.ArticleShow)
     e.GET("/articles/:articleID/edit",handler.ArticleEdit)

     // HTMLではなく、JSONを返す処理はapiとする
     e.GET("/api/articles",handler.ArticleList)
     e.POST("/api/articles",handler.ArticleCreate)
     e.DELETE("/api/articles/:articleID",handler.ArticleDelete)
     e.PATCH("/api/articles/:articleID",handler.ArticleUpdate)
     // Webサーバーを8080で起動する
     port := os.Getenv("PORT")
     if port == "" {
          e.Logger.Fatal("$PORT must be set")
     }
     e.Logger.Fatal(e.Start(":"+port))
}

func connectDB() *sqlx.DB  {
     dsn := os.Getenv("DSN")
     db,err := sqlx.Open("mysql",dsn)
     if err != nil {
          e.Logger.Fatal(err)
     }
     if err := db.Ping(); err != nil {
          e.Logger.Fatal(err)
     }
     log.Println("Db connection succeeded")
     return db
}

func createMux() *echo.Echo  {
     //アプリケーションインスタンスの作成
     e := echo.New()

     //各種ミドルウェアの設定
     e.Use(middleware.Recover())
     e.Use(middleware.Logger())
     e.Use(middleware.Gzip())
     e.Use(middleware.CSRF())

     //静的ファイルの設定
     e.Static("/css","src/css")
     e.Static("/js","src/js")

     e.Validator = &CustomValidator{validator: validator.New()}

     return e
}

type CustomValidator struct {
     validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error  {
     return cv.validator.Struct(i)
}
