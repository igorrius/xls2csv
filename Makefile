RELEASE_FOLDER_NAME=release
EXECUTABLE_NAME=xls2csv
INSTALL_PATH=/usr/local/bin

all: test

test:
		go test ./...

clean:
	rm -rf ${RELEASE_FOLDER_NAME}

install: build_linux_64
		cp ${RELEASE_FOLDER_NAME}/${EXECUTABLE_NAME} ${INSTALL_PATH}/
		make clean

build: build_linux_64 build_windows_64

release: release_linux_64 release_windows_64

build_windows_64:
		mkdir -p ${RELEASE_FOLDER_NAME}
		GOOS=windows GOARCH=amd64 go build -o ${RELEASE_FOLDER_NAME}/${EXECUTABLE_NAME}.exe

build_linux_64:
		mkdir -p ${RELEASE_FOLDER_NAME}
		GOOS=linux GOARCH=amd64 go build -o ${RELEASE_FOLDER_NAME}/${EXECUTABLE_NAME}

release_windows_64: build_windows_64
		cd ${RELEASE_FOLDER_NAME}; \
		zip -9 -m windows-amd64.zip ${EXECUTABLE_NAME}.exe

release_linux_64: build_linux_64
		cd ${RELEASE_FOLDER_NAME}; \
		zip -9 -m linux-amd64.zip ${EXECUTABLE_NAME}
