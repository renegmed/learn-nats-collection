init-project:
	go mod init product-write 

build:
	go build -race -o write . 

run: build 
	./write -port 8080 -s nats://localhost:4222 -streamname PRODUCT -topic PRODUCT.inserted 


build-docker:
	docker build . --tag renegmedal/cqrs-product-write:v0.0.1

push:
	docker push renegmedal/cqrs-product-write:v0.0.1

write1:
	curl -XPOST localhost:8080/product -d '{"name": "Computer Table M1202", "sku": "CT-134424-4535", "stock_count": 54}'

write2:
	curl -XPOST localhost:8080/product -d '{"name": "Desk chair CH1102-112", "sku": "CH-776563-6678", "stock_count": 210}'

write3:
	curl -XPOST localhost:8080/product -d '{"name": "Table Lamp TL8990-4432", "sku": "TL-332346-7798", "stock_count": 321}'

get-products:
	curl localhost:8081/product 