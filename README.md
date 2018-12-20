# Sakudog in Golang
さくらのクラウドからAWSへのトラフィックをメトリクスとしてDatadogに送信しDatadog上で監視できるようにするスクリプトをGolangで書いたものです。

## 使用するもの
* Golang
* Datadog

## セットアップ
1.  リポジトリをclone
```
git clone git@github.com:moneyforward/sakudog.git
```

2.  go-datadog-apiをインストール
```
go get gopkg.in/zorkian/go-datadog-api.v2
```

## 使い方

1.  それぞれ自分のクレデンシャルに置き換える

**main.go**に定義されているクレデンシャルの部分を自分のものと置き換えてください。
```go
// Sakura cloud
const (
	// UserID
	token = "YOUR_API_KEY"
	// Password
	secret = "YOUR_SECRET_TOKEN"
	// BaseURL 
	url = "https://secure.sakura.ad.jp/cloud/zone/${ZONE_NAME}/api/cloud/1.1//commonserviceitem/${COMMONSERVICEITEMID}/activity/awsdirectconnect/monitor"
)

// Datadog
const (
	apiKey   = "YOUR_API_KEY"
	appKey   = "YOUR_APP_KEY"
	screenId = "YOUR_SCREEN_ID"
)
```

ドキュメント参照: https://developer.sakura.ad.jp/cloud/api/1.1/

2.  メトリクスを設定する
```go
    receive := datadog.Metric{
			Metric: datadog.String("sakudog.dx.receive_bytes_per_s"),
			Type:   datadog.String("gauge"),
			Points: []datadog.DataPoint{
				// TODO:-convert custom type(val) to float64
				{ConvertStingToFloat64(key), ConvertInt64ToFloat64(val.ReceiveBytesPerSec)},
			},
		}

		send := datadog.Metric{
			Metric: datadog.String("sakudog.dx.send_bytes_per_s"),
			Type:   datadog.String("gauge"),
			Points: []datadog.DataPoint{
				// TODO:-convert custom type(val) to float64
				{ConvertStingToFloat64(key), ConvertInt64ToFloat64(val.SendBytesPerSec)},
			},
		}

```

そのままでも使えますが、メトリクスの各種パラメーターの値は必要に応じて追加・変更してください。

ドキュメント参照: https://godoc.org/gopkg.in/zorkian/go-datadog-api.v2#Metric

3. スクリプトを実行

cmdディレクトリに移動して、

```
go run main.go
```

でスクリプトを実行します。

Datadogのサイドメニューから、Metrics→Explorerを選んでGraphの覧にメトリクス名を入力して表示されれば成功です。


あとはこのスクリプトが5分毎に定期実行される環境を用意すればDatadog上で常に最新のメトリクスを監視できます。


AWSのLambdaとCloudWatchを使えば無料の枠で5分間隔の定期実行ができるのでオススメです。


## ライセンス
This software is released under the MIT License, see LICENSE.


## 参照
https://github.com/zorkian/go-datadog-api
