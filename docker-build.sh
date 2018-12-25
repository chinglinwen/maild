CGO_ENABLED=0 go build
docker build -t maild .
docker tag maild harbor.haodai.net/ops/maild:v1
docker push harbor.haodai.net/ops/maild:v1
