# Simple Wallet RESTful API

The name of the project is pretty much descriptive enough. To add a little more on the name, the application does the
stuff you need from a digital wallet i.e.
   1. Registration and Login
   2. Account Deposit
   3. Account Withdrawal
   4. Balance Enquiry
   5. Mini Statement
   
PS. This is fully a backend application with no front end client whatsoever. You can use curl, postman or your favourite
http tool to interact with the application.

You can follow the below steps to install and setup the application.

## Installation

To begin with the application uses postgres as the backend database.

### Database Installation

The application uses postgres as the database server. So here are instructions on how to setup postgresql
on your machine using docker.

Get the official postgres docker image.
```bash
$ docker pull postgres
``` 

Then create a container from the image with the following variables
```bash
$ docker create \
--name wallet-db \
-e POSTGRES_USER=wallet \
-e POSTGRES_PASSWORD=wallet \
-p 5432:5432 \
postgres
```

Run the following command to start the container
```bash
$ docker start wallet-db
```

### Simple Wallet Installation

Running the application is as simple as running any other go application but first we need to copy and create our
configuration.

#### Lets begin. Cloning ...

```bash
$ git clone https://github.com/SirWaithaka/simple-wallet.git
```

#### Configuring
```bash
$ cd simple-wallet
$ cp config.yml.example config.yml
```

This configuration file looks something like this

```yaml
database:
  host: "127.0.0.1"
  port: "5432"
  user: "wallet"
  password: "wallet"
  dbname: "wallet"

app_secret_key: "eQig7GS4cHO2su"
```

You can change the config variables depending on your database setup, here i choose to follow the default setup shown at
database installation step.

#### Building and running

##### Using the Binary
```bash
$ cd main
$ go build wallet-server.go
$ ./wallet-server
```

It will install all dependencies required and produce a binary for your platform.

##### Using the Dockerfile
Make sure you have docker installed and working properly.

```bash
$ docker build -t simple-wallet:latest .
$ docker container create --name wallet-server -p 6700:6700 --restart unless-stopped simple-wallet
$ docker container start wallet-server
```

The server will start at port `6700`.

Enjoy.

## API

A description of the api.

### Endpoints

All the routes exposed in the application are all defined in this function
```go
func apiRouteGroup(g fiber.Router, domain *registry.Domain, config app.Config) {

	g.Post("/login", user_handlers.Authenticate(domain.User, config))
	g.Post("/user", user_handlers.Register(domain.User))

	g.Get("/account/balance", middleware.AuthByBearerToken(config.Secret), account_handlers.BalanceEnquiry(domain.Account))
	g.Post("/account/deposit", middleware.AuthByBearerToken(config.Secret), account_handlers.Deposit(domain.Account))
	g.Post("/account/withdrawal", middleware.AuthByBearerToken(config.Secret), account_handlers.Withdraw(domain.Account))
	g.Post("/account/withdraw", middleware.AuthByBearerToken(config.Secret), account_handlers.Withdraw(domain.Account))

	g.Get("/account/statement", middleware.AuthByBearerToken(config.Secret), account_handlers.MiniStatement(domain.Transaction))
}
```

The routes are mounted on the prefix `/api` so your requests should point to
```
POST /api/login
POST /api/user # for registration
GET /api/account/balance
POST /api/account/deposit
POST /api/account/withdrawal
GET /api/account/statement
```

#### To Register
A user can be registered to the api with the following `POST` parameters
`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Sir \
  --data lastName=Waithaka \
  --data email=newme@email.com \
  --data phoneNumber=254700000000 \
  --data password=mnbvcxz
```

Response example
```json
{
    "status": "success",
    "message": "user created",
    "user": {
        "email": "newme@email.com",
        "userId": "b4b00501-ba22-49fb-827d-b25d969c58bb"
    }
}
```

#### To Login
You can use `email and password` or `phoneNumber and password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/login \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data password=mnbvcxz \
  --data phoneNumber=254700000000
```

Response example

```json
{
    "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijg0ODA5YTAyLTkwODItNGFlMy05MDQ3LTM4NDA5NDhjNTdjZiIsImVtYWlsIjoiaGFsbEBlbWFpbC5jb20ifSwiZXhwIjoxNTg0MDI0NTQ0LCJpYXQiOjE1ODQwMDI5NDR9.qZHLJWtYK7_ClgnaPJbGuaiPW8ssd1Ra9xFJWdg6iwE"
}
```

**NOTE**: The remaining endpoints require the token acquired above for authentication

