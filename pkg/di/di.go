// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package di

import (
	"strings"
	"sync"

	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/sassoftware/sas-ggdk/pkg/result"
)

type runState int

const (
	stopped runState = 0
	started runState = 1
)

var lock sync.RWMutex
var state = stopped
var closeableCache closerCacheInterface = new(closerCache)
var creationStrategyCache creationStrategyCacheInterface = new(strategyCache)

// FactoryFn defines the creation function that is registered with the framework
// for creating instances.
type FactoryFn func() result.Result[any]

// StopFn defines the stop function that is returned by Start and
// StartAllowReplaced.
type StopFn func() error

// Get returns the instance of the given id in the type T using the strategy and
// factory registered for the id. If the returned instance implements io.Close it
// will be closed when the StopFn returned from Start or StartAllowReplaced is
// called. The caller must not call Close on the returned instance.
func Get[T any](id string) result.Result[T] {
	lock.RLock()
	defer lock.RUnlock()
	if state == stopped {
		err := errors.New(`the current state is 'stopped' but must be 'started'`)
		return result.Error[T](err)
	}
	strategyMaybe := creationStrategyCache.Get(id)
	strategyResult := result.FromMaybe(strategyMaybe, errors.New(`no factory found for '%s'`, id))
	instanceResult := result.FlatMap(func(strategy creationStrategyInterface) result.Result[any] {
		return strategy.Instance()
	}, strategyResult)
	return result.FlatMap(func(instance any) result.Result[T] {
		typedInstance, ok := instance.(T)
		if !ok {
			err := errors.New(`requested interface not implemented by instance`)
			return result.Error[T](err)
		}
		return result.Ok(typedInstance)
	}, instanceResult)
}

// RegisterLazySingletonFactory registers a FactoryFn for the given id. If a
// FactoryFn is registered for the given id then the new FactoryFn will replace
// it. The decision of whether this is an error or not is deferred to the start.
// Start will return an error if any FactoryFns were replaced.
// StartAllowReplaced will not return an error.
func RegisterLazySingletonFactory(id string, factory FactoryFn) {
	lock.RLock()
	defer lock.RUnlock()
	strategy := &singletonFactoryCreationStrategy{factory: factory}
	closeableCache.Add(strategy)
	creationStrategyCache.Add(id, strategy)
}

// Reset will restore the framework to its uninitialized state. It will to call
// Close on any existing instances that support io.Closer. It will destroy all
// registered FactoryFns. The framework will be in the "stopped" state after this
// function executes. Any errors from Close will be returned but the framework
// will still be reset.
func Reset() error {
	lock.Lock()
	defer lock.Unlock()
	var err error
	if state != stopped {
		err = stopNoLock()
	}
	closeableCache = new(closerCache)
	creationStrategyCache = new(strategyCache)
	return err
}

// Start the framework. If any FactoryFns were registered for ids that already
// had FactoryFns registered then an error is returned. A StopFn is returned
// that callers must call when they are finished with the framework.
func Start() (StopFn, error) {
	lock.Lock()
	defer lock.Unlock()
	if state == started {
		return nil, errors.New(`the current state is 'started' but must be 'stopped'`)
	}
	replacedIDs := creationStrategyCache.Replaced()
	if len(replacedIDs) > 0 {
		replacedS := strings.Join(replacedIDs, ", ")
		return nil, errors.New(`the following ids were replaced: %s`, replacedS)
	}
	state = started
	return createStopFn(), nil
}

// StartAllowReplaced starts the framework. If any FactoryFns were registered for
// ids that already had FactoryFns registered then NO error is returned. A StopFn
// is returned that callers must call when they are finished with the framework.
func StartAllowReplaced() (StopFn, error) {
	lock.Lock()
	defer lock.Unlock()
	if state == started {
		return nil, errors.New(`the current state is 'started' but must be 'stopped'`)
	}
	state = started
	return createStopFn(), nil
}

func createStopFn() StopFn {
	return func() error {
		lock.Lock()
		defer lock.Unlock()
		return stopNoLock()
	}
}

func stopNoLock() error {
	if state == stopped {
		return errors.New(`the current state is 'stopped' but must be 'started'`)
	}
	// Set state to stopped here before calling closeableCache.Close(). The
	// closeableCache will attempt to close everything even if some Close calls
	// fail. A second call to close will be an error anyway.
	state = stopped
	return closeableCache.Close()
}
