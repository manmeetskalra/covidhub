package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vonage/vonage-go-sdk"
)

var (
	Layout                = "2006-01-02 15:04:05"
	CovidHubAPITimeLayout = "2006-01-02T00:00:00Z"
	MySQLHost             = ""
	MySQLPort             = "29291"
	MySQLUser             = "root"
	MySQLPassword         = ""
	DefaultNumber         = ""
	NEXMO_API_KEY         = ""
	NEXMO_API_SECRET      = ""
)

type PhonePayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Frequency   int64  `json:"frequency,omitempty"`
	Country     string `json:"countries,omitempty"`
}

type TextSender struct {
	DB             *sql.DB
	processTextsCh chan PhonePayload
	smsClient      *vonage.SMSClient
}

func (em *TextSender) Run() {
	em.processTextsCh = make(chan PhonePayload)
	go em.sendTextAndUpdateDueTime()

	for {
		texts, err := em.getTextsDueCurrently()
		if err != nil {
			log.Println("failed to getTextsDueCurrently %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, text := range texts {
			em.processTextsCh <- text
		}

		// sleep for 1 second always
		time.Sleep(1 * time.Second)
	}
}

func (em *TextSender) sendTextAndUpdateDueTime() {
	for {
		select {
		case payload := <-em.processTextsCh:
			to := time.Now()
			from := time.Now().Add(-(time.Duration(payload.Frequency)) * time.Second)
			covidAPIResults, err := retreiveDataFromCovidApi(payload.Country, from, to)
			if err != nil {
				log.Println(err)
				continue
			}

			_, err = em.smsClient.Send(DefaultNumber, payload.PhoneNumber, covidApiDataToText(covidAPIResults), vonage.SMSOpts{})
			if err != nil {
				log.Println(err)
			}
		default:
		}
	}
}

func (em *TextSender) getTextsDueCurrently() ([]PhonePayload, error) {
	texts := []PhonePayload{}

	currentDueTime := time.Now().Format(Layout)
	getPhoneNumbersDueQuery := fmt.Sprintf("SELECT phoneNumber, dueTime, frequency, country from covidhub.texts where dueTime = '%s'", currentDueTime)
	results, err := em.DB.Query(getPhoneNumbersDueQuery)
	if err != nil {
		log.Printf("failed to get texts due at: %s", currentDueTime)
		return texts, err
	}

	var (
		phoneNumber, country string
		frequency            int64
		dueTime              time.Time
	)

	for results.Next() {
		err := results.Scan(&phoneNumber, &dueTime, &frequency, &country)
		if err != nil {
			log.Printf("error in retreiving texts due currently: %v", err)
			return texts, err
		}
		texts = append(texts, PhonePayload{PhoneNumber: phoneNumber, Frequency: frequency, Country: country})
	}

	return texts, nil
}

type CovidCountryInfo struct {
	Country   string `json:"country"`
	Confirmed int64  `json:"Confirmed"`
	Recovered int64  `json:"Recovered"`
	Deaths    int64  `json:"Deaths"`
	Date      string `json:"Date"`
}

func retreiveDataFromCovidApi(country string, from time.Time, to time.Time) ([]CovidCountryInfo, error) {
	log.Println("[DEBUG] retrieveDataFromCovidApi ----------------------------------")
	results := []CovidCountryInfo{}
	uri := fmt.Sprintf("https://api.covid19api.com/country/%s", country)
	log.Println(uri)
	resp, err := http.Get(uri)
	if err != nil {
		return results, fmt.Errorf("error in retreiving covid-api data %v", err)
	}
	defer resp.Body.Close()

	localResults := []CovidCountryInfo{}
	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	err = json.Unmarshal(body, &localResults)
	if err != nil {
		return results, fmt.Errorf("error in retreiving covid-api data %v", err)
	}

	for _, localResult := range localResults {
		previousTime, err := time.Parse(CovidHubAPITimeLayout, localResult.Date)
		if err != nil {
			return results, fmt.Errorf("error in retreiving covid-api data %v", err)
		}
		if from.Before(previousTime) && previousTime.Before(to) {
			results = append(results, localResult)
		}
	}

	return results, nil
}

func covidApiDataToText(data []CovidCountryInfo) string {
	buf := new(bytes.Buffer)
	for _, d := range data {
		text := fmt.Sprintf("%s has %d new confirmed, %d new deaths, and %d newly recovered.",
			d.Country,
			d.Confirmed,
			d.Deaths,
			d.Recovered,
		)
		fmt.Fprintln(buf, text)
	}
	return buf.String()
}

type Server struct {
	DB        *sql.DB
	smsClient *vonage.SMSClient
}

func (s *Server) subscribe(w http.ResponseWriter, req *http.Request) {
	var ep PhonePayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.processCreatePhoneNumber(ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) processCreatePhoneNumber(payload PhonePayload) error {
	var (
		phoneNumber string
		frequency   int64
		dueTime     time.Time
	)
	getPhoneNumberQuery := fmt.Sprintf("SELECT phoneNumber, dueTime, frequency from covidhub.texts where phoneNumber = '%s'", payload.PhoneNumber)
	result := s.DB.QueryRow(getPhoneNumberQuery)
	err := result.Scan(&phoneNumber, &dueTime, &frequency)
	if err != nil {
		dueTime := time.Now().Add(time.Second * time.Duration(payload.Frequency)).Format(Layout)
		createPhoneNumberQuery := fmt.Sprintf("INSERT INTO covidhub.texts (phoneNumber, frequency, dueTime, country) VALUES ('%s', '%d', '%s', '%s')",
			payload.PhoneNumber,
			payload.Frequency,
			dueTime,
			payload.Country,
		)
		_, err := s.DB.Query(createPhoneNumberQuery)
		if err != nil {
			log.Printf("error in creating phoneNumber for: %s due to %v", payload.PhoneNumber, err)
			return err
		}
	} else {
		newDueTime := dueTime.Add(time.Second * time.Duration(payload.Frequency-frequency)).Format(Layout)
		createPhoneNumberQuery := fmt.Sprintf("UPDATE covidhub.texts SET phoneNumber = '%s', frequency = '%d', dueTime = '%s', country = '%s' WHERE phoneNumber = '%s'",
			payload.PhoneNumber,
			payload.Frequency,
			newDueTime,
			payload.Country,
			payload.PhoneNumber,
		)

		_, err = s.DB.Query(createPhoneNumberQuery)
		if err != nil {
			log.Printf("error in updating frequency, phoneNumber for: %s due to %v", payload.PhoneNumber, err)
			return err
		}
	}
	return nil
}

func (s *Server) processDeletingPhoneNumber(payload PhonePayload) error {
	deletePhoneNumberQuery := fmt.Sprintf("DELETE FROM covidhub.texts WHERE phoneNumber = '%s'", payload.PhoneNumber)
	_, err := s.DB.Query(deletePhoneNumberQuery)
	if err != nil {
		log.Printf("error in deleting phoneNumber for: %s due to %v", payload.PhoneNumber, err)
		return err
	}

	return nil
}

func (s *Server) unsubscribe(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		http.Error(w, "404 not found.", http.StatusBadRequest)
	}
	var ep PhonePayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.processDeletingPhoneNumber(ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) subscribeNow(w http.ResponseWriter, req *http.Request) {
	var ep PhonePayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := retreiveDataFromCovidApi(ep.Country, time.Now().AddDate(0, 0, -3), time.Now())
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to retreive covid-api data", http.StatusBadRequest)
		return
	}

	log.Println(covidApiDataToText(data))

	_, err = s.smsClient.Send(DefaultNumber, ep.PhoneNumber, covidApiDataToText(data), vonage.SMSOpts{})
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to send text", http.StatusBadRequest)
		return
	}

}

