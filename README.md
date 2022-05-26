# User Record REST API

## Individual Entry stored in the databsae

```json
    {
        "_id": "628e5ac214322b31dac15601",
        "name": "Habeeb Ullah"
        "dob": "11 feb 2002",
        "address": "8194 Euclid City, Hyderabad",
        "description": "Backend Developer",
        "createdAt": "2022-05-25 22:05:14.426684 +0530 IST m=+55.164231301",
    }
```

## Endpoints Description

### Get All Users

```
    URL - *http://localhost:6000/users*
    Method - GET
```

### Get User By ID

```
    URL - *http://localhost:6000/user/<User-ID>*
    Method - GET
```

### Create a new User

```
    URL - *http://localhost:6000/user*
    Method - POST
    Request Header - (Content-Type : application/json)
    Request Body - 

    {
        "name": "Habeeb Ullah"
        "dob": "11 feb 2002",
        "address": "8194 Euclid City, Hyderabad",
        "description": "Backend Developer",
        "createdat": "2022-05-25 22:05:14.426684 +0530 IST m=+55.164231301",
    }

```


### Update User

```
    URL - *http://localhost:6000/user/<User-ID>*
    Method - PUT
    Request Header - (Content-Type : application/json)
    Request Body - 

    {
        "name": "Habeeb Ullah"
        "dob": "11 feb 2002",
        "address": "8194 Euclid City, Hyderabad",
        "description": "Backend Developer",
    }
```

### Delete User

```
    URL - *http://localhost:6000/user/<User-ID>*
    Method - DELETE
```

## Test Driven Development Description

To run all the unit test cases, please do the following -

1. `go test -v`

## Hope everything works. Thank you.
