# blockchain
Blockchain storing heart rate data on IPFS

_This is a proof-of-concept implementation. It has not been tested for production._

Requires [IPFS](https://docs.ipfs.io/introduction/install/) and [Go](https://golang.org/doc/install)

## Run
Init and start IPFS node
```sh
$ ipfs init
$ ipfs daemon
```

Start blockchain node
```sh
$ go run main.go
```

Or restore blockchain from file
```sh
$ go run main.go --restore test_file
```

### Add block
```sh
$ curl -d '{"BPM": 50}' -H "Content-Type: application/json" -X POST http://localhost:8081
```

Run twice to see that the same data results in different IPFS hashes.

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
    "IPFSHash": "QmTGvLDqmnJH2HLKBvJovxznsv1bpUCCR2vxRwXYFdSo4y",
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

### Get block data
```sh
$ curl http://localhost:8081/block/Qmbntg92Ub7HJfz2xB1X9Zttjd5jcUaigRJCkfq4Wkn3wx
{"BPM":55}
```

### Save blockchain to disk
```sh
$ curl http://localhost:8081/save/test_file
```

## References
This implementation is based on [this](https://medium.com/@mycoralhealth/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc) Medium post by Coral Health.