func (s *Server) Run() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		Debug:          true,
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/subscribe", s.subscribe)
	mux.HandleFunc("/unsubscribe", s.unsubscribe)
	mux.HandleFunc("/subscribe/now", s.subscribeNow)

	log.Println("text server listening at port 8002")
	handler := c.Handler(mux)
	http.ListenAndServe(":8002", handler)
}

func init() {
	if MySQLHost = os.Getenv("MYSQL_HOST"); MySQLHost == "" {
		log.Fatalf("MYSQL_HOST not set")
	}
	if MySQLPort = os.Getenv("MYSQL_PORT"); MySQLPort == "" {
		log.Fatalf("MYSQL_PORT not set")
	}
	if MySQLUser = os.Getenv("MYSQL_USER"); MySQLUser == "" {
		log.Fatalf("MYSQL_USER not set")
	}

	MySQLPassword = os.Getenv("MYSQL_PASSWORD")

	if DefaultNumber = os.Getenv("DEFAULT_NUMBER"); DefaultNumber == "" {
		log.Fatalf("DEFAULT_NUMBER not set")
	}
	if NEXMO_API_KEY = os.Getenv("NEXMO_API_KEY"); NEXMO_API_KEY == "" {
		log.Fatalf("NEXMO_API_KEY not set")
	}
	if NEXMO_API_SECRET = os.Getenv("NEXMO_API_SECRET"); NEXMO_API_SECRET == "" {
		log.Fatalf("NEXMO_API_SECRET not set")
	}
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?multiStatements=true&parseTime=true",
		MySQLUser,
		MySQLPassword,
		MySQLHost,
		MySQLPort,
	)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("cannot establish connection to the mysql database: %v", err)
	}

	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("cannot reach to db at %s due to %v", dsn, err)
	}

	auth := vonage.CreateAuthFromKeySecret(NEXMO_API_KEY, NEXMO_API_SECRET)
	textSender := &TextSender{
		DB:        sqlDB,
		smsClient: vonage.NewSMSClient(auth),
	}
	go textSender.Run()

	server := &Server{
		DB:        sqlDB,
		smsClient: vonage.NewSMSClient(auth),
	}
	server.Run()
}
