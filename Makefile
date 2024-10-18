all:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler --config .layerg-crawler.yaml

api:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler api --config .layerg-crawler.yaml