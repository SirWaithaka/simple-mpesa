# Simple MPESA RESTful API

As the name of the project suggests, this is a simple implementation of MPESA. If you are East-African you already know what
MPESA is, but if you don't, I will have a description herein to help you understand inner workings that you may not have been
aware of.

## What is MPESA
>MPesa (M for mobile, pesa is Swahili for money) is a mobile phone-based money transfer service, payments and
>micro-financing service, launched in 2007 by Vodafone Group plc and Safaricom, the largest mobile network operator in
>Kenya. M-Pesa allows users to deposit, withdraw, transfer money, pay for goods and services (Lipa na MPesa), access
>credit and savings, all with a mobile device.
>
>The service allows users to deposit money into an account stored on their
>cell phones, to send balances using PIN-secured SMS text messages to other users, including sellers of goods and services,
>and to redeem deposits for regular money. Users are charged a small fee for sending and withdrawing money using the service.
>MPesa is a branchless banking service; MPesa customers can deposit and withdraw money from a network of agents that
>includes airtime resellers and retail outlets acting as banking agents.


## How does MPESA work
This is the big question. How does it work generally and how does it operate in the interior. We start with an overview.

### Overview
We can use the description above to make general inferences on how MPESA works. 

#### Who are MPESA customers (users)
To begin with, we can make inferences about what kind of customers MPESA operates with. We can get this information from
the description.

>The service allows users to deposit money into an account stored on their cell phones, to send balances using
>PIN-secured SMS text messages to other users, including sellers of goods and services, and to redeem deposits for
>regular money.
>
>MPesa customers can deposit and withdraw money from a network of agents that includes airtime resellers and retail
>outlets acting as banking agents.

1.  `Subscribers` - customers who can deposit and withdraw money into an account on their cell phones
2.  `Merchants` - customers with accounts that accept payments of goods and services
3.  `Agents` - customers that have the function of allowing other customers to deposit cash or withdraw cash from their
     accounts

We have managed to infer 3 types of customers that MPESA serves.

#### What are the MPESA operations
MPESA is a financial service, so we are more interested in the operations that involve the financial part. We are going
to ignore user registration and management at this point for now.
Let's look at this description:-

>M-Pesa allows users to deposit, withdraw, transfer money, pay for goods and services (Lipa na MPesa), access credit and
>savings, all with a mobile device.

1. `Deposit` - customers can credit their accounts
2. `Withdraw` - customers can debit their accounts
3. `Transfer` - customers can make transfers of moneys to other customers, through payments or regular transfers.

Those are the major 3 operations but if we go deeper we can also infer that customers can have some form of `Savings account`
and can `access credit (loans)`. From the above description we can't infer enough to know whether the savings account accrue
interest and whether the loans(credit) are some form of overdraft, short term or long term loan.

*TO KEEP THE API SIMPLE WE ARE GOING TO CONCENTRATE ON THE 3 MAJOR OPERATIONS.*


### MPESA Customers and their Operations

#### 1. Subscriber
A subscriber is a regular or standard customer with a cell phone and can perform the following operations

1. deposit money into their account
2. make transfers and receive money to/from other customers
3. make payments to merchants
4. can withdraw money 

On top of the above operations a subscriber `can retrieve full and mini statements` on their accounts and `fetch account
balance`.

#### 2. Merchant
A merchant is a customer that sells goods or services and usually has an account to majorly assist in receiving payments.
However, a merchant can do other operations with their account.

1.  receive payments
2.  make transfers to other customers
3.  can withdraw money

#### 3. Agent
An agent has 2 primary functions. From the descriptions, agents receive cash and can make a deposit on behalf of another
customer and agents can make a withdrawal on behalf of another customer and in turn give me cash.

1. make and receive transfers
2. can withdraw money
3. can deposit money
4. can make payments

An agent usually operates with a float account, and a separate commission account. The float account is an account they
use to facilitate withdrawals and deposits. During onboarding, an agent is required to have a minimum amount in
their float account.

##### - `Float Account`
When doing a deposit for a customer, the agent receives cash from the customer and initiates the transaction. Money will
be debited from the agent's float account and moved to the customer's account.

When a customer is withdrawing from an agent, money is debited from the customer's account and credited to the agent's
float account. The agent then gives the customer cash.

##### - `Commission Account`
An agent earns some commission from the system for enabling customer withdrawals. The commission is a percentage or a
fixed amount of the transaction cost that the customer pays when withdrawing.

That's a good brief overview of MPESA and how everything generally comes together.

### MPESA Components
Let's now discuss the components that make up MPESA as a whole. From what we have now gathered about MPESA we can infer
a some components. If we are following a Domain Driven Design we can call these components Domains, and if we are following
a service oriented design we can call then services. We are going to infer the major ones.

#### 1. Transaction (Component, Domain, Service)
We begin with what makes up one of the most important business aspect of a Financial Service. A transaction domain is
necessary to enable and record all transactions that are happening in the system between accounts and between users. It's
functions include:
1. enabling and recording all transactions regarding transfers of money between accounts
2. enabling and recording all transactions regarding deposits of money to accounts
3. enabling and recording all transactions regarding withdrawals of money from accounts
4. apply transaction cost on all applicable transactions

We can lay out all possible transactions that customers can perform between each other. We know we have the following
operations `transfer`, `deposit` and `withdraw`

##### i. `TRANSFER`
This is moving money from one account to another between customers. All customer types can make transfers, however some
transfers are treated specially. Also, some customer types cannot make transfers to other particular customers. Let's
define possible transfers that can happen.

