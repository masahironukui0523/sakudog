# Sakudog in Golang
さくらのクラウドからAWSへのトラフィックをメトリクスとしてDatadogに送信しDatadog上で監視できるようにするスクリプトをGolangで書いたものです。

## Dependency
* Golang
* Datadog
* Lambda(定期実行できるものであれば何でもいい) 

## Setup
1.  リポジトリをclone
```
cd $GOPATH/src/
git clone git@github.com:moneyforward/sakudog.git
cd sakudog
```

2.  go-datadog-apiをインストール
```
go get gopkg.in/zorkian/go-datadog-api.v2
```

3.  それぞれ自分のクレデンシャルに置き換える
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

5.  メトリクスを設定する
```go
    receive := datadog.Metric{
			// メトリクス名
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

## Usage

さくらのクラウドのAPIから得られるメトリクスは5分間隔なので、Lambdaにmain.goを5分間おきにスクリプトを実行するようにします。



## Licence
This software is released under the MIT License, see LICENSE.


## Authors
@masahironukui0523(https://github.com/moneyforward/sakudog)


## References
https://github.com/zorkian/go-datadog-api
