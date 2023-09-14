apply:
	kubectl apply -f ./yml

build-image:
	docker build -f Dockerfile -t onehand/webhook:v3 .
