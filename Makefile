.PHONY: build/image

build/image:
	docker build -t pipeflow:latest .
