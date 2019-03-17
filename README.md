# blockchain
Blockchain storing heart rate

## Run
```sh
$ go run main.go
```

### Add block
```sh
$ curl -d '{"BPM": 50}' -H "Content-Type: application/json" -X POST http://localhost:8080
```

### Get blockchain
```sh
$ curl http://localhost:8080
```
 