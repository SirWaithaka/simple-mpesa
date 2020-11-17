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
their float account before they can begin transacting. Agents can acquire this float from a `super-agent`.

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
