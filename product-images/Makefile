rest:
	curl -svX POST localhost:8080/images/1/test.png --data-binary @test.png
multi:
	curl -svX POST localhost:8080/images/1/test.png -F file=@test.png
gzip:
	curl -v --compressed localhost:8080/images/1/test.png -o zipped.png
normal:
	curl -v localhost:8080/images/1/test.png -o normal.png