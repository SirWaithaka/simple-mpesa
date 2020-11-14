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
6.

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

### Simple MPESA Application Installation

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