create-namespace-nats:
	kubectl create ns nats 

create-namespace-publisher:
	kubectl create ns publisher 

create-namespace-subscriber:
	kubectl create ns subscriber
repo:
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/

dryrun-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=true --dry-run > test-nats.yaml 


install-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=true


dryrun-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222 --dry-run > test-nack.yaml  


install-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222 

uninstall-nats:
	helm uninstall nats 

uninstall-nack:
	helm uninstall nack 
	
uninstall: uninstall-nack uninstall-nats uninstall-cqrs-prodread uninstall-cqrs-prodwrite 
 

portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 


#---------------- product read -------------------

dryrun-cqrs-prodread: 
	helm install prodread cqrs-read --values cqrs-read/values.yaml --dry-run > test-cqrs-read.yaml  

install-cqrs-prodread: 
	helm install prodread cqrs-read --values cqrs-read/values.yaml 

portforward-cqrs-prodread:
	kubectl port-forward -n subscriber svc/product-read 8081:8080 

uninstall-cqrs-prodread: 
	helm uninstall prodread 

#---------------- product write -------------------

dryrun-cqrs-prodwrite:
	helm install prodwrite cqrs-write --values cqrs-write/values.yaml --dry-run > test-cqrs-write.yaml  


install-cqrs-prodwrite:
	helm install prodwrite cqrs-write --values cqrs-write/values.yaml  

portforward-cqrs-prodwrite:
	kubectl port-forward -n publisher svc/product-write 8080:8080 
 
uninstall-cqrs-prodwrite:
	helm uninstall prodwrite
	
#--------------- test ----------


write1:
	curl -XPOST localhost:8080/product -d '{"name": "Computer Table M1202", "sku": "CT-134424-4535", "stock_count": 54}'

write2:
	curl -XPOST localhost:8080/product -d '{"name": "Desk chair CH1102-112", "sku": "CH-776563-6678", "stock_count": 210}'

write3:
	curl -XPOST localhost:8080/product -d '{"name": "Table Lamp TL8990-4432", "sku": "TL-332346-7798", "stock_count": 321}'

get-products:
	curl localhost:8081/product 

