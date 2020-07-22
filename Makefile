PB_FILES=$(shell find ./assets/protos/ -maxdepth 1 -mindepth 1 -type f -name *.proto -exec basename {} \;)

gen-pb:
	protoc --proto_path=assets/protos --go_out=internal/pkg/message/pb --go_opt=paths=source_relative $(PB_FILES)
