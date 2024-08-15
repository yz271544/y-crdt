

.PHONY: build
build:
	@cp tests-ffi/include/libyrs.h yffigo/yrs/include/
	@cp target/debug/libyrs.a yffigo/yrs/lib/
	go build -o yffigo/yffigo yffigo/main.go
