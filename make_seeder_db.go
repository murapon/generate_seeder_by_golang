package main

import(
    "database/sql"
	"flag"
	"fmt"
	"io/ioutil"
    "os"
	"strings"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main(){

	// コマンドラインからテーブル名を取得する
	tableName := flag.String("table", "", "テーブル名")
	flag.Parse()

	// DB接続
	dbconf := "root:root@tcp(127.0.0.1:53306)/sample?charset=utf8mb4"
    db, err := sql.Open("mysql", dbconf)
    // 接続が終了したらクローズする
    defer db.Close()  

    if err != nil {
        fmt.Println(err.Error())
    }
    err = db.Ping()
    if err != nil {
		fmt.Println("データベース接続失敗")
        return
    }

	// テーブルのカラム名を取得
	var columns []string
	rows, err := db.Query("SHOW COLUMNS FROM users;")
	if err != nil {
		fmt.Println("getRows db.Query error err:%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var field, fieldType, Null, Key, Extra string
		var Default sql.NullString
		if err := rows.Scan(&field, &fieldType, &Null, &Key, &Default, &Extra); err != nil {
			fmt.Println(err)
			return
		}
		columns = append(columns, field)
	}

	// データを取得
	var records [][]string
	rows, err = db.Query("SELECT * FROM users;")
	if err != nil {
		fmt.Println("getRows db.Query error err:%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			return
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
			return
		}
		var record []string
		for _, col := range values {
			record = append(record, string(col))
		}
		records = append(records, record)
	}

	// Seeder生成用テンプレートファイルを読み込む
    templateFile, err := os.Open("seederTemplate.php")
    if err != nil{
        fmt.Println("error")
    }
    defer templateFile.Close()
	b, err := ioutil.ReadAll(templateFile)
	seederTemplate := string(b)

    // クラス名を生成
	className := snakeToUpper(*tableName) + "TableSeeder"

	seederReplaced := strings.Replace(seederTemplate, "{{CLASS_NAME}}", className, 1)
	seederReplaced = strings.Replace(seederReplaced, "{{TABLE_NAME}}", *tableName, 1)

	count := 0
	seederList := ""
	for i := 0; i < len(records); i++ {
		column := ""
		record := records[i]
		for i := 0; i < len(columns); i++ {
		    column = column + "                    '" + columns[i] + "' => '" + record[i] + "',\n"
		}
		seederList = seederList +
		"            " + strconv.Itoa(count) + " =>\n" +
		"                array(\n" + column + 
		"                ),\n"
	}
	seederReplaced = strings.Replace(seederReplaced, "{{SEEDER_RECORD}}", seederList, 1)

    // ファイル書き込み
	fSeeder, err := os.Create(className + ".php")
    data := []byte(seederReplaced)
    fSeeder.Write(data)
}

func snakeToUpper(tableName string) string {
	className := strings.Replace(tableName, "_", " ", -1)
	className = strings.Title(className)
	className = strings.Replace(className, " ", "", -1)
    return className;
}