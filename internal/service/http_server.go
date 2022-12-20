// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// httpServer wraper for gin.engine and http.Server
type httpServer struct {
	sync.RWMutex

	e            *gin.Engine
	server       *http.Server
	serverStatus uint8
}

func (s *httpServer) status() uint8 {
	s.RLock()
	defer s.RUnlock()

	return s.serverStatus
}

func (s *httpServer) start() error {
	s.Lock()
	if s.serverStatus == _statusServerStarted || s.serverStatus == _statusServerStoped {
		return nil
	}
	oldStatus := s.serverStatus
	s.serverStatus = _statusServerStarted
	s.Unlock()

	if err := s.server.ListenAndServe(); err != nil {
		s.Lock()
		s.serverStatus = oldStatus
		s.Unlock()

		return err
	}
	return nil
}

func (s *httpServer) stop() error {
	s.Lock()
	defer s.Unlock()

	if s.serverStatus == _statusServerStoped || s.serverStatus == _statusServerInitilized {
		return nil
	}
	if err := s.server.Shutdown(context.Background()); err != nil {
		return err
	}
	s.serverStatus = _statusServerStoped
	return nil
}
