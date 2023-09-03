package main

import (
	"fmt"
	"time"

	"mycache/m/airport"
	"mycache/m/cache"
)

func readAside(cache *cache.Cache[airport.Airport], storage airport.Storage, key string, ttl time.Duration) (airport.Airport, bool) {

	val, exists := cache.Get(key)
	if exists {
		fmt.Println("Cache hit")
		return val, exists
	}

	fmt.Println("Cache miss, read from storage")
	val, exists = storage.Read(key)
	if !exists {
		return airport.Airport{}, exists
	}
	// write to cache
	cache.Set(key, val, ttl)
	return val, true
}

func main() {

	c := cache.New[airport.Airport]()

	storage := airport.Storage{
		Datas: [][2]string{
			{"HAN", "Noi Bai"}, {"SGN", "Tan Son Nhat"}, {"DAD", "Da Nang airport"}, {"HPH", "Cat Bi"},
		},
	}

	val, exists := readAside(c, storage, "HAN", time.Second)
	fmt.Println(val, exists)

	val, exists = readAside(c, storage, "HAN", time.Millisecond)
	fmt.Println(val, exists)

	time.Sleep(2 * time.Second)

	val, exists = readAside(c, storage, "HAN", time.Second)
	fmt.Println(val, exists)

	val, exists = readAside(c, storage, "ABC", time.Millisecond)
	fmt.Println(val, exists)
}
