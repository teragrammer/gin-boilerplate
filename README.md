# Golang GIN Boilerplate API
```
A minimal and clean Golang GIN boilerplate for building RESTful APIs quickly. 
This starter template includes essential features like routing, middleware setup, error handling, and 
environment configuration to help you kickstart your API development with best practices.
```

### Getting Started
- Clone the repository
```
$ git clone https://github.com/teragrammer/gin-boilerplate.git
$ cd gin-boilerplate
```

- Configure your .env (.env.example) for docker port (optional)
- Configure your env.json (env.json.example)

- Initialize Docker
```
$ docker compose up -d
$ sh cmd/docker/start.sh
```

- To Update all Packages (Optional)
```
$ go get -u ./...
```

- Install dependencies
```
$ go mod download
$ go mod tidy
```

- Migration & Seed
```
$ clear && go run ./cmd/database/migrate.go [-seed]
```

- Unit Test
```
$ clear && go clean -testcache && go test ./... -v
```

- Running for Development Mode
```
$ clear && nodemon --watch 'internal' --signal SIGTERM --exec 'go' run ./cmd/server/main.go
```

### Hire Me
```
If you like this project and need help with development, customization, or integration, feel free to reach out!

I’m available for freelance work, consulting, and collaboration.

Thank you for checking out Golang GIN Boilerplate API!
Feel free to contribute or open issues.
```