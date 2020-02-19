package apiserver

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

var htmlTemplate = template.Must(template.ParseFiles("./front/mainPage.html"))

const pathToStandardFile = "./internals/app/apiserver/standard.txt"

// const pathToShadowFile = ""
// const pathToThinkyFile = ""

//New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

//Start function...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("starting api server...")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			text := r.FormValue("text")
			ascii := r.FormValue("ascii")

			if ascii == "default" {
				value := openFile(text, pathToStandardFile)
				fmt.Print(value)
				myvar := map[string]interface{}{"MyVar": value}
				htmlTemplate.Execute(w, myvar)
			}
		default:
			htmlTemplate.Execute(w, nil)
		}
	}
}

func openFile(enteredValue string, fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return "cant find file " + fileName
	}

	var asciiArray []string
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		asciiArray = append(asciiArray, scanner.Text())
		i++
	}

	if err1 := scanner.Err(); err1 != nil {
		fmt.Println(err1.Error())
		return "cant read file " + fileName
	}
	file.Close()
	return printASCII("", asciiArray, enteredValue)
}

func printASCII(value string, asciiArray []string, enteredValue string) string {
	var wordIndex int
	for i := 0; i <= 7; i++ {
		for j := range enteredValue {
			if enteredValue[j] == 92 && enteredValue[j+1] == 110 {
				wordIndex = j
				break
			} else {
				value += asciiArray[(int(enteredValue[j])-32)*9+1+i]
			}
		}
		value += "\n"
	}
	if wordIndex != 0 {
		printASCII(value, asciiArray, enteredValue[wordIndex+2:len(enteredValue)])
	}
	return value
}
