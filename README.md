# Pension-ISA


## PostgreSQL

- **ACID compliance** Handling of financial data requires ACID compliance so SQL is best vs noSQL
    - Atomicity: 
    - Consistency:
    - Isolation:
    - Durability:
- **Query Complexity** Support for complex queries may be required
- Well supported by AWS via both AWS RDS for PostgreSQL and Aurora Managed DB
- **AWS DBMS Support** Easy migration to Aurora if managed service is desired via AWS DBMS (Database migration service) or we can opt for lower costs with AWS RDS
- **Data Availability** Both provide multi-AZ deployment options for higher data availabilty and durability
- **Database Seeding** Test data seeding provided for development environment only
- **Generated UUIDs** Randomly generated UUIDs provided by uuid-ossp. The preference for this over sequenced IDs is to allow for the potential of database sharding in a distributed sysem. This also has the added benefit of enhanced security as IDs cannot be easily iterated through.
- **Data Normalisation** 

## React
- **Scope** Minimalist frontend implementation. More focus has been given to project structure, design choices and backend functionality
- **Vite** Vite used for quick frontend deployment
- **Typescript** Typescript used on the frontend rather than Javascript

## Go
- We'll use Go in the backend
- Environment variables

## Testing
- Jest + React Testing Library
- Go unit tests
- DB integration testing

## API Design
- Versioning Strategy
- Middleware for things like auth

## Containerisation
-Docker Compose

## Microservices
- Domain Driven Design
- Event driven
- Env vars for service discovery

## Security
- **Secrets** For now I have stored them in an .env file. This is not ideal but is suitable for the current implementation. For our actual environments we would want to make use of a dedicated secret storage that supports secret rotation such as AWS Secrets Manager of Hashicorp Vault.

- Input validation
- Secure headers
- Rate limiting
- Secure DB connections

## Monitoring and Metrics ##
- Prometheus

## Potential AWS Integration ##


