# Problem Set "Payments"


## API Schema

- POST /auth/signup
- POST /auth/login
- GET /users
- GET /users/{id}
- PUT /users/{id}
- GET /accounts
- GET /accounts/{id}
- PUT /accounts/{id}
- GET /payments
- GET /payments/{id}
- PUT /payments/{id}

## DB Schema

`user` table:

- id
- name
- email
- password_hash
- is_admin
- is_active

`account` table:

- id
- user_id
- name
- balance
- is_active

`payment` table:

- id
- sender_id
- recipient_id
- amount
- created_at
- completed_at
- status (prepared, completed) // [https://www.postgresql.org/docs/current/datatype-enum.html]

`log` table:

- id
- date
- category (block_card, unblock_card, create_user, block_user, unblock_user, balance)
- user_id
- description

## Docker commands

```term
docker run --name payments_db -p 127.0.0.1:5432:5432 -e POSTGRES_PASSWORD=lthgfhjk -e POSTGRES_USER=payments -e POSTGRES_DB=payments -d postgres
```

## CURL commands

```term
# auth/signup

curl -i -X POST -H "Content-Type: application/json" -d '{"name":"Mike","email":"mike@somemail.com","password":"qwerty"}' localhost:8080/auth/signup

```
