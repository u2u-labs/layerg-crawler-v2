all:
	go build crawler/main.go
	chmod +x main
	./main --config .layerg-crawler.yaml

api:
	go build crawler/main.go
	chmod +x main
	./main --config .layerg-crawler.yaml api