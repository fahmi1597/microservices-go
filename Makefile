get:
	curl -s localhost:3000/products
get_1:
	curl -s localhost:3000/products/1
post:
	curl -s localhost:3000/products -d '{"name":"tea", "description":"A cup of tea", "price":0.25, "sku":"xyz-321-ab"}'
put:
	curl -sX PUT localhost:3000/products -d '{"id":1,"name":"tea", "description":"A cup of tea", "price":0.25, "sku":"xyz-321-ab"}'
