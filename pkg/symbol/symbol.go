package symbol

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Result struct {
	Symbol    string
	Mean      float64
	Timestamp time.Time
}

type Message struct {
	Symbol    string  `json:"s"`
	Price     float64 `json:"p"`
	Timestamp int64   `json:"t"`
}

type Response struct {
	Type string    `json:"type"`
	Data []Message `json:"data"`
}

type SMACalculator struct {
	Name     string
	Window   int
	Count    float64
	Mean     float64
	Messages []Message
	File     *os.File
}

func (s *SMACalculator) Add(msg Message) {
	if len(s.Messages) < s.Window {
		s.Messages = append(s.Messages, msg)
	} else {
		s.Messages = append(s.Messages[1:], msg)
	}
}

func (s *SMACalculator) Calculate(out chan Result) {
	if len(s.Messages) < s.Window {
		return
	}
	for _, msg := range s.Messages {
		s.Count++
		s.Mean += (msg.Price - s.Mean) / s.Count
	}
	out <- Result{Timestamp: time.Now().UTC(), Symbol: s.Name, Mean: s.Mean}
}

func (s *SMACalculator) Write(result Result) error {
	output := fmt.Sprintf("%s\t%s\t%f\n", result.Timestamp, result.Symbol, result.Mean)
	if _, err := s.File.WriteString(output); err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}

type Service struct {
	w       *websocket.Conn
	symbols map[string]*SMACalculator
	results chan Result
	errors  chan error
}

func NewService(n int, apikey string, symbolNames []string) (*Service, error) {

	results := make(chan Result, 1)
	errors := make(chan error, 1)

	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", apikey), nil)
	if err != nil {
		return nil, err
	}

	symbols := map[string]*SMACalculator{}
	for _, name := range symbolNames {
		msg, err := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": name})
		if err != nil {
			return nil, err
		}
		if err := w.WriteMessage(websocket.TextMessage, msg); err != nil {
			return nil, err
		}
		file, err := os.OpenFile(fmt.Sprintf("data/%s", name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		symbols[name] = &SMACalculator{
			Name:   name,
			Window: n,
			File:   file,
		}
	}

	return &Service{
		w:       w,
		symbols: symbols,
		results: results,
		errors:  errors,
	}, nil
}

func (svc *Service) Run() error {
	go func() {
		for {
			var response Response
			if err := svc.w.ReadJSON(&response); err != nil {
				svc.errors <- err
			}
			if response.Type == "trade" {
				if len(response.Data) < 1 {
					svc.errors <- fmt.Errorf("no messages recieved")
				}
				msg := response.Data[len(response.Data)-1]
				svc.symbols[msg.Symbol].Add(msg)
				svc.symbols[msg.Symbol].Calculate(svc.results)
			}
		}
	}()
	for {
		select {
		case err := <-svc.errors:
			return err
		case result := <-svc.results:
			svc.symbols[result.Symbol].Write(result)
		}
	}
}

func (svc *Service) Close() error {
	if err := svc.w.Close(); err != nil {
		return err
	}
	for _, symbol := range svc.symbols {
		if err := symbol.File.Close(); err != nil {
			return err
		}
	}
	return nil
}
