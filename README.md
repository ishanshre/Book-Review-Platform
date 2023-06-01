# Book-Review-Platform
A simple server side rendering web application


## Required
    - Golang: https://go.dev/doc/install
    - Golang-Migrate: https://github.com/golang-migrate/migrate
    - Docker(Optional): https://docs.docker.com/engine/install/
      - For running database. Alternatively you can install database to your machine or use cloud
    - Makefile:
      - It is a script used for compiling or building binaries
## Create .env file 
```
m_db_username=your database username
m_db_password=your database password
m_db_dbname=your database name
postgres="user=database_username password=database_password dbname=name_of_database sslmode=disable"
test="user=database_username password=database_password dbname=name_of_database sslmode=disable"

```