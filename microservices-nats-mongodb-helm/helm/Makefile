
repo-mongodb:
	helm repo add bitnami https://charts.bitnami.com/bitnami

create-namespace-nats:
	kubectl create ns nats 

namespace-bookings:
	kubectl create ns bookings 

delete-namespace-bookings:
	kubectl delete ns bookings 
	
uninstall-all:	uninstall-mongodb-showtimes uninstall-mongodb-bookings uninstall-mongodb-movies uninstall-mongodb-website \
	uninstall-website uninstall-showtimes uninstall-bookings uninstall-movies uninstall-nack uninstall-nack

# ----------- nats and nack ------------------

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

portforward-nats:
	kubectl port-forward -n nats svc/nats 4222:4222 

#----- mongodb booking ---------

dryrun-mongodb-bookings:
	# helm install bookings-mongodb bitnami/mongodb --set=global.namespaceOverride=bookings --dry-run > test-mongo-bookings.yaml
	helm install bookings-mongodb bitnami/mongodb -f mongodb/values-bookings.yaml --dry-run > test-mongo-bookings.yaml

install-mongodb-bookings:
	# helm install bookings-mongodb bitnami/mongodb --set=global.namespaceOverride=bookings 
	helm install bookings-mongodb bitnami/mongodb -f mongodb/values-bookings.yaml

upgrade-mongodb-bookings:
	# helm upgrade bookings-mongodb bitnami/mongodb --set=global.namespaceOverride=booking 
	helm upgrade bookings-mongodb bitnami/mongodb -f mongodb/values-bookings.yaml

portforward-mongodb-bookings:
	kubectl port-forward --namespace bookings svc/bookings-mongodb 27017:27017
	
uninstall-mongodb-bookings:
	helm uninstall bookings-mongodb 



#----- booking ---------



dryrun-bookings: 
	helm install bookings bookings -f bookings/values.yaml --dry-run > test-bookings.yaml

install-bookings:
	helm install bookings bookings -f bookings/values.yaml

upgrade-bookings:
	helm upgrade bookings bookings -f bookings/values.yaml

uninstall-bookings:
	helm uninstall bookings 

portforward-bookings:
	kubectl port-forward --namespace bookings svc/bookings 4000:8080

get-bookings:
	curl http://localhost:4000/api/bookings/ -v

post-bookings:
	curl -XPOST http://localhost:4000/api/bookings/ -d '{"userid":"30001","showtimeid":"180034","movies":["Die Hard","Die Harder"]}'



#-------- mongo movies ----------

namespace-movies:
	kubectl create ns movies 

dryrun-mongodb-movies:	 
	helm install movies-mongodb bitnami/mongodb -f mongodb/values-movies.yaml --dry-run > test-mongo-movies.yaml

install-mongodb-movies: 
	helm install movies-mongodb bitnami/mongodb -f mongodb/values-movies.yaml

upgrade-mongodb-movies: 
	helm upgrade movies-mongodb bitnami/mongodb -f mongodb/values-movies.yaml

portforward-mongodb-movies:
	kubectl port-forward --namespace movies svc/movies-mongodb 27018:27017
	
uninstall-mongodb-movies:
	helm uninstall movies-mongodb 

delete-namespace-movies:
	kubectl delete ns movies 

# ------- movies ------------
 

dryrun-movies: 
	helm install movies movies -f movies/values.yaml --dry-run > test-movies.yaml

install-movies:
	helm install movies movies -f movies/values.yaml

upgrade-movies:
	helm upgrade movies movies -f movies/values.yaml

uninstall-movies:
	helm uninstall movies 

portforward-movies:
	kubectl port-forward --namespace movies svc/movies 8282:8080

#-------- mongodb showtimes ----------

namespace-showtimes:
	kubectl create ns showtimes 

dryrun-mongodb-showtimes:	 
	helm install showtimes-mongodb bitnami/mongodb -f mongodb/values-showtimes.yaml --dry-run > test-mongo-showtimes.yaml

