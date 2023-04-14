package recache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gobkc/to"
	"log"
	"reflect"
	"time"
)

var kindMap = map[reflect.Kind]func(value any) reflect.Value{
	reflect.Int: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[int](value))
	},
	reflect.Int32: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[int32](value))
	},
	reflect.Int64: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[int64](value))
	},
	reflect.Float32: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[float32](value))
	},
	reflect.Float64: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[float64](value))
	},
	reflect.String: func(value any) reflect.Value {
		return reflect.ValueOf(to.Any[string](value))
	},
}

var tls time.Duration = 1 * time.Hour

func SetTls(interval time.Duration) {
	tls = interval
}

func QueryDefault[result any](client *redis.Client, key string, def func() (defValue any)) *result {
	dest := new(result)
	v := client.Get(key).Val()
	kind := reflect.TypeOf(dest).Elem().Kind()
	isMap := kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice
	var value any
	if len(v) > 0 {
		value = v
	} else {
		value = def()
		if vk := reflect.TypeOf(value).Kind(); vk == reflect.Struct || vk == reflect.Map || vk == reflect.Slice {
			newVal, err := Marshal(value)
			if err != nil {
				log.Printf(`[WARNING]Marshal Redis Byte Error:%s\n`, err.Error())
			}
			value = string(newVal)
		}
	}
	if isMap {
		s := to.String(value)
		err := Unmarshal([]byte(s), dest)
		if err != nil {
			log.Printf(`[WARNING]Unmarshal Redis Byte Error:%s\n`, err.Error())
		}
		client.Set(key, s, tls)
	} else {
		convert, ok := kindMap[kind]
		if ok {
			newVal := convert(value)
			reflect.ValueOf(dest).Elem().Set(newVal)
		}
		if err := client.Set(key, to.String(*dest), tls).Err(); err != nil {
			log.Printf(`[WARNING]Set Redis Value Error: %v to %v %s\n`, key, *dest, err.Error())
		}
	}
	return dest
}

func SaveFlush(client *redis.Client, key string, def func() (defValue any)) {
	value := def()
	if vk := reflect.TypeOf(value).Kind(); vk == reflect.Struct || vk == reflect.Map || vk == reflect.Slice {
		newVal, err := json.Marshal(value)
		if err != nil {
			log.Printf(`[WARNING]Marshal Redis Byte Error:%s\n`, err.Error())
		}
		value = string(newVal)
	}
	s := to.String(value)
	if err := client.Set(key, s, tls).Err(); err != nil {
		log.Printf(`[WARNING]SaveFlush Error: %v to %v %s\n`, key, s, err.Error())
	}
}
