up:
	@docker compose up -d --build

up-api:
	@docker compose up -d --build rinha-api-1 rinha-api-2 rinha-api-3

down:
	docker compose down

down-api:
	docker compose down rinha-api-1 rinha-api-2 rinha-api-3

logs-app-1:
	docker-compose logs -f rinha-api-1