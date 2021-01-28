package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

func performForEmail(name string, threads int, body map[string]interface{}) (duration int64) {
	wg := &sync.WaitGroup{}
	startTime := time.Now()
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func(wg *sync.WaitGroup, name string, body map[string]interface{}) {
			defer wg.Done()

			client := resty.New()
			for j := 0; j < (50 / threads); j++ {
				body["email"] = "mkm" + strconv.Itoa(i) + strconv.Itoa(j) + "@gmail.com"
				_, err := client.R().SetBody(body).SetHeader("Content-Type", "application/json").Post("http://ec2co-ecsel-1jni8txjq1cve-122992530.us-west-2.elb.amazonaws.com:8001/subscribe")
				if err != nil {
					log.Println("[", name, "]", err)
					continue
				}
			}
		}(wg, name, body)
	}

	wg.Wait()

	return time.Now().Sub(startTime).Milliseconds()
}

func performForText(name string, threads int, body map[string]interface{}) (duration int64) {
	wg := &sync.WaitGroup{}
	startTime := time.Now()
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func(wg *sync.WaitGroup, name string, body map[string]interface{}) {
			defer wg.Done()

			client := resty.New()
			for j := 0; j < (50 / threads); j++ {
				if i < 10 {
					body["phoneNumber"] = "+100000000" + strconv.Itoa(i) + "0" + strconv.Itoa(j)
				} else {
					body["phoneNumber"] = "+100000000" + strconv.Itoa(i) + strconv.Itoa(j)
				}
				_, err := client.R().SetBody(body).SetHeader("Content-Type", "application/json").Post("http://ec2co-ecsel-1jni8txjq1cve-122992530.us-west-2.elb.amazonaws.com:8002/subscribe")
				if err != nil {
					log.Println("[", name, "]", err)
					continue
				}
			}
		}(wg, name, body)
	}

	wg.Wait()

	return time.Now().Sub(startTime).Milliseconds()
}

func main() {
	log.Println("[performEmail] with 1 thread", performForEmail("test", 1, map[string]interface{}{"email": "mkm@gmail.com", "frequency": 864000000}))

	log.Println("[performEmail] with 3 threads", performForEmail("test", 3, map[string]interface{}{"email": "mkm@gmail.com", "frequency": 864000000}))

	log.Println("[performEmail] with 5 threads", performForEmail("test", 5, map[string]interface{}{"email": "mkm@gmail.com", "frequency": 864000000}))

	log.Println("[performForText] with 1 thread", performForText("test", 1, map[string]interface{}{"phoneNumber": "+17783239700", "frequency": 864000000}))

	log.Println("[performForText] with 3 threads", performForText("test", 3, map[string]interface{}{"phoneNumber": "+17783239700", "frequency": 864000000}))

	log.Println("[performForText] with 5 threads", performForText("test", 5, map[string]interface{}{"phoneNumber": "+17783239700", "frequency": 864000000}))

}
