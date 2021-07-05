namespace:
	kubectl create ns nats 

dryrun-nats:
	helm install nats nats --values nats/values.yaml --dry-run > test-nats.yaml 

dryrun-stan:
	helm install stan stan --values stan/values.yaml --dry-run > test-stan.yaml 

install-nats:
	helm install nats nats --values nats/values.yaml 

install-stan:
	helm install stan stan --values stan/values.yaml

uninstall-nats:
	helm uninstall nats 

uninstall-stan:
	helm uninstall stan 
	
uninstall: uninstall-stan uninstall-nats


portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 
