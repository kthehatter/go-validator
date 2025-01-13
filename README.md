# go-validator
  
A custom validator package for Go applications.
## Installation
 

```bash

go get github.com/kthehatter/go-validator/validator

```

## Usage

  Import the package in your Go code:

```go
import "github.com/kthehatter/go-validator/validator"
```

For Gin adapter:
```go

import "github.com/kthehatter/go-validator/validator/ginadapter"
```

## Test the Package
#### Step 1: Basic usage

```go
package main

import (

"github.com/gin-gonic/gin"

"github.com/kthehatter/go-validator/validator"

"github.com/kthehatter/go-validator/validator/ginadapter"

)

func main() {

    // Create a new Gin router

    r := gin.New()


    // Define validation rules

    validationOptions := []validator.ValidationOption{

        {
        
            Key: "username",
            
            IsOptional: false,
            
            Validators: []validator.Validator{
            
            validator.CreateValidator(validator.IsNotEmpty, "Username is required"),
            
            validator.CreateValidator(validator.IsAlphanumeric, "Username must be alphanumeric"),
        
            },
        
        },

        {
        
            Key: "email",
            
            IsOptional: false,
            
            Validators: []validator.Validator{
            
            validator.CreateValidator(validator.IsNotEmpty, "Email is required"),
            
            validator.CreateValidator(validator.IsEmail, "Invalid email address"),
        
            },
        
        },
        
        {
        
            Key: "password",
            
            IsOptional: false,
            
            Validators: []validator.Validator{
            
            validator.CreateValidator(validator.IsNotEmpty, "Password is required"),
            
            validator.CreateValidator(validator.MinLength(6), "Password must be at least 6 characters"),
        
            },
        
        },

    }

    // Apply the validation middleware to the POST /user endpoint

    r.POST("/user", ginadapter.Middleware(validationOptions), func(c *gin.Context) {

    // Retrieve the validated body from the context

    body := c.MustGet("validatedBody").(gin.H)

    // Process the data (for demonstration, just return it)

    c.JSON(200, gin.H{

        "message": "User created successfully",

        "data": body,

    })

    })  

    // Start the Gin server

    r.Run(":8080")

}

```

#### Step 2: Run the Application

1. **Build and run the application**:

```bash
go run main.go
```

2. **Verify the server is running**:
    
    - Open a web browser or use `curl` to send a POST request to `http://localhost:8080/user` with a JSON body.
        

#### Step 3: Test the Endpoint

**Test with Valid Data**:

```bash
curl -X POST http://localhost:8080/user -H "Content-Type: application/json" -d '{
    "username": "user123",
    "email": "user@example.com",
    "password": "password123"
}'
```

Expected response:

```json
{
    "message": "User created successfully",
    "data": {
        "username": "user123",
        "email": "user@example.com",
        "password": "password123"
    }
}
```

**Test with Invalid Data**:

```bash
curl -X POST http://localhost:8080/user -H "Content-Type: application/json" -d '{
    "username": "",
    "email": "invalid-email",
    "password": "pass"
}'
```

Expected response:

```json
{
    "message": "Username is required"
}
```
