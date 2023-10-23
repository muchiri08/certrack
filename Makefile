.PHONY: run
run:
	@go run ./cmd/web

.PHONY: watch
watch: 
	@npx tailwindcss -i ./ui/static/main.css -o ./ui/static/output.css --watch