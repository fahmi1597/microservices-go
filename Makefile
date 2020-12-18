get:
	curl -s localhost:3000/product
post:
	curl -s localhost:3000/product -d '{"name":"tea", "description":"A cup of tea", "price":0.25, "sku":"xyz-321-ab"}'
put:
	curl -sX PUT localhost:3000/product -d '{"id":1,"name":"tea", "description":"A cup of tea", "price":0.25, "sku":"xyz-321-ab"}'
