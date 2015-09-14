// The service package is a wrapper tool to automaticaly declare and
// register a microservice with vulcand as a backend.

/*
Package service is a wrapper tool to automaticaly declare and
register a microservice with vulcand as a backend.

The wrapper takes standard http handlers and will fall back to the
DefaultServeMux if no Handler is provided

The package also allows you to declare a service as private requiring.
If a package is defined as private it will register a private key for
the service and update its randomly generated security key
*/
package service
