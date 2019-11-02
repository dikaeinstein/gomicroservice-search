package features

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/dikaeinstein/gomicroservice-search/data"
)

var criteria string
var response *http.Response
var err error

func iHaveNoSearchCriteria() error {
	if criteria != "" {
		return fmt.Errorf("Criteria should be empty")
	}

	return nil
}

func iCallTheSearchEndpoint() error {
	u := "http://localhost:8082/search?criteria="

	if criteria != "" {
		u = u + url.QueryEscape(criteria)
	}

	response, err = http.Get(u)
	return err
}

func iShouldReceiveABadRequestMessage() error {
	if response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("Should have received a bad response")
	}

	return nil
}

func iHaveAValidSearchCriteria() error {
	criteria = "Fat Freddy's Cat"

	return nil
}

func iShouldReceiveAListOfKittens() error {
	var body []byte
	body, err := ioutil.ReadAll(response.Body)

	if len(body) < 1 || err != nil {
		return fmt.Errorf("Response body is empty")
	}

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	if response == nil || response["kittens"] == nil {
		return fmt.Errorf("No kittens were returned in the response")
	}

	return nil
}

var server *exec.Cmd
var store *data.MySQLStore
var outb, errb bytes.Buffer

func startServer() {
	outb = bytes.Buffer{}
	errb = bytes.Buffer{}

	server = exec.Command("go", "build", "../cmd/search.go")
	server.Run()

	server = exec.Command("./search")

	if os.Getenv("DEBUG") == "true" {
		server.Stderr = os.Stderr
		server.Stdout = os.Stdout
	}

	go server.Run()

	time.Sleep(3 * time.Second)
	fmt.Printf("Server running with pid: %v\n", server.Process.Pid)
}

func waitForDB() error {
	var err error

	for i := 0; i < 30; i++ {
		store, err = data.NewMySQLStore(os.Getenv("MYSQL_CONNECTION"))
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return err
}

func clearDB() {
	store.DeleteAllKittens()
}

func setupData() {
	store.CreateSchema()

	err := store.InsertKittens(
		[]data.Kitten{
			data.Kitten{
				ID:     "1",
				Name:   "Felix",
				Weight: 12.3,
			},
			data.Kitten{
				ID:     "2",
				Name:   "Fat Freddy's Cat",
				Weight: 20.0,
			},
			data.Kitten{
				ID:     "3",
				Name:   "Garfield",
				Weight: 35.0,
			},
		})

	if err != nil {
		log.Fatalln("Unable to insert data:", err)
	}
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have no search criteria$`, iHaveNoSearchCriteria)
	s.Step(`^I call the search endpoint$`, iCallTheSearchEndpoint)
	s.Step(`^I should receive a bad request message$`, iShouldReceiveABadRequestMessage)
	s.Step(`^I have a valid search criteria$`, iHaveAValidSearchCriteria)
	s.Step(`^I should receive a list of kittens$`, iShouldReceiveAListOfKittens)

	s.BeforeScenario(func(interface{}) {
		startServer()
		clearDB()
		setupData()
	})

	s.AfterScenario(func(interface{}, error) {
		server.Process.Signal(syscall.SIGINT)
		fmt.Print(outb.String())
		fmt.Print(errb.String())
	})

	err := waitForDB()
	if err != nil {
		log.Fatalln("Unable to connect to DB:", err)
	}
}
