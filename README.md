# store-api
A RESTful API built with Go and PostgreSQL for simple online store application, including products catalouge and ordering. Automated testing with Github Actions are also included.

## Stack
- [Go](https://go.dev/)
- [Fiber](https://gofiber.io/)
- [Gorm](https://gorm.io/)
- [Casbin](https://casbin.org/)
- [Viper](https://github.com/spf13/viper)
- PostgreSQL
- Docker

## Installation and running
If you have docker and docker compose installed in your system.
```bash
git clone https://github.com/n9mi/store-api.git
cd store-api
cp .env.example .env
docker compose up --build
```
or if you don't 
```bash
git clone https://github.com/n9mi/store-api.git
cd store-api
cp .env.example .env # configure .env according to your machine configurations
go mod tidy
go run cmd/web/main.go
```

## Structure
```bash
.
└── store-api/
    ├── .github/
    │   └── workflows/
    │       └── cicd.yml
    ├── cmd/                         # main entry point of the project for starting and running the app
    │   └── web/
    │       └── main.go
    ├── config/                      # configuration for Gorm, Fiber, Casbin, ...
    ├── database/                    # for creating, dropping, and seeding database on every app start/restart
    ├── internal/                    # internal code for the project
    │   ├── casbin/                  # casbin configuration
    │   │   ├── model.conf
    │   │   └── policy.csv
    │   ├── custom_handler/          # custom handler for Fiber error handling
    │   ├── delivery/                # contains code related to data delivery (in this project, HTTP implementation)
    │   │   └── http/
    │   │       ├── controller/      # contains handler for mapping users input/request and presented it back to user as relevant responses
    │   │       ├── middleware/      # contains code for handling incoming request, ex: authorization 
    │   │       └── route/           # specifying route for each controller handlers
    │   ├── dto/                     # data structures for request and response objects
    │   ├── entity/                  # domain objects
    │   ├── repository/              # contains handler for accessing db to perform set of manipulations on records 
    │   └── service/                 # contains set of logic to process the data
    ├── test/                        # contains testing for each controllers 
    ├── util/                        # helper functions 
    ├── .env.example
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    └── go.sum
```

## API
```bash
.
└── /api/v1/
    ├── auth/
    │   ├── register/
    │   │   └── POST:    # registering an user with 'customer' role
    │   └── login/
    │       └── POST:    # get a token for registered user by inputting user credentials
    └── customer/        # route for customer role
        ├── products/
        │   └── GET:     # get product list with its details
        ├── addresses/
        │   ├── GET:     # get current customer addresses
        │   └── POST:    # create current customer address
        ├── cart/
        │   ├── GET:     # get all items in customer cart
        │   └── POST:    # insert product and product quantity to customer cart
        ├── cart/:productId
        │   └── DELETE:  # delete product from customer cart
        └── order/
            └── POST:    # make an order from all products in customer cart and generate payment code
```
