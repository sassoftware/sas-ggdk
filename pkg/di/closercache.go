// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package di

import (
	"io"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Ensure closerCache implements closerCacheInterface.
var _ closerCacheInterface = (*closerCache)(nil)

type closerCacheInterface interface {
	io.Closer
	Add(c io.Closer)
}

type closerCache struct {
	lock    sync.RWMutex
	closers []io.Closer
}

// Add the given io.Closer to the cache.
func (c *closerCache) Add(closer io.Closer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closers == nil {
		c.closers = make([]io.Closer, 0, 1)
	}
	c.closers = append(c.closers, closer)
}

// Close the cache. When the cache is closed everything in the cache is closed;
// close failures are returned.
func (c *closerCache) Close() error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var allErrors error
	for _, closer := range c.closers {
		err := closer.Close()
		if err != nil {
			allErrors = multierror.Append(allErrors, err)
		}
	}
	return allErrors
}
