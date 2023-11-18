# fin-planning-backend

Welcome to the backend repository of the fin-planning project!  
Here is the frontend of the project: https://github.com/SgtMilk/fin-planning

This project uses containers:
- Development is better in vs-code. Re-open in a devcontainer.
- For Production, use `docker-compose up`

Some Commands: 
- `go run main.go` to run the server.
- `go test ./tests` to run tests.

## PostgreSQL DB
You will need to create a `.env.local` file from the `.env` file in the main directory of the repo. Inside, change your credentials.

## UML Diagram
![alt text](./assets/fin-planning.jpg)

## Where to find things
- For anything test related, the `./test` folder
- For the router and route handling, the `./controller` folder
- For middleware to the router, the `./middleware` folder
- For models and database instantiation, the `./database` folder
