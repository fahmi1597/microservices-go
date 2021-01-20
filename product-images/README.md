# Product Images

## Test Uploading 

### REST approach:
```
curl -svX POST localhost:9001/images/1/test.png --data-binary @test.png
```

### Multipart approach:
```
curl -svX POST localhost:9001/images/2/test.png -F file=@test.png
```

### Get files in normal:
Transfer file in normal size
```
curl -v localhost:9001/images/1/test.png -o normal.png
```

### Get files in gzip compressed:
File is compressed before transfer, automatically decompressed once it arrives
```
curl -v --compressed localhost:9001/images/1/test.png -o zipped.png
```