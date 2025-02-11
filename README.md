# Pension-ISA
This solution is mostly focused on the backend, database, data structures, API and overall repo structure.

Attempts were made to try and adhere to some DDD principles within the limited scope of the project. Care was given to ensure clear separation of concerns on the backend.

Postman and pgAdmin were used to test the APIs and DB.

Particular points of note are tagged throughout the codebase with the **Note** tag.
Planned improvements are tagged through the codebase with the **TODO** tag.

## Database
- **PostgreSQL** PostgreSQL was chosen for the reasons outlined below
- **ACID compliance** Handling of financial data requires ACID compliance so SQL is best vs noSQL
    - Atomicity: 
    - Consistency:
    - Isolation:
    - Durability:
- **Query Complexity** Support for complex queries may be required
- **AWS Support** Well supported by AWS via both AWS RDS for PostgreSQL and Aurora Managed DB
- **AWS DBMS Support** Easy migration to Aurora if managed service is desired via AWS DBMS (Database migration service) or we can opt for lower costs with AWS RDS
- **Data Availability** Both provide multi-AZ deployment options for higher data availabilty and durability
- **Database Seeding** Test data seeding provided for development environment only
- **Generated UUIDs** Randomly generated UUIDs provided by uuid-ossp. The preference for this over sequenced IDs is to allow for the potential of database sharding in a distributed sysem. This also has the added benefit of enhanced security as IDs cannot be easily iterated through.
- **Data Normalisation** Currently database is at 3NF. If time allows after initial implementation we will increase this
- **Materialised Views** A materialised view is used to determine the total investment that a user has made to a particular fund
- **Indexing** Indexing on certain keys to provide faster lookups

## Backend
- **Go** We use Go as our language of choice for the backend.
- **Chi** Chi was selected for some convenience with routing. It was specifically chosen as it is lightweight and uses only stdlib + net/http
- **Separation of Concerns**
    - Handler layer: Only deals with HTTP concerns
    - Service Layer: Contains any business logic between Handler and Respository layers
    - Repository Layer: Handles data access
- **Repository Pattern** This pattern allows easy switch out of database choices
- **Helpers and Middleware** Helper methods and middleware provide shared and reusable functionality
- **Pagination** Custom pagination middleware to allow configurable page sizes of returned data
- **Materialized View Refresh** We trigger a refresh on the materialized view after each investment
- **Transaction Rollbacks** If we fail to update the materialized view we rollback the transaction to ensure data consistency (Likely not ideal behaviour in the real world but it's pretty neat)
- Environment variables

## API Design
- **Versioning** Versioning implemented from the start
- **Deprecation** Deprecation functionality added but not actively used
- **Pagination** Middleware added for pagination of List requests
- Auth 

## Containerisation
- **Docker Compose** Docker is used to run the Postgres DB in a container. If time permits I will also add the backend and frontend to containers

## Frontend
- **Scope** Minimalist frontend implementation. More focus has been given to project structure, design choices and backend functionality
- **Vite** Vite used for quick frontend deployment
- **Typescript** React and Typescript used on the frontend

## Security
- **Secrets** For now I have stored them in text files to be read into docker-compose. This is not ideal but is suitable for the current implementation. For our actual environments we would want to make use of a dedicated secret storage that supports secret rotation such as AWS Secrets Manager of Hashicorp Vault.
- **Secret Generation** A script is used to generate secrets

- Input validation
- Secure headers
- Rate limiting
- Secure DB connections

## Testing
- Go unit tests
- Jest + React Testing Library
- DB integration testing

## Microservices
- Domain Driven Design
- Event driven
- Env vars for service discovery

## Monitoring and Metrics ##

## Potential AWS Integration ##


