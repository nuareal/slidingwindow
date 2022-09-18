package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nuareal/slidingwindow/persistent"
)

type Response struct {
	NumRequests int `json:"numRequests"`
}

type Server struct {
	data     *persistent.RequestData
	Server   http.Server
	muData   sync.RWMutex
	muUpdate sync.Mutex
}

func NewServer(filePath string, port string) *Server {

	server := Server{data: &persistent.RequestData{Timestamp: time.Now()}}
	_, err := os.Open(filePath)
	if err == nil {
		err = server.GetData().ReadFromFile(filePath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				log.Println("Failed to removed cache file", err)
			}
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(server.HandleRequest))

	server.Server = http.Server{
		Addr:    port,
		Handler: mux,
	}
	return &server
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {

	data := s.GetData()
	go s.updateData()
	numRequests := data.GetNumRequests()

	response := &Response{NumRequests: numRequests + 1}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println("Failed to respond", err)
		}
	}
}

func (s *Server) updateData() {

	s.muUpdate.Lock()
	defer s.muUpdate.Unlock()
	newData := *s.GetData()

	currentTime := time.Now()
	currentSeconds := currentTime.Second()

	diffTime := int(currentTime.Sub(newData.Timestamp).Seconds())
	if diffTime < 60 {
		preIndex := newData.Timestamp.Second()
		for index := (preIndex + 1) % 60; index != (currentSeconds+1)%60; index = (index + 1) % 60 {

			newData.Acummulated[index] = newData.Acummulated[preIndex] - newData.Count[index]
			newData.Count[index] = 0
			preIndex = index
		}
	} else {
		for index := 0; index < 60; index++ {
			newData.Acummulated[index] = 0
			newData.Count[index] = 0
		}
	}
	newData.Timestamp = currentTime
	newData.Count[currentSeconds]++
	newData.Acummulated[currentSeconds]++

	s.SetData(&newData)
}

func (s *Server) GetData() *persistent.RequestData {
	s.muData.RLock()
	defer s.muData.RUnlock()
	return s.data
}

func (s *Server) SetData(data *persistent.RequestData) {
	s.muData.Lock()
	s.data = data
	s.muData.Unlock()
}