#### To Deposit
You only need the `amount` parameter

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/account/deposit \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijk4YmNmMmY1LWRiY2ItNDk1NS04NTU0LTc0OWYxMTVhZjU5OCIsImVtYWlsIjoiIn0sImV4cCI6MTYwNDA2Mjg0NiwiaWF0IjoxNjA0MDQxMjQ2fQ.Z0oFwOV3wEiQzpwLg4LH5NZIBUsllDhcJefgvceMiHw' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=1000
```

Response example

```json
{
  "balance": 1000,
  "message": "Amount successfully deposited new balance 1000",
  "userId": "98bcf2f5-dbcb-4955-8554-749f115af598"
}
``` 

#### To Withdraw
You only need the `amount` parameter

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/account/withdrawal \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijk4YmNmMmY1LWRiY2ItNDk1NS04NTU0LTc0OWYxMTVhZjU5OCIsImVtYWlsIjoiIn0sImV4cCI6MTYwNDA2OTE0MywiaWF0IjoxNjA0MDQ3NTQzfQ.IYyclrC66aweehs_A4Sigmc83a27udmPofM2yOeut9Q' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=40
```

Response example

```json
{
  "balance": 880,
  "message": "Amount successfully withdrawn. New balance 880",
  "userId": "98bcf2f5-dbcb-4955-8554-749f115af598"
}
```

#### To Query Balance
This is just a `GET` request, no params

Curl request example
```bash
curl --request GET \
  --url http://localhost:6700/api/account/balance \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijk4YmNmMmY1LWRiY2ItNDk1NS04NTU0LTc0OWYxMTVhZjU5OCIsImVtYWlsIjoiIn0sImV4cCI6MTYwNDA2Mjg0NiwiaWF0IjoxNjA0MDQxMjQ2fQ.Z0oFwOV3wEiQzpwLg4LH5NZIBUsllDhcJefgvceMiHw'
```

Response example

```json
{
    "message": "Your current balance is 4700",
    "balance": 4700
}
```

#### To Get Mini Statement
This is just a `GET` request, no params

Curl request example
```bash
curl --request GET \
  --url http://localhost:6700/api/account/statement \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijk4YmNmMmY1LWRiY2ItNDk1NS04NTU0LTc0OWYxMTVhZjU5OCIsImVtYWlsIjoiIn0sImV4cCI6MTYwNDA2OTE0MywiaWF0IjoxNjA0MDQ3NTQzfQ.IYyclrC66aweehs_A4Sigmc83a27udmPofM2yOeut9Q'
```

Response example

```json
{
    "message": "ministatement retrieved for the past 5 transactions",
    "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
    "transactions": [
        {
            "transactionId": "1088b880-1aa1-4cf0-929b-dfab01c52c13",
            "transactionType": "balance_enquiry",
            "timestamp": "2020-03-12T13:11:09.863693Z",
            "amount": 4700,
            "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
            "accountId": "fd1ce4e4-e467-4eac-8ea0-ea0c9d4f76fe"
        },
        {
            "transactionId": "8b64ca58-b869-47b6-964a-0846957d4c7f",
            "transactionType": "withdrawal",
            "timestamp": "2020-03-12T12:30:37.355034Z",
            "amount": 4700,
            "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
            "accountId": "fd1ce4e4-e467-4eac-8ea0-ea0c9d4f76fe"
        },
        {
            "transactionId": "4514a94c-f303-4324-b010-a4e7c3dd3f77",
            "transactionType": "withdrawal",
            "timestamp": "2020-03-12T12:30:36.278053Z",
            "amount": 4710,
            "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
            "accountId": "fd1ce4e4-e467-4eac-8ea0-ea0c9d4f76fe"
        },
        {
            "transactionId": "871035ad-456a-4467-8260-414b464b6d86",
            "transactionType": "withdrawal",
            "timestamp": "2020-03-12T12:30:35.446227Z",
            "amount": 4720,
            "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
            "accountId": "fd1ce4e4-e467-4eac-8ea0-ea0c9d4f76fe"
        },
        {
            "transactionId": "c3c3a19a-fc3a-4080-9b52-5d3e19853cd7",
            "transactionType": "withdrawal",
            "timestamp": "2020-03-12T12:30:34.646326Z",
            "amount": 4730,
            "userId": "84809a02-9082-4ae3-9047-3840948c57cf",
            "accountId": "fd1ce4e4-e467-4eac-8ea0-ea0c9d4f76fe"
        }
    ]
}
```

## Testing

Tests have not been written for the application but i have very important talks i would share here that i cant recommend
enough about how to go about testing the application.

1. [Ian Cooper - TDD, Where Did It All Go Wrong](https://www.youtube.com/watch?v=EZ05e7EMOLM)
2. 

An approach to testing this application would be something in the following lines.

1. Test the code in the interactor files
2. Test the code in the repository files

The above files carry the bulk of the behaviour of the whole application, they are the business logic of the application.
The rest of the files are just implementation details that could change rapidly and the tests written for them would
certainly fail after change.

e.g.
the http handler functions in the application use [gofiber](https://github.com/gofiber/fiber), writing unit tests for
them is good but not desired, because [gofiber](https://github.com/gofiber/fiber) can be replaced with [mux](https://github.com/gorilla/mux) easily and that
would break your tests.