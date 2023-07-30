build:
	sudo docker stop $$(sudo docker ps | grep ":27017" | awk '{print $$1}')
	go build ./main.go
	sudo docker run -d -p 27017:27017 -v lenscape:/data/db mongo:latest

setup:
	$(MAKE) -i build
	./main

deploy:
	go build -o main
	scp main ubuntu@13.41.110.229:/home/ubuntu

run:
	./main
