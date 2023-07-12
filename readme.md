

# To Run this project 

```json
   Docker compose up
```



# OVERVIEW OF THIS PROJECT 
### Here are some of the things that I have dene in this project:
> - I've done a lot of work on the database layer. This is important, as the database is the foundation of any application.
>-  I've also implemented a number of security features, such as JWT and PASETO tokens, authentication middleware, and authorization rules. This is important for ensuring the security of your application.
>- I've used a number of different technologies, such as Go, Postgres, Gin, gRPC, Swagger, and Redis. This shows my expertise  with a wide range of technologies and can use them to build a production-ready application
>- I've also automated a number of tasks, such as building and pushing Docker images, deploying to Kubernetes, and issuing TLS certificates. This shows that I'm using modern DevOps practices to build and deploy your application.


# Project breakdown
1. Designed a database schema and generated SQL code using diagram.io
2. Implemented migrations
3. Generated CRUD code for Golang using SQL for Postgres
4. Efficiently handled database transaction locks and deadlocks
5. Came up with an efficient strategy for deadlock avoidance in Golang
6. Set up Github Actions for Golang + Postgres to run automated tests
7. Implemented a RESTful HTTP API using Gin
8. Ran a mock database for testing the HTTP API in Go and achieved 100%
9. test coverage
10. Implemented a transfer money API with a custom parameters validator
11. Implemented a JWT and PASETO token
12. Implemented a user session manager with refresh token
13. Built a minimal Golang Docker image using a multi-stage Dockerfile
14. Introduced gRPC to the project
15. Defined a gRPC API and generated Go code from protobuf
16. Implemented the Golang gRPC server for the API
17. Integrated the gRPC API to create and manage users
18. Made the API compatible with both HTTP and gRPC requests
19. Automated the generation and serving of Swagger documentation from the
20. Go server using OpenAPI
21. Implemented a validator for gRPC parameters and sent friendly error messages
22. Added authorization to protect the gRPC API
23. Implemented structured logging for the C APIs
24. Implemented a HTTP logger middleware
25. Integrated a background worker with Redis and asynq for the email server
26. Handled all possible errors and printed logs for Asynq
27. Implemented the send and verification email with Gomail
28. Wrote tests for the gRPC API that requires authentication
29. Implemented automatic building and pushing of the Docker image to AWS ECR with Github Actions
30. Created a Postgres database on AWS RDS
31. Stored and retrieved production secrets with AWS Secrets Manager
32. Created an EKS cluster on AWS
33. Deployed the API_BANK to a Kubernetes cluster on AWS EKS
34. Registered a domain and set up a record in Route53
35. Used ingress to route traffic to different services in Kubernetes
36. Implemented an automatic issuance of TLS certificates in Kubernetes with Let'sencrypt
37. Finally, automated deployment to the internet (EKS) with Github Actions



<!-- ### Postgres


### Create a migration 
migrate create -ext sql -dir db/migration -seq init_schema
<br /><br />

### SQLC
-   Very fast and easy to user 
-   Automatic ode generation 
-   Catch SQL query errors before genration codes 
-   Full  support Postgres.
(TODO: Publis a toturial on sqlc golang)
<br /><br />


## Database Transactions 
### Why do we need database transactions 
1.   To provide a reliable and consistent unit of work, even inc ase of system failure 
2.   To provider isolation between programs tha access the database concurrently <br><br>
**To Achieve 1 and 2 the database must maintain ACID property**
> ### 1.  Atomicity (A)
> Either all operations complete successfully or the transaction fails and the db is unchanged

> ### 2.  Consistency (C)
> The db state must be vallid after the transaction. All constraints mmust be satisfied.


> ### 3.  Isolation (I)
> Concurrent transactions must not affect each other.

>  ### 4.  Durability (D)
> Data written  by a successful transaction must be recorded in persistent storage.
>

## DBML Docs page
https://dbdocs.io/joshuamorka4/shiny_bank_project?schema=public&view=relationships&table=accounts


### Mock Database  -->
