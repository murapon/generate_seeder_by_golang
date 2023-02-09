# generate_seeder_by_golang
GO言語を用いたLaravel Seeder生成処理
csvファイルから生成する方法と、DBのデータから生成する方法の2種類。

## CSVファイルから生成する方法
./test.csvのような、1行目に生成するテーブルカラム、2行目以降に各カラムに入れたいデータがあるcsvファイルを用意する。
```
go run make_seeder.go -csv=[csvファイル名] -table=[テーブル名]
```
上記コマンドを実行する。

### 実行例
./test.csvを使い、
```
go run make_seeder.go -csv=test.csv -table=sample
```
と実行すると、sampleテーブルに、csvファイルからデータを投入するSeederができる。

## 　DBテーブルから生成する方法
事前にDBテーブルを用意。
make_seeder_db.goのデータ接続部分
```
dbconf := "root:root@tcp(127.0.0.1:53306)/sample?charset=utf8mb4"
```
を実際の環境に合わせて修正する。
```
go run make_seeder_db.go -table=[テーブル名]
```
上記コマンドを実行する。

## Seederのテンプレート
Seederの元になるテンプレートは、`seederTemplate.php`で定義されているので、
変えたい場合は、適宜修正する




