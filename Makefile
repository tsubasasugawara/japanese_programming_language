build:
	docker build -t golang-jpl:alpine .

run:
	docker run -it -v $(PWD)/src:/go/src -w /go/src --name golang-jpl golang-jpl:alpine

start:
	clear
	docker start -ia golang-jpl
