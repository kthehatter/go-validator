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

## Example

```go

package main

import (
    "github.com/kthehatter/go-validator/validator"
    "github.com/kthehatter/go-validator/validator/ginadapter"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()
    options := []validator.ValidationOption{
        // Define validation options
    }
    r.POST("/user", ginadapter.Middleware(options), func(c *gin.Context) {
        // Handle request
    })
    r.Run()
}
```

### Step 7: Test the Package

Create a sample project to test importing and using your package:

sample-app/  
├── go.mod  
└── main.go

In `main.go`:

```go
package main

import (
    "github.com/kthehatter/go-validator/validator"
    "github.com/kthehatter/go-validator/validator/ginadapter"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()
    options := []validator.ValidationOption{
        // Define validation options
    }
    r.POST("/user", ginadapter.Middleware(options), func(c *gin.Context) {
        // Handle request
    })
    r.Run()
}
```