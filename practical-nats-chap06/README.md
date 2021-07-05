
## Learning NATS with Cluster ##

This work is based on W. Quevedo's Practical NATS Rides Management (chap 6)
that riders can request a ride. The rides managers will find an 
available driver with specific type of vehicle on behalf of the requester. 

Note that each application can be made to be deployed independently. Common
functions should be moved to a separate github for common access.

To run these applications, do the following:

1. Create and run containers for nats cluster, 

        make up
2. Request a driver with 'regular' type of vehicle. This should be available

        make request-regular

3. Request a driver with SUV type of vehicle. This should be available

        make request-suv

4. Request a driver with bus type of vehicle. This should not be available

        make request-bus
