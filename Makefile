.PHONY: postgres

up:
	docker-compose up -d

up-db:
	docker-compose up -d postgres


build:
	docker build -t ai-finance:local .