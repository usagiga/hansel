# tetona
「画像も貼らずにスレ立てとな!?」 *ー[ガゾーモ・ハラズニ・スレタ＝**テトナ**](https://twitter.com/2chijin/status/552787150127656961)　(１７０１～１９９６)ー*

## 概要
* DiscordからMinecraftゲームサーバ(EC2)を操作するDiscordBot
* Minecraftゲームサーバ: [sureta](https://github.com/shokkunrf/sureta)

## 機能
* EC2の起動(トリガーメッセージ: start)
* EC2の停止(トリガーメッセージ: sleep)
* EC2のPublicIPAddress確認(トリガーメッセージ: status)

## 使い方
### 必要なもの
* docker-compose

### 事前準備
* EC2インスタンスの作成
* DiscordBotの作成

### pull
```sh
$ git pull https://github.com/shokkunrf/tetona.git
$ cd tetona
```

### .envの作成
```
BOT_ID=<DiscordBotのTOKEN>
INSTANCE_ID=<EC2インスタンスid>
AWS_ACCESS_KEY_ID=<EC2にアクセス可能なcredentialsのkey>
AWS_SECRET_ACCESS_KEY=<EC2にアクセス可能なcredentialsのsecret_key>
AWS_DEFAULT_REGION=<EC2インスタンスのリージョン>
AWS_DEFAULT_OUTPUT=json
```

### 実行
```sh
$ docker-compose up --build
```