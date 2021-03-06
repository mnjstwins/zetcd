// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !zkdocker

package integration

import (
	"net"
	"testing"

	"github.com/coreos/zetcd"

	"github.com/coreos/etcd/integration"
)

type zkCluster struct {
	zkClientAddr string

	// don't use these in tests since they may change with build tags

	etcdClus *integration.ClusterV3
	cancel   func()
	donec    <-chan struct{}
}

func newZKCluster(t *testing.T) *zkCluster {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	donec := make(chan struct{})

	// TODO use unix socket
	ln, err := net.Listen("tcp", ":30000")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer close(donec)
		c := clus.RandClient()
		zetcd.Serve(c.Ctx(), ln, zetcd.NewAuth(c), zetcd.NewZK(c))
	}()
	return &zkCluster{
		zkClientAddr: "127.0.0.1:30000",

		etcdClus: clus,
		cancel:   func() { ln.Close() },
		donec:    donec,
	}
}

func (zkclus *zkCluster) Close(t *testing.T) {
	zkclus.etcdClus.Terminate(t)
	zkclus.cancel()
	<-zkclus.donec
}
