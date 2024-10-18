all:
	go build -o main
	chmod +x main
	./main --config .layerg-crawler.yaml

api:
	go build -o main
	chmod +x main
	./main api --config .layerg-crawler.yaml