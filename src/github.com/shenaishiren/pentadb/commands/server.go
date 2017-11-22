// Contains the implementation of server-command of levelDB

/* BSD 3-Clause License

Copyright (c) 2017, Guan Jiawen, Li Lundong
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"os"
	"net"
	"flag"
	"net/http"
	"net/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/shenaishiren/pentadb/opt"
	"github.com/shenaishiren/pentadb/server"
	"github.com/shenaishiren/pentadb/log"
	"fmt"
)

var LOG = log.NewLog(os.Stdout)

var helpPrompt = `Usage: pentadb [--port <port>] [--path <path>] [options]

A PentaDB rpc server, backed by LevelDB

Options:
	--help           		Display this help message and exit
	--port <port>    		The port to listen on (default: 4567)
	--path <path>    		The path to use for the LevelDB store
`

type Server struct {
	Node *server.Node
}

func (s *Server) listen(port string, path string) error {
	s.Node = server.NewNode("127.0.0.1:" + port)
	db, err := leveldb.OpenFile(path, nil)

	if err != nil {
		LOG.Error(err.Error())
		return err
	}
	s.Node.DB = db
	rpc.Register(s.Node)
	rpc.HandleHTTP()      // bind prc to http service

	l, e := net.Listen("tcp", ":" + port)
	if e != nil {
		LOG.Error("listen error:", e)
	}
	LOG.Info("listening at http://0.0.0.0:" + port)
	http.Serve(l, nil)
	return nil
}


func main() {
	var (
		help bool
		port string
		path string
	)
	flag.BoolVar(&help, "h", false, "Display this help message and exit")
	flag.StringVar(&port, "p", "4567", "The port to listen on (default: 4567)")
	flag.StringVar(&path, "a", opt.DeafultPath, "The path to use for the LevelDB store")

	// run
	flag.Parse()

	// help command
	if help {
		fmt.Print(helpPrompt)
	} else {
		server := new(Server)
		server.listen(port, path)
	}
}