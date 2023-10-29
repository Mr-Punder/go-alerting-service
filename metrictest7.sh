#!/bin/bash
cd cmd/agent
go build -o agent *.go

cd ../server
go build -o server *.go
 cd ../..


            SERVER_PORT=8099
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration7$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.