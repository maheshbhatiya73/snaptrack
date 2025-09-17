dev:
	cd web && npm run dev & \
	cd server && air

build:
	cd web && npm run build

serve:
	go run server/main.go
