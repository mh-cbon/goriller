#!/bin/sh

set -ex

rm `which goriller`
go install


goriller - demo/Controller:ControllerGoriller | grep -F 'router.HandleFunc("/{id}", t.embed.GetByID).Methods("GET")' || exit 1;
goriller -mode rpc - demo/Controller:ControllerGorillerRPC | grep -F 'router.HandleFunc("GetByID", t.embed.GetByID)' || exit 1;
goriller - demo/Controller:ControllerGoriller | grep "package main" || exit 1;
goriller -p nop - demo/Controller:ControllerGoriller | grep "package nop" || exit 1;

goriller - demo/Controller:ControllerGoriller | grep "embed Controller" || exit 1;
goriller - demo/*Controller:ControllerGoriller | grep -F "embed *Controller" || exit 1;

goriller - demo/Controller:*ControllerGoriller | grep "embed Controller" || exit 1;
goriller - demo/*Controller:*ControllerGoriller | grep -F "embed *Controller" || exit 1;

rm -fr gen_test
goriller demo/Controller:gen_test/ControllerGoriller || exit 1;
ls -al gen_test | grep "controllergoriller.go" || exit 1;
cat gen_test/controllergoriller.go | grep -F 'router.HandleFunc("/{id}", t.embed.GetByID).Methods("GET")' || exit 1;
cat gen_test/controllergoriller.go | grep "package gen_test" || exit 1;
rm -fr gen_test

rm -fr demo/*gen.go
go generate demo/main.go
ls -al demo | grep "controllergoriller.go" || exit 1;
cat demo/controllergoriller.go | grep "package main" || exit 1;
cat demo/controllergoriller.go | grep "NewControllerGoriller(" || exit 1;
go run demo/*.go | grep "Red" || exit 1;

rm -fr demo/*gen.go
go generate github.com/mh-cbon/goriller/demo
ls -al demo | grep "controllergoriller.go" || exit 1;
cat demo/controllergoriller.go | grep "package main" || exit 1;
cat demo/controllergoriller.go | grep "NewControllerGoriller(" || exit 1;
go run demo/*.go | grep "Red" || exit 1;
# rm -fr demo/gen # keep it for demo

# go test


echo ""
echo "ALL GOOD!"
