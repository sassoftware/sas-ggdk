// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package di

import (
	"io"
	"sync"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// Ensure singletonFactoryCreationStrategy implements creationStrategyInterface.
var _ creationStrategyInterface = (*singletonFactoryCreationStrategy)(nil)

type creationStrategyInterface interface {
	io.Closer
	Instance() result.Result[any]
}

type singletonFactoryCreationStrategy struct {
	lock     sync.RWMutex
	factory  FactoryFn
	instance any
}

// Close the creation strategy. When the creation strategy is closed any instance
// previously created is closed (if the instance implements io.Closer). Any
// close failures are returned.
func (s *singletonFactoryCreationStrategy) Close() error {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.instance == nil {
		return nil
	}
	closer, ok := s.instance.(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}

// Instance creates an instance from the current FactoryFn if one does not exist
// and saves that instance. If an instance already exists, that existing instance
// is returned.
func (s *singletonFactoryCreationStrategy) Instance() result.Result[any] {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.instance != nil {
		return result.Ok[any](s.instance)
	}
	instanceResult := s.create()
	return result.MapNoError(func(i any) any {
		s.instance = i
		return i
	}, instanceResult)
}

func (s *singletonFactoryCreationStrategy) create() result.Result[any] {
	return s.factory()
}
