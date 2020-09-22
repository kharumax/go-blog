package repository

import "github.com/jmoiron/sqlx"

// repository内でのグローバル変数を宣言する
var db *sqlx.DB

// 引数にデータベースの接続情報を持った構造体を受け取り、グローバル変数にセットする
func SetDB(d *sqlx.DB)  {
	db = d
}