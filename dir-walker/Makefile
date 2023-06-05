BINARY_NAME=dir-walker

build:
	go build -o ${BINARY_NAME} dir-walker.go

clean:
	go clean
	rm ${BINARY_NAME}

install: build
	install -D -m 0755 ${BINARY_NAME} $(DESTDIR)/opt/${BINARY_NAME}

run: build
	./${BINARY_NAME}