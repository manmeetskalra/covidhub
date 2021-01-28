package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	gomail "gopkg.in/mail.v2"
)

var (
	Layout                = "2006-01-02 15:04:05"
	CovidHubAPITimeLayout = "2006-01-02T00:00:00Z"
	MySQLHost             = ""
	MySQLPort             = "29291"
	MySQLUser             = "root"
	MySQLPassword         = ""
	DefaultEmail          = ""
	DefaultEmailPassword  = ""
)

type EmailPayload struct {
	Email     string `json:"email"`
	Frequency int64  `json:"frequency,omitempty"`
	Countries string `json:"countries,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (ep EmailPayload) countriesFromString() []string {
	return strings.Split(ep.Countries, "|")
}

type EmailSender struct {
	DB             *sql.DB
	processEmailCh chan EmailPayload
}

func (em *EmailSender) Run() {
	em.processEmailCh = make(chan EmailPayload)
	go em.sendEmailAndUpdateDueTime()

	for {
		emails, err := em.getEmailsDueCurrently()
		if err != nil {
			log.Println("failed to getEmailsDueCurrently %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, email := range emails {
			em.processEmailCh <- email
		}

		// sleep for 1 second
		time.Sleep(1 * time.Second)
	}
}

func (em *EmailSender) sendEmailAndUpdateDueTime() {
	for {
		select {
		case payload := <-em.processEmailCh:
			to := time.Now()
			from := time.Now().Add(-(time.Duration(payload.Frequency)) * time.Second)
			covidAPIResults, err := retrieveDataFromCovidApi(payload.countriesFromString(), payload.Type, from, to)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := sendEmail(payload.Email, covidAPIResults); err != nil {
				log.Println(err)
			}

			updatedDueTime := time.Now().Add(time.Second * time.Duration(payload.Frequency)).Format(Layout)
			updateDueTimeQuery := fmt.Sprintf("UPDATE covidhub.emails SET dueTime = '%s' WHERE email = '%s'", updatedDueTime, payload.Email)
			_, err = em.DB.Query(updateDueTimeQuery)
			if err != nil {
				log.Println("failed to update dueTime for email %s due to %v", payload.Email, err)
			}

		default:
		}
	}
}

func sendEmail(to string, covidAPIResults []CovidCountryInfo) error {
	var body bytes.Buffer
	t, err := template.ParseFiles("emailtemplate.html")
	if err != nil {
		log.Println("failed to parse covidhub template %v", err)
	}
	t.Execute(&body, covidAPIResults)

	m := gomail.NewMessage()
	m.SetHeader("From", DefaultEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "CovidHub Updates")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, DefaultEmail, DefaultEmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s due to %v", to, err)
	}

	return nil
}

type CovidCountryInfo struct {
	Country string `json:"country"`
	Cases   int64  `json:"Cases"`
	Status  string `json:"Status"`
	Date    string `json:"Date"`
}

func retrieveDataFromCovidApi(countries []string, infoType string, from time.Time, to time.Time) ([]CovidCountryInfo, error) {
	results := []CovidCountryInfo{}
	log.Println("[DEBUG] retrieveDataFromCovidApi ----------------------------------")
	for _, country := range countries {
		uri := fmt.Sprintf("https://api.covid19api.com/country/%s/status/%s", country, infoType)
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
	}
	return results, nil
}

func (em *EmailSender) getEmailsDueCurrently() ([]EmailPayload, error) {
	emails := []EmailPayload{}

	currentDueTime := time.Now().Format(Layout)
	getEmailsDueQuery := fmt.Sprintf("SELECT email, dueTime, frequency, countries, informationType from covidhub.emails where dueTime = '%s'", currentDueTime)
	results, err := em.DB.Query(getEmailsDueQuery)
	if err != nil {
		log.Printf("failed to get emails due at: %s", currentDueTime)
		return emails, err
	}

	var (
		email, infoType, countries string
		frequency                  int64
		dueTime                    time.Time
	)

	for results.Next() {
		err := results.Scan(&email, &dueTime, &frequency, &countries, &infoType)
		if err != nil {
			log.Printf("error in retreiving emails due currently: %v", err)
			return emails, err
		}
		emails = append(emails, EmailPayload{
			Email:     email,
			Frequency: frequency,
			Countries: countries,
			Type:      infoType,
		})
	}

	return emails, nil
}

type Server struct {
	DB *sql.DB
}

func (s *Server) subscribeToEmail(w http.ResponseWriter, req *http.Request) {
	var ep EmailPayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.processCreateEmails(ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) unsubscribeToEmail(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		http.Error(w, "404 not found.", http.StatusBadRequest)
	}
	var ep EmailPayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.processDeletingEmails(ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) subscribeNow(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusBadRequest)
	}
	var ep EmailPayload
	err := json.NewDecoder(req.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := retrieveDataFromCovidApi(ep.countriesFromString(), ep.Type, time.Now().AddDate(0, 0, -5), time.Now())
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to retreive covid-api data", http.StatusBadRequest)
		return
	}

	err = sendEmail(ep.Email, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to send email", http.StatusBadRequest)
		return
	}
}

func (s *Server) processCreateEmails(payload EmailPayload) error {
	var (
		email     string
		frequency int64
		dueTime   time.Time
	)
	getEmailQuery := fmt.Sprintf("SELECT email, dueTime, frequency from covidhub.emails where email = '%s'", payload.Email)
	result := s.DB.QueryRow(getEmailQuery)
	err := result.Scan(&email, &dueTime, &frequency)
	if err != nil {
		dueTime := time.Now().Add(time.Second * time.Duration(payload.Frequency)).Format(Layout)
		createEmailQuery := fmt.Sprintf("INSERT INTO covidhub.emails (email, frequency, dueTime, countries, informationType) VALUES ('%s', '%d', '%s', '%s', '%s')",
			payload.Email,
			payload.Frequency,
			dueTime,
			payload.Countries,
			payload.Type,
		)
		_, err := s.DB.Query(createEmailQuery)
		if err != nil {
			log.Printf("error in creating email entry for: %s due to %v", payload.Email, err)
			return err
		}
	} else {
		newDueTime := dueTime.Add(time.Second * time.Duration(payload.Frequency-frequency)).Format(Layout)
		updateEmailQuery := fmt.Sprintf("UPDATE covidhub.emails SET email = '%s', frequency = '%d', dueTime = '%s', countries = '%s', informationType = '%s' WHERE email = '%s'",
			payload.Email,
			payload.Frequency,
			newDueTime,
			payload.Countries,
			payload.Type,
			payload.Email,
		)

		_, err = s.DB.Query(updateEmailQuery)
		if err != nil {
			log.Printf("error in updating email entry for: %s due to %v", payload.Email, err)
			return err
		}
	}
	return nil
}

func (s *Server) processDeletingEmails(payload EmailPayload) error {
	deleteEmailQuery := fmt.Sprintf("DELETE FROM covidhub.emails WHERE email = '%s'", payload.Email)
	_, err := s.DB.Query(deleteEmailQuery)
	if err != nil {
		log.Printf("error in deleting email entry for: %s due to %v", payload.Email, err)
		return err
	}

	return nil
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
	mux.HandleFunc("/subscribe", s.subscribeToEmail)
	mux.HandleFunc("/unsubscribe", s.unsubscribeToEmail)
	mux.HandleFunc("/subscribe/now", s.subscribeNow)

	log.Println("email server listening at port 8001")
	handler := c.Handler(mux)
	http.ListenAndServe(":8001", handler)
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

	if DefaultEmail = os.Getenv("DEFAULT_EMAIL"); DefaultEmail == "" {
		log.Fatalf("DEFAULT_EMAIL not set")
	}
	if DefaultEmailPassword = os.Getenv("DEFAULT_EMAIL_PASSWORD"); DefaultEmailPassword == "" {
		log.Fatalf("DEFAULT_EMAIL_PASSWORD not set")
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

	emailSender := &EmailSender{DB: sqlDB}
	go emailSender.Run()

	server := &Server{
		DB: sqlDB,
	}
	server.Run()
}
