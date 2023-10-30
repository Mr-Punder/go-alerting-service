#!/bin/bash
cd cmd/agent
go build -o agent *.go

cd ../server
go build -o server *.go
cd ../..

go test ./...



metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server

metricstest -test.v -test.run=^TestIteration2[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

metricstest -test.v -test.run=^TestIteration3[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server

SERVER_PORT=8095
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(mktemp)
          metricstest -test.v -test.run=^TestIteration4$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=8095\
            -source-path=.



SERVER_PORT=8099
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration5$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration6$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.




            SERVER_PORT=8099
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration7$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


            SERVER_PORT=8099
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration8$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.



            SERVER_PORT=8099
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=$(random tempfile)
          metricstest -test.v -test.run=^TestIteration9$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -file-storage-path=./data/db.json \
            -server-port=$SERVER_PORT \
            -source-path=.