install-mongodb-showtimes: 
	helm install showtimes-mongodb bitnami/mongodb -f mongodb/values-showtimes.yaml

upgrade-mongodb-showtimes: 
	helm upgrade showtimes-mongodb bitnami/mongodb -f mongodb/values-showtimes.yaml

portforward-mongodb-showtimes:
	kubectl port-forward --namespace showtimes svc/showtimes-mongodb 27019:27017
	
uninstall-mongodb-showtimes:
	helm uninstall showtimes-mongodb 

delete-namespace-showtimes:
	kubectl delete ns showtimes 

# -------------- showtimes --------------------------

dryrun-showtimes: 
	helm install showtimes showtimes -f showtimes/values.yaml --dry-run > test-showtimes.yaml

install-showtimes:
	helm install showtimes showtimes -f showtimes/values.yaml

upgrade-showtimes:
	helm upgrade showtimes showtimes -f showtimes/values.yaml

uninstall-showtimes:
	helm uninstall showtimes 

portforward-showtimes:
	kubectl port-forward --namespace showtimes svc/showtimes 6000:8080


#-------- mongodb users ----------

namespace-users:
	kubectl create ns users 

dryrun-mongodb-users:	 
	helm install users-mongodb bitnami/mongodb -f mongodb/values-users.yaml --dry-run > test-mongo-users.yaml

install-mongodb-users: 
	helm install users-mongodb bitnami/mongodb -f mongodb/values-users.yaml

upgrade-mongodb-users: 
	helm upgrade users-mongodb bitnami/mongodb -f mongodb/values-users.yaml

portforward-mongodb-users:
	kubectl port-forward --namespace users svc/users-mongodb 27020:27017
	
uninstall-mongodb-users:
	helm uninstall users-mongodb 

delete-namespace-users:
	kubectl delete ns users 


# ----------- users -------------

dryrun-users: 
	helm install users users -f users/values.yaml --dry-run > test-users.yaml

install-users:
	helm install users users -f users/values.yaml

upgrade-users:
	helm upgrade users users -f users/values.yaml

uninstall-users:
	helm uninstall users 

portforward-users:
	kubectl port-forward --namespace users svc/users 8080:8080


#-------- website ----------
 
namespace-website:
	kubectl create ns website

dryrun-website: 
	helm install website website -f website/values.yaml --dry-run > test-website.yaml

install-website:
	helm install website website -f website/values.yaml

upgrade-website:
	helm upgrade website website -f website/values.yaml

uninstall-website:
	helm uninstall website 

delete-namespace-website:
	kubectl delete ns website 

portforward-website:
	kubectl port-forward --namespace website svc/website 8080:8080

# ------------ ingress website -------------

dryrun-ingress-website: 
	helm install ingress-website ingress-website -f ingress-website/values.yaml --dry-run > test-ingress-website.yaml

install-ingress-website:
	helm install ingress-website ingress-website -f ingress-website/values.yaml

upgrade-ingress-website:
	helm upgrade ingress-website ingress-website -f ingress-website/values.yaml

uninstall-ingress-website:
	helm uninstall ingress-website 

# ------------ ingress users -------------


dryrun-ingress-users: 
	helm install ingress-users ingress-users -f ingress-users/values.yaml --dry-run > test-ingress-users.yaml

install-ingress-users:
	helm install ingress-users ingress-users -f ingress-users/values.yaml

upgrade-ingress-users:
	helm upgrade ingress-users ingress-users -f ingress-users/values.yaml

uninstall-ingress-users:
	helm uninstall ingress-users


# ------------ test ------------

list-users: 
	curl http://192.168.49.2/api/users/ -v

post-user1:
	curl -XPOST http://192.168.49.2/api/users/ -d '{"name":"John","lastname":"Doe"}' -v 

post-user2:
	curl -XPOST http://192.168.49.2/api/users/ -d '{"name":"Jane","lastname":"Denali"}'  

post-user3:
	curl -XPOST http://192.168.49.2/api/users/ -d '{"name":"Peter","lastname":"Ducker"}'  