# Frete Rápido API

### Running the application

1. Clone the repository:
```sh
  git clone https://github.com/belmadge/freteRapidoTask.git
```

2. Create a `.env` file:
```env
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=frete_rapido
   DB_HOST=mysql
   DB_PORT=3306
```

3. Build and run the application using Docker Compose:
```sh
  docker-compose up --build
```

4. Testes:
- Para executar os testes, utilize o comando:
```sh
  go test ./...
```

5. To see the API documentation access the directory "doc" in this project
