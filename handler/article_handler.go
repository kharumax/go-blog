package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-tech-blog/repository"
	"go-tech-blog/model"
	"github.com/labstack/echo/v4"
)

// ArticleIndex ...
func ArticleIndex(c echo.Context) error {
	// /articlesでのリクエストがあったら、/にリダイレクトする
	// これで、GAなどで計測する時にパスが統一される
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect,"/")
	}
	articles, err := repository.ArticleListByCursor(0)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	data := map[string]interface{}{
		"Articles": articles, // 記事データをテンプレートエンジンに渡す
		"Cursor": cursor,
	}
	return render(c, "article/index.html", data)
}

// ArticleNew ...
func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}

	return render(c, "article/new.html", data)
}

// ArticleShow ...
func ArticleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))
	article,err := repository.ArticleGetByID(id)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/show.html", data)
}

// ArticleEdit ...
func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	aritcle,err := repository.ArticleGetByID(id)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"Article": aritcle,
	}
	return render(c, "article/edit.html", data)
}

type ArticleCreateOutput struct {
	Article *model.Article
	Message string
	ValidationErrors []string
}

func ArticleCreate(c echo.Context) error  {
	//　曹仁されてくるフォームの内容を格納する構造体を宣言
	var article model.Article
	var out ArticleCreateOutput

	// フォームの内容を構造体に埋め込む
	if err := c.Bind(&article); err != nil {
		c.Logger().Error(err.Error())
		// リクエストの解釈に失敗した場合は、400エラーを返す
		return c.JSON(http.StatusBadRequest,out)
	}

	if err := c.Validate(&article); err != nil {
		c.Logger().Error(err.Error())
		out.ValidationErrors = article.ValidationErrors(err)

		return c.JSON(http.StatusUnprocessableEntity,out)
	}
	// 紐付けたデータを元に、SQLを発行する
	res,err := repository.ArticleCreate(&article)
	if err != nil {
		c.Logger().Error(err.Error())
		// サーバー内のエラーが発生した場合は、500エラーを返す
		return c.JSON(http.StatusInternalServerError,out)
	}
	// SQL実行結果から、作成されたレコードのIDを取得
	id,_ := res.LastInsertId()
	article.ID = int(id)
	out.Article = &article

	return c.JSON(http.StatusOK,out)
}

func ArticleDelete(c echo.Context) error  {
	id,_ := strconv.Atoi(c.Param("articleID"))
	if err := repository.ArticleDelete(id); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError,"")
	}
	return c.JSON(http.StatusOK,fmt.Sprintf("Article %d is deleted.",id))
}

func ArticleList(c echo.Context) error  {
	cursor,_ := strconv.Atoi(c.QueryParam("cursor"))
	articles,err := repository.ArticleListByCursor(cursor)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError,"")
	}
	return c.JSON(http.StatusOK,articles)
}

type ArticleUpdateOutput struct {
	Article *model.Article
	Message string
	ValidationErrors []string
}

func ArticleUpdate(c echo.Context) error  {
	// リクエスト送信元のパスを取得
	ref := c.Request().Referer()
	refID := strings.Split(ref,"/")[4]
	reqID := c.Param("articleID")

	if refID != reqID {
		return c.JSON(http.StatusBadRequest,"")
	}
	var article model.Article
	var out ArticleUpdateOutput

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest,out)
	}
	if err := c.Validate(&article); err != nil {
		out.ValidationErrors = article.ValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity,out)
	}
	articleID,_ := strconv.Atoi(reqID)
	article.ID = articleID
	_,err := repository.ArticleUpdate(&article)
	if err != nil {
		out.Message = err.Error()
		return c.JSON(http.StatusInternalServerError,out)
	}
	out.Article = &article
	return c.JSON(http.StatusOK,out)
}
