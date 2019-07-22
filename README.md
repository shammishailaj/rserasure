# rserasure
Sample usage of Reed-Solomon Erasure Coding

**Please note:** *This codebase has been tested using version:* ```go version go1.12.5 darwin/amd64```

By default, the server runs on *port 3000* of your *localhost*.

Following is an example of the usage

Assumptions made:
- You have a file named ```rsetry.mp4``` at location ```/tmp/rsetry.mp4```
- You have a directory ```/tmp/rse``` which is writeable by the application
- You have a directory ```/tmp/rse/merged``` which is writeable by the application

Encoder
```
curl -X POST "http://localhost:3000/encode?source=/tmp/rsetry.mp4&targetdir=/tmp/rse&datashards=10&parityshards=5"
```

Decoder
```
curl -X POST "http://localhost:3000/decode?source=/tmp/rse/rsetry.mp4&targetdir=/tmp/rse/merged/rsetry.mp4&datashards=10&parityshards=5"
```
