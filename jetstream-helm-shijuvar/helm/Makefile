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
	
uninstall: uninstall-nack uninstall-nats


portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 