1. A transaction between a `subscriber` and `subscriber`
2. A transaction between an `agent` and an `agent`
3. A transaction between a `merchant` and a `subscriber`
4. A transaction between a `subscriber` and a `merchant` is considered a `PAYMENT` but its inner workings is a transfer.
5. A transaction between an `agent` and a `merchant` is also considered a `PAYMENT` but its inner workings is a transfer.
6. 

##### ii. `DEPOSIT`
This is when a customer makes a cash deposit at an agent. Deposits can only be done at agents. We can define which customers
can make deposits.

1. A deposit between a `subscriber` and `agent`
2. A deposit between an `agent` and `agent`

##### iii. `WITHDRAW`
This when a customer makes a cash withdrawal at an agent. Just like deposits, withdrawals can only be done at agents. We
similarly define which customers can make withdrawals.

1. A withdrawal at an `agent` by a `subscriber`
2. A withdrawal at an `agent` by a `agent`
3. A withdrawal at an `agent` by a `merchant`

#### 2. Tariff (Component, Domain, Service)
> Users are charged a small fee for sending and withdrawing money using the service.

In order to return a profit (RoI), a financial service such as MPESA will charge a small fee for using the service on
selected operations.

MPESA applies a transaction fee on the following services:
- querying balance
- getting mini statement and full statement
- making a transfer to another customer
- making a payment to a merchant
- withdrawing from an agent

The tariff domain will be responsible for the following functions:

1. configuring and maintaining a tariff on selected operations.
2. given a type of transaction, tariff domain will return the configured cost.

When discussing the `Transaction domain`, we detailed and defined all possible transactions between customers. All these
transactions are possible areas a tariff could be applied. Apart from the transactions, a configured tariff can also be
applied to services offered such as `fetching mini and full statement` and `fetching balance` depending on the business
decisions.

#### 3. Account (Component, Domain, Service)
> The service allows users to deposit money into an account stored on their cell phones,

Once registered, customers need to have an account so that they can begin transacting and using the service. We can also
refer to them as wallets, but since for the Simple MPESA application we are going to allow a customer to hold only one
wallet, then we are going to refer to them as accounts. 

The account domain will be responsible for the following functions:
1. maintaining state and status of customer accounts
2. managing access to customer accounts i.e if an account is not marked active, we can't transact with it
3. recording and updating account balance after transactions

> _On a side note_
>
> In MPESA, your cell phone number is your account number, in contrast to a bank account where, you are given an
> account number when you successfully open an account. In the Kenyan variant of MPESA, that is owned by a Mobile Network
> Operator called Safaricom, the process of registration into MPESA is nearly similar to the process of registering your
> new SIM Card or is done at the same time. When getting a new Safaricom SIM Card a customer's KYC details are taken and in
> the registration form, a customer can opt in to use MPESA services. Your details for SIM registration are then going to
> be used to register and open an account for you with MPESA.

#### 4. Users/Customers (Component, Domain, Service)

We have 3 types of customers, we will have to find a way of managing KYC (Know Your Customer) data of our customers and limiting access to
different resources within the system by authenticating and authorizing based on the type of customer. On the other hand,
the system needs some form of administrator user(s) who runs the service and makes sure its available and reliable for our
customers to use and experience.

So we can infer that we have 2 types of users of our system i.e `administrator(s)` and `customers` and further, we have
3 types of customers i.e `subscribers`, `merchants` and `agents`. Each user of the system can perform specific actions,
based on how they authenticate.

Our user and customer domain will have the following responsibilities:

1. Manage and record user information(data)
2. Authenticate and authorize users
3. Provide and manage access to user data
4. 


## System Design And Architecture

To this point we have been using various resources we could acquire to help us infer about MPESA, and we have managed to
answer; `what it is?`, `how it works?`, `what is it made of?`.

Going forward we cannot infer anything about the System Design and Architecture unless we have access to the proprietary
code, so we build our own implementation that would work closely to how we experience using MPESA.

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
$ go build -o bin/wallet-server cmd/wallet-server.go 
$ ./bin/wallet-server
```

It will install all dependencies required and produce a binary for your platform.

##### Using the Dockerfile
Make sure you have docker installed and working properly.

```bash
$ docker build -t simple-mpesa:latest .
$ docker container create --name mpesa-server -p 6700:6700 --restart unless-stopped simple-mpesa
$ docker container start mpesa-server
```

The server will start at port `6700`.

Enjoy.

## API

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
POST /api/login/<user_type>                     <-- user_type can be either of agent, admin, merchant, subscriber
POST /api/user/<user_type> # for registration   <-- user_type can be either of agent, admin, merchant, subscriber
POST /api/admin/assign-float
POST /api/admin/update-charge
GET /api/admin/get-tariff
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
An agent can be registered to the api with the following `POST` parameters

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

##### 1. Assigning Float
Before you can start transacting, you need to login as an administrator and assign float to your `super-agent` using the
following endpoint

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

##### 2. Transfer Float to agents
The `super-agent` is limited to depositing to agents only. You will need to transfer the acquired float to other agents
you have registered.

##### 3. Configure Tariff
The default tariff in the system is set to zero amount for all chargeable transactions. You could begin testing transactions
using the default tariff and later choose to configure your own tariff. Choose your poison :-).

You can configure a tariff by updating the available charges. The system doesn't allow you to add any other charge band.

`GET /api/admin/get-tariff` - use this endpoint to get the available configured transaction charges
`POST /api/admin/update-charge` - use this endpoint to update a charge using its `id`. The amount should be in `cents`.

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