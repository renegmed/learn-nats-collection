#------------ nats/nack ---------------

create-namespace-nats:
	kubectl create ns nats 

create-namespace-drones:
	kubectl create ns drones 
 
uninstall-all: uninstall-events uninstall-query uninstall-commands uninstall-mongo uninstall-nack uninstall-nats

remove-namespace:
	kubectl delete ns drones 
	 
repo-nats:
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
 

portforward-nats:
	kubectl port-forward -n nats svc/nats 4222:4222 

#-------------- mongodb ------------------------------

repo-bitnami:
	helm repo add bitnami https://charts.bitnami.com/bitnami

create-namespace-drones:
	kubectl create ns drones 

dryrun-mongo:
	helm install drones-mongodb -f mongo/mongo-values.yaml bitnami/mongodb --dry-run > test-mongo.yaml

install-mongo:
	helm install drones-mongodb -f mongo/mongo-values.yaml bitnami/mongodb 

#
# drones-mongodb.drones.svc.cluster.local:27017
#
# mongo:
# 	kubectl exec -it -n default pod/my-mongodb-c6dc46fbb-sck7l -- bash
# 	# > mongo 
# 	# > 

upgrade-mongo:
	helm upgrade drones-mongodb -f mongo/mongo-values.yaml bitnami/mongodb

portforward-mongo:
	kubectl port-forward svc/drones-mongodb 27017:27017 -n drones

uninstall-mongo:
	helm uninstall drones-mongodb


#---------------- commands -------------------

dryrun-commands: 
	helm install commands commands --values commands/values.yaml --dry-run > test-commands.yaml  

install-commands: 
	helm install commands commands --values commands/values.yaml  

portforward-commands:
	kubectl port-forward -n drones svc/drones-commands 7070:8080 

uninstall-commands: 
	helm uninstall commands 

#---------------- events -------------------

dryrun-events: 
	helm install events events --values events/values.yaml --dry-run > test-events.yaml  

install-events: 
	helm install events events --values events/values.yaml  
 
uninstall-events: 
	helm uninstall events 

#--------------- query ----------
dryrun-query: 
	helm install query query --values query/values.yaml --dry-run > test-query.yaml  

install-query: 
	helm install query query --values query/values.yaml  
 
uninstall-query: 
	helm uninstall query 

#-------------- test --------------


change-position:
	curl -XPOST http://localhost:7070/api/cmds/position -d '{"drone_id":"drone-1", "latitude": 81.231, "longitude": 43.1231, "altitude": 2301.1, "current_speed": 41.3, "heading_cardinal": 1}'

query-position:
	curl http://localhost:7070/api/cmds/query-position/drone-1


add-telemetry:
	curl -XPOST http://localhost:7070/api/cmds/telemetry -d '{"drone_id":"drone-1", "battery": 72, "uptime": 6941, "core_temp": 21 }'

query-telemetry:
	curl http://localhost:7070/api/cmds/query-telemetry/drone-1


add-alert:
	curl -XPOST http://localhost:7070/api/cmds/alertsignal -d '{"drone_id":"drone-1", "fault_code": 12, "description": "all the things are failing"}'

query-alert:
	curl http://localhost:7070/api/cmds/query-alertsignal/drone-1
