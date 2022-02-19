# コンパイル用ステージ
FROM golang:1.17.7-alpine as builder
ENV ROOT=/go/src/app
WORKDIR ${ROOT}
RUN apk update && apk --no-cache add git
COPY go.mod go.sum ./
RUN go mod download all
COPY ./cmd ${ROOT}
RUN GOOS=linux GOARCH=amd64 go build -o main {$ROOT}/main.go

# 開発用ステージ
FROM golang:1.17.7-alpine as dev
ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY go.mod go.sum ./
RUN apk --no-cache add git

RUN go mod download all

RUN go get -u github.com/cosmtrek/air \
  && go get -u golang.org/x/tools/cmd/goimports \
  && go get github.com/go-sql-driver/mysql


EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# 本番用ファイナルステージ
# バイナリファイルが実行できればいいので無駄は全て削ぎ落とした
# 今回は不要だが今後のイメージのために作成した
FROM scratch as prod
ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder ${ROOT} .
EXPOSE 8080
CMD ["./main"]