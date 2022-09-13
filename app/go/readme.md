## GO implementation

```sh
go run index.go   


docker build . -t goahoj  
docker run -p 8090:8090 goahoj

curl http://localhost:8090/go\?timeout\=800
```