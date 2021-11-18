# Simple MPESA RESTful API

As the name of the project suggests, this is a simple implementation of MPESA. If you are East-African you already know what
MPESA is, but if you don't, I will have a description herein to help you understand inner workings that you may not have been
aware of.

[*Click here to learn about What MPESA is*](https://github.com/SirWaithaka/simple-mpesa/blob/primary/mpesa.md)

## System Design And Architecture

Whatever follows, is the design and architecture of our simple MPESA application.

### System Architecture

This application is built as a `Monolith` application following Domain Driven Design principles. The application uses a
postgres database for storage.

### System Design

Domain Driven Design is at the heart of our simple mpesa application.

#### Bounded Contexts
DDD principles have concepts called contexts. The application uses the following bounded contexts:

1. Admin
2. Agent
3. Merchant
4. Subscriber
5. Transaction
6. Account
7. Statement
8. Tariff
9. Auth
10. Customer

##### 1. Admin Context
The application needs some form of administration by a super user charge with the responsibility of running and maintaining
the application to ensure reliability and stability. This user is an `admin` and is given their own bounded context. Some
responsibilities/actions of this user are:

1. Can login to system or register.
2. Can assign float to a Super Agent.
3. Can configure tariff
4. Can suspend/change status of a customer account
5. Can view/edit/delete customer accounts

As the application grows and scales the administrator context would have more responsibilities.

1. The application would need more than one administrator and more so more than one category of administrators. In the
case of MPESA, some examples of administrators with their roles include:

    i. Customer Care - is a part admin who would assist customers with information about the system and troubleshoot
    problems.
    
    ii. Finance - is a part admin whose responsibilities would be financial and accounting aspect in the system.
    
    iii. IT - an admin whose responsible for the infrastructure that the system runs on.
    
    iv. 

##### 2. Agent Context
We have acquired or developed a wallet and money transfer service for a Telco and we have been given the go ahead by the
Central Bank to deploy the application and get some customers to use our system. There are however some initial steps
the business has to perform to start onboarding new customers. The system should be able to have some level of autonomy
when it comes to the flow of money. That is where agents come in.

> The initial obstacle in the pilot was gaining the agent’s trust and encouraging them to process cash withdrawals
> and agent training.
>
> *Source [wikipedia](https://en.wikipedia.org/wiki/M-Pesa#cite_note-8)*

Our first initial steps before we can roll out

1. Acquire and entity licensed to hold public money. A bank.
2. Create a super agent(s) whose task would be depositing money to our bank account.
3. Once a super agent deposits to our account, we assign them with an equivalent amount of float they can sell
to other agents.
4. When we onboard an ordinary agent, they will have a balance of zero, and they will approach the super agent to get
float.

Agents are important customers to the system. They can also have various categories depending on the business use case.
For our example we have 2 types of agents:

1. Super Agent
2. Ordinary Agent

##### 3. Merchant Context
MPESA has 2 types of merchants.

1. A merchant to provides utility services to their customers
2. A merchant that sells goods and services to customers

Both merchants have unique ways of how customers pay for their services/goods. However both merchants have an account number.

1.  Pay bill number as an account number for a merchant - Customer provides the `pay bill number`, `a customer
account number` and the `amount`.
2.  Till number as an account number for a merchant - Customer provides the `till number` and `amount`.

`Pay bill number` is usually given to utility companies that need to identify from whom the payment is coming from by the
`customer account number`.

`Till number` is usually given to small scale traders that want to accept payment via MPESA from their customers.

For our example we stick to one general merchant that accepts payment.

##### 4. Subscriber Context
A subscriber does not have much going on. They can authenticate and perform a transaction.

##### 5. Transaction Context
Contains all business logic in regard to transactions happening in the system. It enforces the transaction rules and
business policy.

Business policies:

1. A transaction cannot happen between identical customers i.e a customer cannot transact with themselves
2. A deposit cannot be done by none other customer than an agent
3. A customer cannot perform a withdrawal with no other customer than an agent
4. A super agent is however only allowed to do deposits for other agents only
5. Customers are not allowed to deposit, withdraw or transfer money below the minimum amount allowed
6. Apply transaction fee as per the tariff configured

##### 6. Account Context
The main responsibility of this context is managing customer accounts/wallets. Responsibilities:

1. Updating account balances, credit/debit accounts
2. Updating system ledger after changing account balances
3. 

##### 7. Statement Context
The main responsibility of this context is managing the system ledger. If we scale the system, we can view this ledger
as the statements/transactions event store. Borrowing from `event sourcing` design, our statement context is a record
of every event with customer transactions.

##### 8. Tariff Context
This context has a responsibility of configuring and maintaining the tariff used in various transactions.

##### 9. Auth Context
The system has 4 different types of users, `admin`, `agent`, `subscriber` and `merchant`. The auth context is responsible
for authenticating and authorizing these users into the system.

##### 10. Customer Context
This context is mainly an aggregator of the `agent`, `merchant` and `subscriber` contexts. It exposes common functionality
for, which can be used in other core contexts.



## Installation

To begin with, the application uses postgres as the backend database.

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
--name mpesa-db \
-e POSTGRES_USER=mpesa \
-e POSTGRES_PASSWORD=mpesa \
-p 5432:5432 \
postgres
```

Run the following command to start the container
```bash
$ docker start mpesa-db
```

### Simple MPESA Application Installation

Running the application is as simple as running any other go application but first we need to copy and create our
configuration.

#### Lets begin. Cloning ...

```bash
$ git clone https://github.com/SirWaithaka/simple-mpesa.git
```

#### Configuring
```bash
$ cd simple-mpesa
$ cp config.yml.example config.yml
```

This configuration file looks something like this

```yaml
database:
  host: "127.0.0.1"
  port: "5432"
  user: "mpesa"
  password: "mpesa"
  dbname: "mpesa"

app_secret_key: "eQig7GS4cHO2su"
```

You can change the config variables depending on your database setup, here I choose to follow the default setup shown at
database installation step.

#### Building and running

##### Using the Binary
```bash
$ mkdir bin
$ go build -o bin/mpesa-server cmd/mpesa-server.go 
$ ./bin/mpesa-server
```

It will install all dependencies required and produce a binary for your platform.

##### Using the Dockerfile
Make sure you have docker installed and working properly.

```bash
$ docker build -t simple-mpesa:latest .
$ docker container create --network=host --name mpesa-server --restart unless-stopped simple-mpesa
$ docker container start mpesa-server
```

The server will start at port `6700`.

Enjoy.

## API Usage

A description of the api.

### Endpoints

All the routes exposed in the application are all defined in this function
```go
func apiRouteGroup(api fiber.Router, domain *registry.Domain, config app.Config) {

	api.Post("/login/:user_type", user_handlers.Authenticate(domain, config))
	api.Post("/user/:user_type", user_handlers.Register(domain))

	// create group at /api/admin
	admin := api.Group("/admin", middleware.AuthByBearerToken(config.Secret))
	admin.Post("/assign-float", user_handlers.AssignFloat(domain.Admin))
	admin.Post("/update-charge", user_handlers.UpdateCharge(domain.Tariff))
	admin.Get("/get-tariff", user_handlers.GetTariff(domain.Tariff))
	admin.Put("/super-agent-status", user_handlers.UpdateSuperAgentStatus(domain.Agent))

	// create group at /api/account
	account := api.Group("/account", middleware.AuthByBearerToken(config.Secret))
	account.Get("/balance", account_handlers.BalanceEnquiry(domain.Account))
	account.Get("/statement", account_handlers.MiniStatement(domain.Statement))

	// create group at /api/transaction
	transaction := api.Group("/transaction", middleware.AuthByBearerToken(config.Secret))
	transaction.Post("/deposit", transaction_handlers.Deposit(domain.Transactor))
	transaction.Post("/transfer", transaction_handlers.Transfer(domain.Transactor))
	transaction.Post("/withdraw", transaction_handlers.Withdraw(domain.Transactor))
}
```

The routes are mounted on the prefix `/api` so your requests should point to
```
POST /api/login/<user_type>                     <-- user_type can be either of agent, admininistrator, merchant, subscriber
POST /api/user/<user_type> # for registration   <-- user_type can be either of agent, admininistrator, merchant, subscriber
POST /api/admin/assign-float
POST /api/admin/update-charge
GET /api/admin/get-tariff
PUT /api/admin/super-agent-status
GET /api/account/balance
POST /api/account/statement
POST /api/transaction/deposit
POST /api/transaction/transfer
POST /api/transaction/withdraw
```

#### To Register

The api can be used to register 4 types of users: `admin`, `agent`, `merchant` and `subscriber`

#### Admin Registration
An admin can be registered to the api with the following `POST` parameters

`firstName`, `lastName`, `email`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/administrator \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Admin \
  --data lastName=Waithaka \
  --data email=admin@email.com \
  --data password=mnbvcxz
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "ac8f944b-b0aa-4029-9caf-dfe67007bc84",
    "userType": "administrator"
  }
}
```

##### Agent Registration
At minimum, you need to create 2 agents, one of which will become a `super agent`. An agent can be registered to the api
with the following `POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/agent \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Agent \
  --data lastName=Waithaka \
  --data email=agent_waithaka@email.com \
  --data phoneNumber=254700000000 \
  --data password=mnbvcxz
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "cca7d227-74ae-4d47-aae8-a0ab952aac28",
    "userType": "agent"
  }
}
```

##### Merchant Registration
A merchant can be registered to the api with the following `POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/merchant \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Merchant \
  --data lastName=Waithaka \
  --data email=merch_waithaka@email.com \
  --data phoneNumber=254700000000 \
  --data password=mnbvcxz
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "c3a81710-ef66-47d9-adc9-f365a324ed5c",
    "userType": "merchant"
  }
}
```

##### Subscriber Registration
A subscriber can be registered to the api with the following `POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/subscriber \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Subscriber \
  --data lastName=Waithaka \
  --data email=subscriber_waithaka@email.com \
  --data phoneNumber=254700000000 \
  --data password=mnbvcxz
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
    "userType": "subscriber"
  }
}
```

#### To Login
You can use the following `POST` parameters for login with any of the 4 users
`email`, `password`

Curl request example for subscriber login
```bash
curl --request POST \
  --url http://localhost:6700/api/login/subscriber \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data email=subscriber_waithaka@email.com \
  --data password=mnbvcxz
```

Response example

```json
{
  "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
  "userType": "subscriber",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDU0MTYwNjAsImlhdCI6MTYwNTM5NDQ2MH0.lAJ4WpF2Mnfg52iuTOoPV8nvbHV3JrMQOC-5xXrQ5EE"
}
```

**NOTE**: The remaining endpoints require the token acquired above for authentication


#### Initial Steps Before Transacting
There are some initial setups that need to be done before you can begin doing transactions.

##### 1. Creating a super agent
Before you can start transacting, you need to login as an administrator and create a super agent by changing the status
of an existing agent. When registering an agent, you ought to have created at minimum 2 agents. It is now that we need
make one of those agents a super agent.

The following endpoint is used to update the `super agent status` of an agent.

`PUT /api/admin/super-agent-status` requires the following post parameters: `email`

Curl request example
```bash
curl --request PUT \
  --url http://localhost:6700/api/admin/super-agent-status \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijc2YmM0YWEzLTAyNWQtNGQ1YS1hNWZiLWY1NDk1NTdmNjM0YSIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDU0NTE4MDUsImlhdCI6MTYwNTQzMDIwNX0.8lTWl9hGr9GTST7WpEpzKdm_gqhMkf4qUellLx4o5bw' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data email=agent_waithaka@email.com
```

Response example

```json
{
  "status": "success",
  "message": "Super Agent Status updated"
}
```

##### 2. Assigning Float
Logged in as an administrator, you need to assign float to your `super-agent` using the following endpoint

`POST /api/admin/assign-float`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/admin/assign-float \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImE3OGVjNjNhLTA0ZWItNDAzNC1iZmVkLTBhNmMwMjU3ZTJlNCIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDUzMjExNzUsImlhdCI6MTYwNTI5OTU3NX0.4fGlMJQB-eylKwOAwa4d16nVQQt3uYgwPbUjYt7j9zA' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data accountNo=agent_waithaka@email.com \
  --data amount=100000
```

Response example
```json
{
  "status": "success",
  "message": "Float has been assigned.",
  "data": {
    "balance": 100000
  }
}
```

##### 3. Transfer Float to agents
The `super-agent` is limited to depositing to agents only. You will need to transfer the acquired float to other agents
you have registered.

##### 4. Configure Tariff
The default tariff in the system is set to zero amount for all chargeable transactions. You could begin testing transactions
using the default tariff and later choose to configure your own tariff. Choose your poison :-).

You can configure a tariff by updating the available charges. The system doesn't allow you to add any other charge band.

`GET /api/admin/get-tariff` - use this endpoint to get the available configured transaction charges

Response example
```json
{
  "status": "success",
  "message": "Tariff retrieved",
  "data": [
    {
      "id": "acf3e6bf-c9de-45b4-a8b6-bf97f92b783a",
      "txnOperation": "WITHDRAW",
      "srcUserType": "subscriber",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "0e5a4aaa-135a-4464-96c9-d021f769bdb7",
      "txnOperation": "WITHDRAW",
      "srcUserType": "merchant",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "243e7ecc-c2dd-41bb-9953-1278050bfb64",
      "txnOperation": "WITHDRAW",
      "srcUserType": "agent",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "f8835176-316c-49de-b001-687e2c4a338d",
      "txnOperation": "TRANSFER",
      "srcUserType": "agent",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "4edeb6d0-37cd-4c67-997a-0b3fa93b722d",
      "txnOperation": "TRANSFER",
      "srcUserType": "subscriber",
      "destUserType": "subscriber",
      "fee": 0
    },
    {
      "id": "94c0ae8b-a131-41b9-b5af-5235b8926fa4",
      "txnOperation": "TRANSFER",
      "srcUserType": "merchant",
      "destUserType": "subscriber",
      "fee": 0
    },
    {
      "id": "450e4baa-58c3-41b3-abe5-a55555492e0c",
      "txnOperation": "TRANSFER",
      "srcUserType": "subscriber",
      "destUserType": "merchant",
      "fee": 0
    },
    {
      "id": "3623a89f-c496-41c8-b6c9-73429cc4ef9d",
      "txnOperation": "TRANSFER",
      "srcUserType": "agent",
      "destUserType": "merchant",
      "fee": 0
    }
  ]
}
```

`POST /api/admin/update-charge` - use this endpoint to update a charge using its `id`.

You need the following `POST` parameters

`amount`, `chargeId` - The amount should be in `cents`.

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/admin/update-charge \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijc2YmM0YWEzLTAyNWQtNGQ1YS1hNWZiLWY1NDk1NTdmNjM0YSIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDU0NTE4MDUsImlhdCI6MTYwNTQzMDIwNX0.8lTWl9hGr9GTST7WpEpzKdm_gqhMkf4qUellLx4o5bw' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=1050 \
  --data chargeId=acf3e6bf-c9de-45b4-a8b6-bf97f92b783a
```

Response example
```json
{
  "status": "success",
  "message": "charge configured"
}
```

#### Performing Transactions
While configuring a charge requires you to provide the amount in `cents`, performing transactions requires the amount to
be in whole units i.e. `shillings`

Transacting also requires you to provide an `accountNo`, use the `email` of the customer as the `accountNo`

`customerType` can be either of `agent`, `merchant` or `subscriber`

##### 1. To Deposit
A deposit is only done by an `agent`. You need an `agent` token to perform this transaction.

You need the following `POST` parameters

`amount`, `accountNo` and `customerType`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/deposit \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNjYTdkMjI3LTc0YWUtNGQ0Ny1hYWU4LWEwYWI5NTJhYWMyOCIsInVzZXJUeXBlIjoiYWdlbnQifSwiZXhwIjoxNjA1MzIxODEyLCJpYXQiOjE2MDUzMDAyMTJ9.jFLfjScuvHaOV68n11sRticy2ntzQRhwbNq5E4sPmQI' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=400 \
  --data accountNo=subscriber_waithaka@email.com \
  --data customerType=subscriber
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
``` 

##### 2. To Withdraw
You need the following `POST` parameters

`amount`, `agentNumber`

Use agent email for `agentNumber`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/withdraw \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDU0MTYwNjAsImlhdCI6MTYwNTM5NDQ2MH0.lAJ4WpF2Mnfg52iuTOoPV8nvbHV3JrMQOC-5xXrQ5EE' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=40 \
  --data agentNumber=agent_waithaka@email.com
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
```

##### 3. To Transfer
You need the following `POST` parameters

`amount`, `accountNo` and `customerType`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/transfer \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDUzMjE5MDMsImlhdCI6MTYwNTMwMDMwM30.vLiHdNTr4onTVqUZbLbdpwgbH98VYzHJJU-JKtFOHVg' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=30 \
  --data accountNo=merch_waithaka@email.com \
  --data customerType=merchant
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
```


#### To Query Balance
This is just a `GET` request, no params

Curl request example
```bash
curl --request GET \
  --url http://localhost:6700/api/account/balance \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDUzNjY3NTMsImlhdCI6MTYwNTM0NTE1M30.-Piib6bXzYqb0S8nLo76SBTyGmWi7UPUMExptIcqBZI'
```

Response example

```json
{
  "status": "success",
  "message": "Your current balance is 690",
  "data": {
    "userID": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
    "balance": 690
  }
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
  "status": "success",
  "message": "mini statement retrieved for the past 5 transactions",
  "data": {
    "message": "mini statement retrieved for the past 5 transactions",
    "userID": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
    "transactions": [
      {
        "transactionId": "97c3ff6d-72d5-479d-8838-85a5c32985a2",
        "transactionType": "DEPOSIT",
        "createdAt": "2020-11-14T01:59:28.613007+03:00",
        "creditedAmount": 400,
        "debitedAmount": 0,
        "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      },
      {
        "transactionId": "4be4a008-b18e-4d4b-95d3-58b660d5b931",
        "transactionType": "TRANSFER",
        "createdAt": "2020-11-14T01:59:05.949066+03:00",
        "creditedAmount": 0,
        "debitedAmount": 30,
        "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      },
      {
        "transactionId": "45da6c6a-03d8-4d58-849a-fd80bbfabbb4",
        "transactionType": "TRANSFER",
        "createdAt": "2020-11-14T01:57:04.622507+03:00",
        "creditedAmount": 0,
        "debitedAmount": 40,
        "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      }
    ]
  }
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
