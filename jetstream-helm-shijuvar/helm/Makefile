create-namespace-nats:
	kubectl create ns nats 

create-namespace-ordering:
	kubectl create ns ordering 

repo:
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/

dryrun-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=true --dry-run > test-nats.yaml 

dryrun-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222 --dry-run > test-nack.yaml 
	# helm install nack nats/nack  --set=namespace=nats --set=jetstream.nats.url=nats://nats:4222 --dry-run > test-nack.yaml 


install-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=true

install-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222
	# helm install nack nats/nack --set=namespace=nats --set=jetstream.nats.url=nats://nats:4222

uninstall-nats:
	helm uninstall nats 

uninstall-nack:
	helm uninstall nack 
	
uninstall: uninstall-nack uninstall-nats


portforward-nats:
	kubectl port-forward -n nats svc/nats 4222:4222 

# ------------------- order monitoring ------------------------------

dryrun-ordering-monitor:
	helm install ordering-monitor monitor --values monitor/values.yaml --dry-run > test-monitor.yaml  


install-ordering-monitor:
	helm install ordering-monitor monitor --values monitor/values.yaml  

# portforward-ordering-monitor:
# 	kubectl port-forward -n publisher svc/product-write 8080:8080 
 
uninstall-ordering-monitor:
	helm uninstall ordering-monitor
	

# ------------------- order creation ------------------------------


dryrun-ordering-new:
	helm install ordering-new order-new --values order-new/values.yaml --dry-run > test-order-new.yaml  


install-ordering-new:
	helm install ordering-new order-new --values order-new/values.yaml  

# portforward-ordering-new:
# 	kubectl port-forward -n publisher svc/product-write 8080:8080 
 
uninstall-ordering-new:
	helm uninstall ordering-new
	
# ------------------- order review ------------------------------


dryrun-ordering-review:
	helm install ordering-review order-review --values order-review/values.yaml --dry-run > test-order-review.yaml  


install-ordering-review:
	helm install ordering-review order-review --values order-review/values.yaml  

# portforward-ordering-review:
# 	kubectl port-forward -n publisher svc/product-write 8080:8080 
 
uninstall-ordering-review:
	helm uninstall ordering-review
	