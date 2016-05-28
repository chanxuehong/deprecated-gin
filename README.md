# Gin Web Framework

Gin is a web framework written in Golang. It features a martini-like API with much better performance, up to 40 times faster thanks to [httprouter](https://github.com/julienschmidt/httprouter). If you need performance and good productivity, you will love Gin.

## Start using it

1. Download and install it:

```sh
go get -u github.com/chanxuehong/gin
```

2. Import it in your code:

```go
import "github.com/chanxuehong/gin"
```

## API Examples

#### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
package main

import (
	"github.com/chanxuehong/gin"
)

func handler(ctx *gin.Context) {
	ctx.String(200, ctx.Request.Method)
}

func main() {
	router := gin.New()

	router.Get("/get", handler)
	router.Post("/post", handler)
	router.Put("/put", handler)
	router.Delete("/delete", handler)
	router.Patch("/patch", handler)
	router.Head("/head", handler)
	router.Options("/options", handler)

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
```

```go
package main

import (
	"github.com/chanxuehong/gin"
	"github.com/chanxuehong/gin/middleware"
)

func handler(ctx *gin.Context) {
	ctx.String(200, ctx.Request.Method)
}

func main() {
	router := gin.New()
	router.Use(middleware.Logger(), middleware.Recovery()) // use middleware

	router.Get("/get", handler)
	router.Post("/post", handler)
	router.Put("/put", handler)
	router.Delete("/delete", handler)
	router.Patch("/patch", handler)
	router.Head("/head", handler)
	router.Options("/options", handler)

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
```

#### Parameters in path

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.Get("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/join/
	router.Get("/user/:name/*action", func(ctx *gin.Context) {
		name := ctx.Param("name")
		action := ctx.Param("action")
		message := name + " is " + action
		ctx.String(http.StatusOK, message)
	})

	router.Run(":8080")
}
```

#### Querystring parameters

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.Get("/welcome", func(ctx *gin.Context) {
		firstname := ctx.DefaultQuery("firstname", "Guest")
		lastname := ctx.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		ctx.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.Run(":8080")
}
```

#### Urlencoded Form

```go
package main

import (
	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	router.Post("/form_post", func(ctx *gin.Context) {
		message := ctx.PostFormValue("message")
		nick := ctx.DefaultPostFormValue("nick", "anonymous")

		ctx.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	router.Run(":8080")
}
```

#### Another example: query + post form

```
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
```

```go
package main

import (
	"fmt"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	router.Post("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("id", "0")
		name := c.PostFormValue("name")
		message := c.PostFormValue("message")

		fmt.Println("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})

	router.Run(":8080")
}
```

```
id: 1234; page: 0; name: manu; message: this_is_great
```

#### Grouping routes

```go
package main

import (
	"github.com/chanxuehong/gin"
)

func handler1(ctx *gin.Context) { ctx.ResponseWriter.WriteString("1") }
func handler2(ctx *gin.Context) { ctx.ResponseWriter.WriteString("2") }
func handler3(ctx *gin.Context) { ctx.ResponseWriter.WriteString("3") }
func handler4(ctx *gin.Context) { ctx.ResponseWriter.WriteString("4") }
func handler5(ctx *gin.Context) { ctx.ResponseWriter.WriteString("5") }
func handler6(ctx *gin.Context) { ctx.ResponseWriter.WriteString("6") }

func main() {
	router := gin.New()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.Get("/login", handler1)
		v1.Get("/submit", handler2)
		v1.Get("/read", handler3)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.Get("/login", handler4)
		v2.Get("/submit", handler5)
		v2.Get("/read", handler6)
	}

	router.Run(":8080")
}
```

#### creating and using middleware

```go
package main

import (
	"github.com/chanxuehong/gin"
)

func middleware1(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("middleware1\r\n")
}
func middleware21(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("middleware21\r\n")
}
func middleware22(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("middleware22\r\n")
}
func handler1(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("handler1\r\n")
}
func handler2(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("handler2\r\n")
}

func groupMiddleware1(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group middleware1\r\n")
}
func groupMiddleware2(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group middleware2\r\n")
}
func groupMiddleware31(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group middleware31\r\n")
}
func groupMiddleware32(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group middleware32\r\n")
}
func groupHandler1(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group handler1\r\n")
}
func groupHandler2(ctx *gin.Context) {
	ctx.ResponseWriter.WriteString("group handler2\r\n")
}

func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(middleware1)

	// Per route middleware, you can add as many as you desire.
	r.Get("/test1", middleware21, handler1)
	r.Get("/test2", middleware22, handler2)

	group := r.Group("/user", groupMiddleware1) // group middleware
	group.Use(groupMiddleware2)                 // group middleware

	// Per route middleware, you can add as many as you desire.
	group.Get("/test1", groupMiddleware31, groupHandler1)
	group.Get("/test2", groupMiddleware32, groupHandler2)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Model binding and validation

To bind a request body into a type, use model binding. We currently support binding of JSON, XML.

Note that you need to set the corresponding binding tag on all fields you want to bind. For example, when binding from JSON, set `json:"fieldname"`.

You can also specify that specific fields are required. If a field is decorated with `validate:"required"` and has a empty value when binding, the current request will fail with an error.

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

// Binding from JSON
type Login struct {
	User     string `json:"user"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

func main() {
	router := gin.New()
	router.DefaultValidator(gin.DefaultValidator) // it not set, it does not validate the struct.

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.Post("/loginJSON", func(ctx *gin.Context) {
		var json Login
		if ctx.BindJSON(&json) == nil {
			if json.User == "manu" && json.Password == "123" {
				ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
```

#### XML and JSON rendering

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	r := gin.New()

	// gin.H is a shortcut for map[string]interface{}
	r.Get("/someJSON", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.Get("/moreJSON", func(ctx *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
		ctx.JSON(http.StatusOK, msg)
	})

	r.Get("/someXML", func(ctx *gin.Context) {
		ctx.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Serving static files

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	router.StaticAlias("/assets/", gin.Dir("./assets/"))  // /assets/1234.jpg --> ./assets/1234.jpg
	router.StaticRoot("/more_static/", http.Dir("/root")) // /more_static/1234.jpg --> /root/more_static/1234.jpg
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
```

#### Redirects

Issuing a HTTP redirect is easy:

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	router.Get("/test", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
```

Both internal and external locations are supported.


#### goroutines inside a middleware

When starting inside a middleware or handler, you **SHOULD NOT** use the original context inside it, you have to use a read-only copy.

```go
package main

import (
	"log"
	"time"

	"github.com/chanxuehong/gin"
)

func main() {
	r := gin.New()

	r.Get("/long_async", func(ctx *gin.Context) {
		// create copy to be used inside the goroutine
		ctxCopy := ctx.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note than you are using the copied context "ctxCopy", IMPORTANT
			log.Println("Done! in path " + ctxCopy.Request.URL.Path)
		}()
	})

	r.Get("/long_sync", func(ctx *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// since we are NOT using a goroutine, we do not have to copy the context
		log.Println("Done! in path " + ctx.Request.URL.Path)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Custom HTTP configuration

Use `http.ListenAndServe()` directly, like this:

```go
package main

import (
	"net/http"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()
	http.ListenAndServe(":8080", router)
}
```

or


```go
package main

import (
	"net/http"
	"time"

	"github.com/chanxuehong/gin"
)

func main() {
	router := gin.New()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
```
