namespace:
	kubectl create ns nats

namespace-trans-service:
	kubectl create ns trans-service 

uninstall: uninstall-nack uninstall-nats  uninstall-api uninstall-agent uninstall-manager 

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
	


portforward:
	kubectl port-forward -n nats svc/nats 4222:4222 

#---------------- api-server ----------------------

dryrun-api:
	helm install api-server api-server --values api-server/values.yaml --dry-run > test-api-server.yaml 

install-api:
	helm install api-server api-server --values api-server/values.yaml

uninstall-api:
	helm uninstall api-server 

portforward-api:
	kubectl port-forward -n trans-service svc/api-server 8282:8080 

#---------------- driver-agent ----------------------

dryrun-agent:
	helm install driver-agent driver-agent --values driver-agent/values.yaml --dry-run > test-driver-agent.yaml 

install-agent:
	helm install driver-agent driver-agent --values driver-agent/values.yaml

uninstall-agent:
	helm uninstall driver-agent 

#---------------- rides-manager ----------------------

dryrun-manager:
	helm install rides-manager rides-manager --values rides-manager/values.yaml --dry-run > test-rides-manager.yaml 

install-manager:
	helm install rides-manager rides-manager --values rides-manager/values.yaml

uninstall-manager:
	helm uninstall rides-manager 


# ------------------ test ----------------------


request-suv:
	curl -X POST http://localhost:8282/rides  -d '{"type":"suv", "location": {"lat": 12342243.344, "lng": 345534555.458}}'
 

request-regular:
	curl -X POST http://localhost:8282/rides  -d '{"type":"regular", "location": {"lat": 12342243.344, "lng": 345534555.458}}'
 

request-bus:
	curl -X POST http://localhost:8282/rides  -d '{"type":"bus", "location": {"lat": 12342243.344, "lng": 345534555.458}}'
 