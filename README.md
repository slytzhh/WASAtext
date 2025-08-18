# WASAText

This repo contains the project for the WASA (Web and Software Application) course.

## How to run (in development mode)

You can launch the backend only using:

Per avviare il backend esegui questo
```shell 
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

Per avviare il frontend esegui i prossimi due comandi
```shell
sudo ./open-node.sh
# (here you're inside the container)
yarn run dev
```

## Docker
### Backend
Build:
```shell
docker build -t wasa-text-backend:latest -f Dockerfile.backend .
```
Run:
```shell
docker run -it --rm -p 3000:3000 wasa-text-backend:latest
```
### Frontend
Build:
```shell
docker build -t wasa-text-frontend:latest -f Dockerfile.frontend .
```
Run:
```shell
docker run -it --rm -p 5173:80 wasa-text-frontend:latest
```
## License

See [LICENSE](LICENSE).
