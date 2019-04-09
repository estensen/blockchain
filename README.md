# blockchain
Blockchain storing heart rate

Requires [IPFS](https://docs.ipfs.io/introduction/install/)

## Run
```sh
$ go run main.go
```

### Add block
```sh
$ curl -d '{"BPM": 50}' -H "Content-Type: application/json" -X POST http://localhost:8081
```

### Get blockchain
```sh
$ curl http://localhost:8080
```
 