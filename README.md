# Introduction
lock-free / mutex-free and thread-safe promise written in Go
Supports 1.18 generics

# Usage Example

## Await

```go
func main() {
	ip, err := getIP().Await()
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("my ip is %s", ip)
	}
}

func getIP() *promise.Promise[string] {
	var prom promise.Promise[string]
	go func() {
		resp, err := http.Get("https://httpbin.org/ip")
		if err != nil {
			prom.Error(err)
			return
		}
		type Response struct {
			Origin string `json:"origin"`
		}
		var response Response

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			prom.Error(err)
			return
		}
		prom.Resolve(response.Origin)
	}()
	return &prom
}
```

## Then / Catch / Finally

```go
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
	getIP().Then(func(ip string) {
        log.Printf("my ip is %s", ip)
    }).Catch(func(err error) {
        log.Printf("error: %v", err)
    }).Finally(wg.Done)
    // wait until promise finishes
    wg.Wait()
}
```

## Chaining multiple promises

```go
package main

func main() {
	loc, err := getCurrentLocation().Await()
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("im in %s, %s", loc.City, loc.Country)
	}
}

// chain multiple promise
// getIP -> getIPLocation
func getCurrentLocation() *promise.Promise[IPLocation] {
	var prom promise.Promise[IPLocation]
	getIP().Then(func(ip string) {
		// forward result
		getIPLocation(ip).
			Then(prom.Resolve).
			Catch(prom.Error)
	}).Catch(prom.Error)
	return &prom
}

func getIPLocation(ip string) *promise.Promise[IPLocation] {
	var prom promise.Promise[IPLocation]
	go func() {
		resp, err := http.Get("http://ip-api.com/json/" + ip)
		if err != nil {
			prom.Error(err)
			return
		}
		var loc IPLocation
		err = json.NewDecoder(resp.Body).Decode(&loc)
		if err != nil {
			prom.Error(err)
			return
		}
		prom.Resolve(loc)
	}()
	return &prom
}

type IPLocation struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}
```
