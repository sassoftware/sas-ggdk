// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package di

import (
	"sync"

	"github.com/sassoftware/sas-ggdk/pkg/maputils"
	"github.com/sassoftware/sas-ggdk/pkg/maybe"
)

// Ensure strategyCache implements creationStrategyCacheInterface.
var _ creationStrategyCacheInterface = (*strategyCache)(nil)

type creationStrategyCacheInterface interface {
	Add(id string, strategy creationStrategyInterface)
	Get(id string) maybe.Maybe[creationStrategyInterface]
	Replaced() []string
}

type strategyCache struct {
	lock     sync.RWMutex
	replaced []string
	cache    map[string]creationStrategyInterface
}

// Add saves the given creation strategy for the given id.
func (c *strategyCache) Add(id string, strategy creationStrategyInterface) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]creationStrategyInterface)
	}
	_, present := c.cache[id]
	if present {
		c.replaced = append(c.replaced, id)
	}
	c.cache[id] = strategy
}

// Get returns the given creation strategy for the given id.
func (c *strategyCache) Get(id string) maybe.Maybe[creationStrategyInterface] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.cache == nil {
		c.cache = make(map[string]creationStrategyInterface)
	}
	return maputils.Get(c.cache, id)
}

// Replaced returns the list of ids that have been used more than once in calls
// to Add. May be empty but will not be nil.
func (c *strategyCache) Replaced() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	size := len(c.replaced)
	result := make([]string, size)
	copy(result, c.replaced)
	return result
}
