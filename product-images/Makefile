rest:
	curl -svX POST localhost:9001/images/1/test.png --data-binary @test.png
multi:
	curl -svX POST localhost:9001/images/2/test.png -F file=@test.png
normal:
	curl -v localhost:9001/images/1/test.png -o normal.png
gzip:
	curl -v --compressed localhost:9001/images/1/test.png -o zipped.png