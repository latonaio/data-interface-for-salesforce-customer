# data-interface-for-salesforce-customer

## 概要
data-interface-for-salesforce-customer は、salesforceの顧客オブジェクトに対する各種アクションに必要なデータの整形、および salesforce から受け取った response の MySQL への格納を行うマイクロサービスです。

## 動作環境
data-interface-for-salesforce-customer は、aion-coreのプラットフォーム上での動作を前提としています。  
使用する際は、事前に下記の通りAIONの動作環境を用意してください。   

* OS: Linux OS     
* CPU: ARM/AMD/Intel     
* Kubernetes     
* AION のリソース     

## セットアップ
1. 以下のコマンドを実行して、docker imageを作成してください。
```
$ cd /path/to/data-interface-for-salesforce-customer
$ make docker-build
```

2. 本マイクロサービスは DB に MySQL を使用します。MySQL に関する設定を、 `data-interface-for-salesforce-customer.yaml` の環境変数に記述してください。

| env_name | description |
| --- | --- |
| MYSQL_HOST | ホスト名 |
| MYSQL_PORT | ポート番号 |
| MYSQL_USER | ユーザー名 |
| MYSQL_PASSWORD | パスワード |
| MYSQL_DBNAME | データベース名 |
| MAX_OPEN_CONNECTION | 最大コネクション数 |
| MAX_IDLE_CONNECTION | アイドル状態の最大コネクション数 |
| KANBANADDR: | kanban のアドレス |
| TZ | タイムゾーン |

## 起動方法
以下のコマンドを実行して、podを立ち上げてください。
```
$ cd /path/to/data-interface-for-salesforce-customer
$ kubectl apply -f data-interface-for-salesforce-customer.yaml
```

## kanban との通信
### kanban から受信するデータ
kanban から受信する metadata に下記の情報を含む必要があります。

| key | value |
| --- | --- |
| connection_type | 文字列 "request" |
| method | 使用する http method |
| object | Customer |
| id | 顧客 ID |
| account_id | 顧客 ID |
| body | salesforce api に送信する流クエストの body |

具体例 1: 
```example
# metadata (map[string]interface{}) の中身

"connection_type": "request"
"method": "get"
"object": "Account"
"id": "xxxx"
```

具体例 2: 
```example
# metadata (map[string]interface{}) の中身

"connection_type": "request"
"method": "get"
"object": "Account"
"account_id": "xxxx"
```

### kanban に送信するデータ
kanban に送信する metadata は下記の情報を含みます。

| key | type | description |
| --- | --- | --- |
| method | string | 使用する http method |
| object | string | 文字列 "Account" を指定 |
| path_param | string | 顧客 ID を指定 |
| connection_key | string | 文字列 customer と、使用する http メソッドを(_)で結合した文字列 (例: customer_get)|


具体例: 
```example
# metadata (map[string]interface{}) の中身

"method": "get"
"object": "Account"
"path_param": "xxxx"
"connection_key": "customer_get"
```

## kanban(salesforce-api-kube) から受信するデータ
kanban からの受信可能データは下記の形式です

| key | value |
| --- | --- |
| key | 文字列 "Account" |
| content | Contract の詳細情報を含む JSON 配列 |
| connection_type | 文字列 "response" |

具体例:
```example
# metadata (map[string]interface{}) の中身

"key": "Account"
"content": "[{xxxxxxxxxxx}]"
"connection_type": "response"
```