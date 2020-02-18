package apiserver

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

// var tpl = template.Must(template.ParseFiles("./front/html/mainPage.html"))
// var js = template.Must(template.ParseFiles("./front/js/script.js"))

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
	if err := s.configureStore(); err != nil {
		return err
	}
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
	// s.router.HandleFunc("/front/js/script.js", s.loadJSFile())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// tpl.Execute(w, nil)
	}
}

// func (s *APIServer) loadJSFile() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		js.Execute(w, nil)
// 	}
// }
