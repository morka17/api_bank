
### Postgres


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




### Mock Database 
