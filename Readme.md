# Loan APP

## Local Setup
1. git clone 
2. create dev.toml from default.toml
3. replace db values with your username and password to connect to mysql
4. run migrations <br>
``
APP_ENV=dev go run cmd/migration/main.go up
``
5. run server <br>
``
   APP_ENV=dev go run cmd/api/main.go
``


## Postman Collection
https://www.postman.com/red-desert-269362/workspace/loan/collection/9524554-39674238-4cb0-46ee-874a-c674be34fb93?action=share&creator=9524554

## Postman Env
https://www.postman.com/red-desert-269362/workspace/loan/environment/9524554-65ffbb65-ae88-4a85-8884-83c118244082


Did not get much time to add a lot of tests. have skipped unit tests and 
added functional test for happy flow for loan.
Have mentioned possible test cases in loan_test.go
