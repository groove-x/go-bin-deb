
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_DIR := $(patsubst %/,%,$(dir $(mkfile_path)))

ARCH?=amd64

BIN?=${PROJECT_DIR}/go-bin-deb
TEST_DEB?=hello-${ARCH}.deb
TEST_DEB_FILE?=demo/${TEST_DEB}

echo:
	@echo "PROJECT_DIR: ${PROJECT_DIR}"
	@echo "BIN: ${BIN}"
	@echo "TEST_DEB: ${TEST_DEB}"

build: ${BIN}
${BIN}:
	@echo 
	@echo "go-bin-debをビルド"
	@echo 
	go build

test: build
	make ${TEST_DEB_FILE}

	@echo 
	@echo "debファイルに期待する中身があることを確認"
	@echo 
	@echo "expect: symlinkが2件含まれている"
	@dpkg -c demo/hello-amd64.deb | grep "asset -> /usr/share/hello/other/asset1"

${TEST_DEB_FILE}:
	@echo 
	@echo "demoディレクトリの内容をパッケージにできるかテスト"
	@echo 
	cd demo && \
		GOOS=linux GOARCH=${ARCH} go build -o build/${ARCH}/hello hello.go && \
		${BIN} generate -a ${ARCH} --version 0.0.1 -w pkg-build/${ARCH}/ -o ${TEST_DEB}

clean:
	rm -rf ${BIN}
	rm -rf demo/pkg-build
	rm -rf demo/build
	rm -rf demo/*.deb
