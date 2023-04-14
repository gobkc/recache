# About recache
A redis cache library

### Contributing
You can commit PR to this repository

### How to get it?
````
go get -u github.com/gobkc/recache
````

### Quick start
````
package main

import (
	"github.com/gobkc/recache"
	"log"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	def := recache.QueryDefault[map[string]interface{}](rdb, "test3", func() (defValue any) {
		return map[string]interface{}{
			"aaa": 123, "bbb": "456",
		}
	})
	recache.SaveFlush(rdb, "test5", func() (defValue any) {
		return map[string]interface{}{
			"aaa": 123, "bbb": "456",
		}
	})
}
````

### License
© Gobkc, 2023~time.Now

Released under the Apache License