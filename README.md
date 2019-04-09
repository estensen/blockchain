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
$ curl http://localhost:8081
[{"Index":0,"Timestamp":"2019-04-09 16:35:02.643395 +0200 CEST m=+0.003511293","IPFSHash":"","PrevHash":"","Hash":""},{"Index":1,"Timestamp":"2019-04-09 16:35:04.783103 +0200 CEST m=+2.143217249","IPFSHash":"Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx","PrevHash":"","Hash":"2118a9b15ebdcaf4b389063bbb602613ad40f6b6e03c07cfc0c701080b9c9a91"}]
```

### Get block data
```sh
$ curl http://localhost:8081/block/Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx
{"BPM":55}
```
