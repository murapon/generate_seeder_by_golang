package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	var header []string

	// コマンドラインからcsvファイル名、テーブル名を取得する
	csvPath := flag.String("csv", "", "csvファイル名")
	tableName := flag.String("table", "", "テーブル名")
	flag.Parse()

	// CSVファイルを読み込む
	csvFile, err := os.Open(*csvPath)
	if err != nil {
		fmt.Println("error")
	}
	defer csvFile.Close()

	// Seeder生成用テンプレートファイルを読み込む
	templateFile, err := os.Open("seederTemplate.php")
	if err != nil {
		fmt.Println("error")
	}
	defer templateFile.Close()
	b, err := ioutil.ReadAll(templateFile)
	seederTemplate := string(b)

	// クラス名を生成
	className := snakeToUpper(*tableName) + "TableSeeder"

	seederReplaced := strings.Replace(seederTemplate, "{{CLASS_NAME}}", className, 1)
	seederReplaced = strings.Replace(seederReplaced, "{{TABLE_NAME}}", *tableName, 1)

	// CSVデータを取り出す
	csvData := csv.NewReader(csvFile)
	csvData.Comma = ','
	csvData.LazyQuotes = true
	isHeader := true
	fmt.Println(*tableName)
	count := 0
	recordList := ""
	for {
		record, err := csvData.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// 最初の1行目はヘッダーとして読むこみ、2行目以降はデータとして読み込む
		if isHeader == true {
			header = record
			isHeader = false
		} else {
			column := ""
			for i := 0; i < len(header); i++ {
				// カラムによって、現在日時などPHPの処理を入れたい場合は分岐処理を入れる
				if record[i] == "CarbonImmutable::now()" {
					column = column + "                    '" + header[i] + "' => " + record[i] + ",\n"
				} else {
					column = column + "                    '" + header[i] + "' => '" + record[i] + "',\n"
				}
				//                column = column + "                    '" + header[i] +
				//"' => '" + record[i] + "',\n"
			}
			recordList = recordList +
				"            " + strconv.Itoa(count) + " =>\n" +
				"                array(\n" + column +
				"                ),\n"
			count = count + 1
		}
	}
	seederReplaced = strings.Replace(seederReplaced, "{{SEEDER_RECORD}}", recordList, 1)

	// ファイル書き込み
	fSeeder, err := os.Create(className + ".php")
	data := []byte(seederReplaced)
	fSeeder.Write(data)
}

func snakeToUpper(tableName string) string {
	className := strings.Replace(tableName, "_", " ", -1)
	className = strings.Title(className)
	className = strings.Replace(className, " ", "", -1)
	return className
}
