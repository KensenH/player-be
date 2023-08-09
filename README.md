
# Player-be
Docker Image
```
docker pull kensenh/player-be:0.0.2
```
### Requirement
```
- go 1.20
- postgres
    port: 5432
    db: player
- redis
    port: 6379
```

### How to run:
```
Option 1:
1) extract postgres_data.zip to current directory
2) docker-compose up --build -d

Option 2:
1) extract postgres_data.zip to current directory
2) docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -v postgres_data:/var/lib/postgresql/data -p 5432:5432 -d postgres:alpine3.18
3) docker run --name some-redis -p 6379:6379 -d redis
4) docker run kensenh/player-be:0.0.2

Option 3:
1) docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
2) create database "player" in postgres
3) docker run --name some-redis -p 6379:6379 -d redis
4) go run cmd/player-be/main.go
```
### Postman link
```
https://api.postman.com/collections/25072264-6deb102f-7e40-4adb-8ae6-69f7db2dd5a1?access_key=PMAT-01H7DD7EP7122Q4MYRSEAZFFZR
```
### Endpoint
sign-up\
/player-be/api/v1/player/signup
```
request-body :
{
    "username" : "lut",
    "password" : "hut",
    "first_name" : "lut",
    "last_name": "hut",
    "phone_number": "+62813173723",
    "email": "luthut@gmail.com"
}
```

sign-in\
/player-be/api/v1/player/signin
```
BasicAuth
base64 of username:password

header example:
- Authorization: Basic bHV0Omh1dA==

cookie will be injected to client
```

sign-out
/player-be/api/v1/player/signout
```
cookie will be removed from client and blacklisted to redis
```

add bank account\
/api/v1/player/addbankaccount
```
Request Body:
{
    "bank_name": "bca",
    "account_owner_name": "luthut",
    "account_number": 12371233810236
}
```

top up\
/player-be/api/v1/player/topup
```
Request Body:
{
    "top_up_amount": 102667722
}
```

Get Player Detail\
/player-be/api/v1/player/detail/:id
```
example get player detail with playerId 1:
localhost:8080/player-be/api/v1/player/detail/1
```

Get Player's Profile\
/player-be/api/v1/player/profile
```
need cookie
```

Get Player's TopUp Histories\
/player-be/api/v1/player/receipts
```
need cookie
```

Get all player / search player / filter player \
/player-be/api/v1/player/search?
```
query-params:
- join_before (ex: 10-08-2023)
- join_after (ex: 10-08-2023)
- max_ingame_currency (int)
- min_ingame_currency (int)
- username (ex: kensen)
- player_id (int)
- bank_name (ex: bca)
- bank_account_name (ex: kensen)
- bank_account_number (int)

```
