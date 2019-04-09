# blockchain
Blockchain storing heart rate data on IPFS

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
[
  {
    "Index": 0,
    "Timestamp": "2019-04-09 17:16:17.334761 +0200 CEST m=+0.003436740",
    "IPFSHash": "",
    "PrevHash": "",
    "Hash": ""
  },
  {
    "Index": 1,
    "Timestamp": "2019-04-09 17:16:19.93325 +0200 CEST m=+2.601923666",
    "IPFSHash": "Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx",
    "PrevHash": "",
    "Hash": "b2bf56b03c2ab09a0bb510f3c0dbdbd3309dc40263e983e48b4c48321891ca32"
  },
  {
    "Index": 2,
    "Timestamp": "2019-04-09 17:16:22.372999 +0200 CEST m=+5.041669596",
    "IPFSHash": "Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx",
    "PrevHash": "b2bf56b03c2ab09a0bb510f3c0dbdbd3309dc40263e983e48b4c48321891ca32",
    "Hash": "13f54f8268937bbe4045650f16d7240f0b69ca76ecc47438bf57a660c34fa5f4"
  }
]
```

Because the data is still unencrypted you can tell that the last two blocks contain the same BPM

### Get block data
```sh
$ curl http://localhost:8081/block/Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx
{"BPM":55}
```
