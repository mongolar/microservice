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

The fields for Service type

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

The service is registered as a backend for Vulcand using this standard

/vulcand/bakends/{Title}.{Version}/backend = json of service
/vulcand/bakends/{Title}.{Version}/servers/{Host}.{Port} = "http://{Host}:{Port}

If service is private

/vulcand/bakends/{Title}.{Version}/privatekey = rotating privatekey
(see private service for more information)

Environment Variables

The below values are required unless set in flags.

ETCD_MACHINES = a pipe delimited list of etcd machines, if the service utilizes
this value it watches the value for changes based on the frequency flag

MICRO_SERVICES_HOST = hostname to declare for microservices on this machine

Flags
-etcd: This flag declares the etcd machines, this overrides the Environment Variable
-frequency:


*/
package service
