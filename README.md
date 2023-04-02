## What is a load balancer

A load balancer is a server that has a load balancing function that distributes requests to multiple servers. It is a kind of reverse proxy that increases the availability of services.

There are two main types of load balancers:
- L7: application layer
- L4: transport layer.

In addition to load balancing, the load balancer has the functions of persistence (session maintenance) and health check.

### Type of load balancing
There are different types of load balancing:
- Static:
	- A typical static method is Round Robin, which distributes requests evenly.
- Dynamic:
	- A typical dynamic method is Least Connection, which distributes requests to the server with the least number of unprocessed requests.

### Types of persistence
Persistence is a function for maintaining a session between multiple servers to which the load balancer is distributed.

There are two main methods:
- Source address affinity persistence:
	- Fix the destination server by looking at the source IP address.
- Cookie persistence:
	- Issues a cookie to maintain a session, looks at the cookie, and fixes the server to which it is distributed.

### Health check type
Health check is a function that the load balancer checks the operating status of the distribution destination server.

An active health check method that checks the health of the load balancer to the server to which it is distributed, and a method that monitors the response to requests from clients.

Active checks can be classified into L3 check, L4 check, and L7 check depending on the protocol used.
