<!-- vim: set fileencoding=utf-8 : -->

# ビルド時の位置とコマンド
リポジトリルートからのビルドである前提だと`Dockerfile`に記載をしているので、
その前提で移動してからビルドする。
あとは、GCRに登録できる名前でイメージ作成をしている。

```
cd /path/to/grpc-microservice-sample
docker build --file docker/Dockerfile -t asia.gcr.io/kubernetes-229910/grpc-user-service:1.0.0 ./
```

## やっておきたいこと
サーバビルドとかコンテナビルドプッシュを一括で管理するMakefileをルートリポジトリの直下に作るかどうか

