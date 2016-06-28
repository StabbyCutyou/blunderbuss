// Package services is to hold the services that wrap business logic and data storage
// concerns.
// Take care not to add circular dependencies between services and other packages.
// Services should make use of the models package, but make sure models do not rely on
// their services.
// Services should each have a NewService style constructor, where all of their dependencies
// are passed in at the time it is created, to help with testing and to maintain a proper
// chain of dependence.
// It is acceptable for services to depend on other services in this manner
package services
