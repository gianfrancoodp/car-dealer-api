# Car Dealer REST API

## Technical Information
This REST API was developed with the following technologies:
* `Language`: Golang 1.20.3
* `IDE`: GoLand by NetBrains.
* `Data Base`: MongoDB. 
* `Data Base Management`: MongoDB Atlas.
* `Web API Framework`: Fiber v2


## Setting Up a Development Environment

```bash
# Get the code
git clone https://https://github.com/gianfrancoodp/car-dealer-api Factly
cd Factly
```

## Starting the Web Server

The web server can be started as shown below (assuming Go 1.20.3 is installed). By default, this application listens for
HTTP connections on port 6000, so point your client at
`localhost:6000`.

 ```bash
 go run main.go
```

## Consuming the REST API Endpoints

The REST API Endpoints could be consumed in *Postman*. The collection file is in the "*routes*" folder, you can open this file in your Postman application:

![Postman-Collection](https://github.com/gianfrancoodp/car-dealer-api/blob/master/doc/Postman_collection.png)

The collection file has the following endpoints:

![Endpoints.png](https://github.com/gianfrancoodp/car-dealer-api/blob/master/doc/endpoints.png)

* `POST:` create a new car
```
http://localhost:6000/car
```
* `GET ALL:` get all cars
```
http://localhost:6000/cars/
```
* `GET:` get a car with its ID number
```
http://localhost:6000/car/carId
```
* `PUT:` update a car
```
http://localhost:6000/car/carId
```
* `DELETE:` delete a car
```
http://localhost:6000/car/carId
```

## REST API Structure

![rest-api-structure.png](https://github.com/gianfrancoodp/car-dealer-api/blob/master/doc/rest-api-structure.png)

## References
1. Non-Relational Databases Naming Convention: https://www.coding-guidelines.lftechnology.com/docs/nosql/documentdb/document-db-naming-convention/).
2. Project Naming Convention: https://blog.devgenius.io/golang-naming-conventions-72bbaf84e959
3. Naming Convention de Packages: https://go.dev/blog/package-names
4. Best practices to write commits in Git: https://www.freecodecamp.org/news/writing-good-commit-messages-a-practical-guide/
5. Git Branching Naming Convention Best Practices: https://codingsight.com/git-branching-naming-convention-best-practices/
6. Golang Fundamentals: https://www.mindbowser.com/golang-language-fundamentals/
7. REST API Fundamentals: https://arunrajeevan.medium.com/fundamentals-of-rest-api-design-d9c425c1b0f6