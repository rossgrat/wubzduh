IMAGE = ghcr.io/rossgrat/wubzduh
SERVER = potatoserver
REMOTE_DIR = ~/services/wubzduh

.PHONY: deploy stop logs

deploy:
	scp deploy/docker-compose.yml $(SERVER):$(REMOTE_DIR)/docker-compose.yml
	scp config.yaml $(SERVER):$(REMOTE_DIR)/config.yaml
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker pull $(IMAGE):latest && docker compose up -d"

stop:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose down"

logs:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose logs -f"
