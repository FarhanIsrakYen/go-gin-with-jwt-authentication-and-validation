build:
	docker-compose up --build

down:
	docker-compose down --rmi all --volumes --remove-orphans

jwt-secret:
	go run commands/JwtSecretKeyGenerator.go