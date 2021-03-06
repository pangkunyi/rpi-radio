export GOPATH=$(shell pwd)
install:
	@go install rpi-radio
run:
	@pkill rpi-radio || echo "no rpi-radio process"
	@mkdir -p ~/.logs/rpi-radio/logs
	@nohup ./bin/rpi-radio 2>&1 >> ~/.logs/rpi-radio/logs/rpi-radio.log &
stop:
	@pkill rpi-radio || echo "no rpi-radio process"
dep:
	@go get -u github.com/pangkunyi/baidu-pcs
	@go get -u github.com/gorilla/mux
