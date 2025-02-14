# Pension-ISA
This solution is mostly focused on the backend, database, data structures, API and overall repo structure.

Any features listed below that are not **bold** are features that were either planned and not added or simply considerations for improvement.

Attempts were made to try and adhere to some DDD principles within the limited scope of the project. Care was given to ensure clear separation of concerns on the backend.

I have tried to provide a breadth of basic implementations. Of course, many of the features can be expanded upon.

Particular points of note are tagged throughout the codebase with the **Note** tag.
Planned improvements are tagged through the codebase with the **TODO** tag.

## How To Run
The database is containerised and can be run in docker via **docker-compose up  --build**
This will also seed the database with some inital data.

The backend can be started via the Makefile using **make run** if it has not also been added to the docker-compose file at the time of viewing.

Postman and pgAdmin were used to test the APIs and DB. Please note that some API routes are protected. Steps to bypass this in Postman can be found under the Security section.

## Database
- **PostgreSQL:** PostgreSQL was chosen for the reasons outlined below
- **ACID compliance:** Handling of financial data requires ACID compliance so SQL is best vs noSQL
    - Atomicity: Ensures financial transactions are all-or-nothing
    - Consistency: Maintains data integrity and business rules
    - Isolation: Prevents interference between concurrent transactions
    - Durability: Ensures committed transactions are permanent
- **Query Complexity:** Support for complex queries may be required
- **AWS Support:** Well supported by AWS via both AWS RDS for PostgreSQL and Aurora Managed DB
- **AWS DBMS Support:** Easy migration to Aurora if managed service is desired via AWS DBMS (Database migration service) or we can opt for lower costs with AWS RDS
- **Data Availability:** Both provide multi-AZ deployment options for higher data availabilty and durability
- **Database Seeding:** Test data seeding provided for development environment only
- **Generated UUIDs:** Randomly generated UUIDs provided by uuid-ossp. The preference for this over sequenced IDs is to allow for the potential of database sharding in a distributed sysem. This also has the added benefit of enhanced security as IDs cannot be easily iterated through.
- **Data Normalisation:** Currently the database adheres to at least BCNF
- **Materialised Views:** A materialised view is used to determine the total investment that a user has made to a particular fund
- **Indexing:** Indexing on certain keys to provide faster lookups

## Backend
- **Go:** We use Go as our language of choice for the backend.
- **Chi:** Chi was selected for some convenience with routing. It was specifically chosen as it is lightweight and uses only stdlib + net/http
- **Separation of Concerns:**
    - Handler layer: Only deals with HTTP concerns
    - Service Layer: Contains any business logic between Handler and Respository layers
    - Repository Layer: Handles data access
- **Repository Pattern:** This pattern allows easy switch out of database choices
- **Helpers and Middleware:** Helper methods and middleware provide shared and reusable functionality
- **Pagination:** Custom pagination middleware to allow configurable page sizes of returned data
- **Materialized View Refresh:** We trigger a refresh on the materialized view after each investment
- **Transaction Rollbacks:** If we fail to update the materialized view we rollback the transaction to ensure data consistency (Likely not ideal behaviour in the real world but it's pretty neat and serves a good example of the atomicity required in financial transations)
- **Environment variables** Environment variables set .env file and read into config
- **Fund Limit:** Customers limited to investing in one fund. I chose to do this in the backend code rather than put limitations within the DB as it is easier to switch out at a later date if this limitation is removed.
- **Database Healthcheck:** A Go routine that periodically pings the database and gathers some statistics

## API Design
- **Versioning:** Versioning implemented from the start
- **Deprecation:** Deprecation functionality added but not actively used
- **Pagination:** Middleware added for pagination of List requests
- **Auth:** Auth middleware used to restrict access to certain API calls 

## Containerisation
- **Docker Compose:** Docker is used to run the Postgres DB in a container. If time permits I will also add the backend and frontend to containers

## Frontend (ABANDONED)
- **Abandoned:** Started then swiftly abandoned once I had clarification on the scope of this project. I am leaving it in as I will likely come back and work on this at a later date.
- **Scope:** Minimalist frontend implementation. More focus has been given to project structure, design choices and backend functionality
- **Vite:** Vite used for quick frontend deployment
- **Typescript:** React and Typescript used on the frontend

## Security
- **Secrets:** For now I have stored them in text files to be read into docker-compose. This is not ideal but is suitable for the current implementation. For our actual environments we would want to make use of a dedicated secret storage that supports secret rotation such as AWS Secrets Manager or Hashicorp Vault.
- **Secret Generation:** A script is used to generate secrets (found under /scripts)
- **JWT Auth:** Authentication used to protect certain API routes. To test these routes via postman you can use the following settings under Authorization in your Postman call:
    - Auth Type: JWT Bearer
    - Add JWT to: Request Header
    - Algorithm: HS256
    - Secret: secret
- **Input validation:** Simple input validation. Could be greatly expanded upon
- Rate limiting
- Query timeouts
- Query retries with exponential backoff
- Proper use of interfaces looser coupling

## Testing
- **Unit Tests** Basic happy path unit tests implemented for the repository files. Obviously testing should be greatly expanded upon in an ideal situation. Testing makes use of mocks and Testify.
- Jest + React Testing Library
- DB integration testing

## Microservices
- Domain Driven Design
- Event driven

## Monitoring and Metrics ##
- **DB Metrics:** Currently gather database connection metrics

## Potential AWS Integration ##


