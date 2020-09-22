package repository

import (
	"database/sql"
	"time"
	"math"

	"go-tech-blog/model"
)

// ArticleList
func ArticleListByCursor(cursor int) ([]*model.Article,error) {
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	query := `SELECT *
    FROM articles
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`
	// クエリの結果を格納するスライスを作成
	// サイズとキャパシティを指定
	articles := make([]*model.Article,0,10)

	if err := db.Select(&articles,query,cursor); err != nil {
		return nil,err
	}
	return articles,nil
}

// ArticleCreate
func ArticleCreate(article *model.Article) (sql.Result,error) {
	now := time.Now()
	article.Created = now
	article.Updated = now

	query := `INSERT INTO articles (title,body,created,updated) VALUES (:title,:body,:created,:updated);`

	// トランザクションの開始
	tx := db.MustBegin()

	// クエリと構造体を引数に渡して、SQLを実行する
	// クエリ内のフィールドは、構造体の定義に沿って変換される(NamedExecの処理）
	res,err := tx.NamedExec(query,article)

	if err != nil {
		tx.Rollback()

		return nil,err
	}
	tx.Commit()

	return res,nil
}

func ArticleDelete(id int) error  {
	query := "DELETE FROM articles WHERE id = ?"
	tx := db.MustBegin()
	// 構造体と紐付けない場合は、Execで
	if _,err := tx.Exec(query,id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func ArticleGetByID(id int) (*model.Article,error)  {
	query := `SELECT * FROM articles WHERE id = ?;`
	var article model.Article
	// 複数取得する場合は、db.Select()で単数の場合は、db.Get()
	if err := db.Get(&article,query,id); err != nil {
		return nil,err
	}
	return &article,nil
}

func ArticleUpdate(article *model.Article) (sql.Result,error)  {
	now := time.Now()
	article.Updated = now
	query := `UPDATE articles
	SET title = :title,body = :body,updated = :updated
	WHERE id = :id
    `
	tx := db.MustBegin()

	res,err := tx.NamedExec(query,&article)
	if err != nil {
		tx.Rollback()
		return nil,err
	}
	tx.Commit()
	return res,nil
}