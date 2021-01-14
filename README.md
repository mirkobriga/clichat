## Create certs

```bash
chmod +x makecert
./makecer
```

## Run Server

```bash
go run clichat.go runserver :8080
```

## Run Client
```bash
go run clichat.go connect :8080
```