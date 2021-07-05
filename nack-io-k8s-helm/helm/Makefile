namespace:
	kubectl create ns nats

dryrun-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=false --dry-run > test-nats.yaml 

dryrun-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222 --dry-run > test-nack.yaml 
 

install-nats:
	helm install nats nats --values nats/values.yaml --set nats.image=synadia/nats-server:nightly --set=nats.jetstream.enabled=false

install-nack:
	helm install nack nack --values nack/values.yaml --set=jetstream.nats.url=nats://nats:4222
	 

uninstall-nats:
	helm uninstall nats 

uninstall-nack:
	helm uninstall nack 
	
uninstall: uninstall-nack uninstall-nats  

portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 


# namespace:
# 	kubectl create ns nats 

# dryrun-nats:
# 	helm install nats nats --values nats/values.yaml --dry-run > test-nats.yaml 

# dryrun-stan:
# 	helm install stan stan --values stan/values.yaml --dry-run > test-stan.yaml 

# install-nats:
# 	helm install nats nats --values nats/values.yaml 

# install-stan:
# 	helm install stan stan --values stan/values.yaml

# uninstall-nats:
# 	helm uninstall nats 

# uninstall-stan:
# 	helm uninstall stan 
	
# uninstall: uninstall-stan uninstall-nats


# portforward:
# 	kubectl port-forward -n nats svc/nats 4222:4222 
