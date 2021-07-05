namespace:
	kubectl create ns nats 
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
	
uninstall: uninstall-nack uninstall-nats uninstall-publisher uninstall-subscriber uninstall-subscriber-update


portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 

#----------------- publisher -----------


namespace-publisher:
	kubectl create ns publisher 

dryrun-publisher:
	helm install publisher publisher --values publisher/values.yaml --dry-run > test-publisher.yaml 

install-publisher:
	helm install publisher publisher --values publisher/values.yaml


uninstall-publisher:
	helm uninstall publisher 

remove-ns-publisher:
	kubectl delete ns publisher 


portforward-publisher:
	kubectl port-forward -n publisher svc/publisher 9090:8080

#----------------- subscriber -----------


namespace-subscriber:
	kubectl create ns subscriber 

dryrun-subscriber:
	helm install subscriber subscriber --values subscriber/values.yaml --dry-run > test-subscriber.yaml 

install-subscriber:
	helm install subscriber subscriber --values subscriber/values.yaml


uninstall-subscriber:
	helm uninstall subscriber 

remove-ns-subscriber:
	kubectl delete ns subscriber 

#----------------- subscriber update -----------
 
dryrun-subscriber-update:
	helm install subscriber-update subscriber --values subscriber/values-update.yaml --dry-run > test-subscriber-update.yaml 

install-subscriber-update:
	helm install subscriber-update subscriber --values subscriber/values-update.yaml


uninstall-subscriber-update:
	helm uninstall subscriber-update 