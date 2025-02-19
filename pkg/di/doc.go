// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package di provides a very small dependency injection framework for Go. It
// supports only lazy initialization of singleton objects. These can be thought
// of as the "services" provided to the rest of the program.
//
// Packages that wish to provide a service, register a factory function with a
// string id in their package init. This function will be called the first time a
// request is made to the framework for a service by that id. If the service
// implements io.Closer, then the framework will call Close() when the framework
// is shut down.
//
// Example of a package exposing a service:
//
//	package service
//	import "github.com/sassoftware/sas-ggdk/di"
//	type aService struct {}
//	func(a *aService) Serve() error {
//	    // implements api.IServe
//	    return nil
//	}
//	func init() {
//	    di.RegisterLazySingletonFactory("service", func() result.Result[any] {
//	        return result.Ok[any](&aService{})
//	    })
//	}
//
// The factory function can call Get to get any other services it needs to
// complete construction. If a cycle is encountered a fatal deadlock error will
// be thrown and the process will stop. It is expected that the factory function
// will return pointers to unexported structs that implement exported interfaces.
// If your service does not have any methods or has many data fields that callers
// are expected to use, that is a good signal that you should not be publishing
// that service using this framework.
//
// Packages that wish to make use of the framework must call Start or
// StartAllowReplaced. Start will return an error if there have been any
// factories registered for names that already existed. StartAllowReplaced will
// not. Both of these functions return a stop function that must be called when
// the caller is finished with the framework. Callers ask for a service by id and
// specify the interface they expect that service to implement. If the id does
// not exist or the service registered to that id does not implement the
// requested interface then an error is returned. This framework only supports
// lazily initialized singletons so multiple calls to Get for the same id will
// get the same service.
//
// Example of a package using a service:
//
//	package client
//	import (
//	    "github.com/sassoftware/sas-ggdk/di"
//	    "myproject/api"
//	)
//	func Perform() error {
//	    stopfn, err := di.Start()
//	    defer stopfn()
//	    service, err := di.Get[api.IServe]("service")
//	    if err != nil {
//	        return err
//	    }
//	    return service.Serve()
//	}
package di
