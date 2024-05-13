# 概要
MaxIdleConns/MaxIdleConnsPerHost/MaxConnsPerHost/IdleConnTimeoutおよび `cenkalti/backoff/v4` の挙動を確認するためのサンプル。

# 使い方
`docker-compose up -d` でいけるはず。clientからひたすらserver1とserver2にリクエストします。

# 観察方法
リクエストのログを出しています。各リクエストにはUUIDを振ってあります。
```
docker-compose logs -t -f
```
netstatで状態毎に集計した結果を観察。
```
docker exec -it test-go-nethttp-client-1 sh -c 'while true; do date; netstat -tanp | grep "8080" | sed -E "s/ +/,/g" | cut -d"," -f6,5 | sort | uniq -c; sleep 1; done'
```

# その他
serverは以下のような実装になっています。
- 同時接続数を制限
- 一定の確率で500エラーを返したり、レイテンシを悪化させたり
clientはsysctlsで以下のように設定しています。
- tcp_tw_reuse=1でクライアントからの接続を使い回すように
- tcp_fin_timeout=5でFIN-WAIT-2からCLOSINGへの待機をデフォルトの60sから5sに
- ip_local_port_rangeを極端に狭めることで枯渇した状態を発生しやすく
