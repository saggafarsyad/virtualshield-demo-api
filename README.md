# VirtualShield: Internet Shield Presentation demo

## Compile
`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go`

## Dockerfile
```
FROM scratch
COPY api /app
ENV VIRTUALSHIELD_DATASOURCE <USERNAME>:<PASSWORD>(<HOST>:<PORT>)/<DATABASE_NAME>?parseTime=true
CMD ["/app"]
```

## Build Container
`docker run -d --name=virtualshield-api-demo -p 8080:8080 saggafarsyad/virtualshield-api:demo`

Saggaf Arsyad <saggaf@area54labs.net>