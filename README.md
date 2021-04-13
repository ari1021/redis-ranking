# redis-ranking
## 概要
このレポジトリは，Redis の Sorted Set を用いたランキング機能を実装したレポジトリです．

MySQL を用いたランキング機能も実装してあるので，Redisを用いた場合 と MySQL を用いた場合で性能比較をすることができます．

## ディレクトリ構成
ディレクトリ構成は以下のようになっています．
```bash
.
├── README.md
├── docker-compose.yml
├── go.mod
├── go.sum
├── init
│   └── main.go # ベンチマーク用に初期データを挿入するファイル
├── mysql
│   ├── init
│   │   └── ddl.sql
│   └── my.cnf
├── redis
│   └── redis.conf
└── src
    ├── db
    │   └── conn.go # MySQL と Redis のコネクションを確立するファイル
    ├── mysql
    │   ├── ranking.go # MySQL を用いたランキング機能を実装したファイル
    │   └── ranking_test.go # MySQL を用いたランキング機能を計測するベンチマークファイル
    └── redis
        ├── ranking.go # Redis を用いたランキング機能を実装したファイル
        └── ranking_test.go # Redis を用いたランキング機能を計測するベンチマークファイル
```

## 使い方
### 1. コンテナを起動する
```
docker-compose build
docker-compose up -d
```

### 2. 環境変数を設定する
```
export MYSQL_USER=root \
    MYSQL_PASSWORD=password \
    MYSQL_HOST=127.0.0.1 \
    MYSQL_PORT=3306 \
    MYSQL_DATABASE=database \
    REDIS_HOST=127.0.0.1 \
    REDIS_PORT=6379
```

### 3. ベンチマーク用に初期データを挿入する
初期データの個数はコマンドライン引数から指定することができます．

初期データの個数を指定しない場合は，1000個になります
```
go run init/main.go [-dataNum={任意の自然数}]
```

### 4. 計測を行う
```
go test -bench=. ./src/redis -benchmem
go test -bench=. ./src/mysql -benchmem
```

## 計測結果
CPU が 1.8GHz Intel Core i5，メモリが 8GB の計算機上で行なった結果が以下の通りです．

ランキング上位100件のデータを取得しています．

左側から，
```
実行したベンチマーク名 / 実行した回数 / 1回あたりの実行時間(ns/op) / 1回あたりの確保容量(B/op) / 1回あたりのアロケーション回数(allocs/op) 
```
となっています．

- データ数が1000の場合
```
# Redis
BenchmarkGetRankings-4   	     328	   3576331 ns/op	   41288 B/op	     908 allocs/op

# MySQL
BenchmarkGetRankings-4   	      63	  18586058 ns/op	  392967 B/op	   18051 allocs/op
```

- データ数が5000の場合
```
# Redis
BenchmarkGetRankings-4   	     348	   3679865 ns/op	   41952 B/op	     908 allocs/op

# MySQL
BenchmarkGetRankings-4   	      20	  58493116 ns/op	 2169418 B/op	   90068 allocs/op
```

- データ数が10000の場合
```
# Redis
BenchmarkGetRankings-4   	     301	   3582936 ns/op	   42096 B/op	     908 allocs/op

# MySQL
BenchmarkGetRankings-4   	      12	 105611875 ns/op	 4443574 B/op	  180075 allocs/op
```