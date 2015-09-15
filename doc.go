// The service package is a wrapper tool to automaticaly declare and
// register a microservice with vulcand as a backend.

/*
The service package takes standard http handlers and will fall back to the
DefaultServeMux if no Handler is provided.

The Service details are loaded from a Service.yaml found in the root
folder of the project.  Thse details are used to broadcast details
about the service to etcd in the format expected by vulcand.
You can skip the service declaration and declare it straight in code
as well by creating a new Service and setting the values manually.

Service type

The fields included in the service type.
	Title: Unique title for the service
	Version: version id of the service
	Type: the Type required by Vulcand
	Private: If the service is private
	Requires: a slice of other services this service requires
		the required services only utilize Title, Version, Private
	Parameters: A map of parameters and their expected format*
	Method: The http method this service utilizes*
	Handler: http.Handler (defaults to DefaultServeMux)

*(for future client auto generation)

Registration paths

The service is registered as a backend for Vulcand using this standard.

	/vulcand/bakends/{Title}.{Version}/backend = json of service
	/vulcand/bakends/{Title}.{Version}/servers/{Host}.{Port} = "http://{Host}:{Port}

If service is marked private

	/vulcand/bakends/{Title}.{Version}/privatekey = rotating privatekey
	(see private service for more information)

Environment Variables

The below values are required unless set in flags.

	ETCD_MACHINES = a pipe delimited list of etcd machines, if the service utilizes
	this value it watches the value for changes based on the frequency flag

	MICRO_SERVICES_HOST = hostname to declare for microservices on this machine

Command Line Flags

Values that are passed via command line.

	-etcd: This flag declares the etcd machines, this overrides the ETCD_MACHINES

	-frequency: The frequency at which the service will refresh its ttl for
	Vulcand, and how often ETCD_MACHINES will be checked (default 10)

	-host: host to declare for itself

	-port: the port to serve this service (this field is required, unless hard coded)


Private Service

Private services are those microservices that are only supposed to be served to
internal clients.  If the service being served is marked as private the service
will attempt to lead private key generation.  It will maintain key values for
2x the frequency but rotate keys at the frequency.

The private key will be expected in the header of all internal requests.  All
other requests will be marked as forbidden.



*/
package service
