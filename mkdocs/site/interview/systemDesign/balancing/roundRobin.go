package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

// В поле current будем хранить номер текущего сервера, на который распределяем нагрузку и обязательно обложимся мьютексом
type LoadBalancer struct {
	Current int
	Mutex   sync.Mutex
}

// Сами сервера будут представлены простой структуркой, с сылкой и параметром здорвовья
type Server struct {
	URL       *url.URL
	IsHealthy bool
	Mutex     sync.Mutex
}

// сердце нашего алгоритма - сам round robin
func (lb *LoadBalancer) getNextServer(servers []*Server) *Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	// Проходимся по всем серверам, проверяя каждый раз здоров ли он
	for i := 0; i < len(servers); i++ {
		// Получаем индекс текущего сервера
		idx := lb.Current % len(servers)
		nextServer := servers[idx]

		// Увеличиваем счетчик
		lb.Current++

		nextServer.Mutex.Lock()
		isHealthy := nextServer.IsHealthy
		nextServer.Mutex.Unlock()

		// Возвращаем сервер, если он здоров
		if isHealthy {
			return nextServer
		}
	}

	return nil
}

// Обратный прокси - это сервер, который находится между клиентов и одним или несколькими серверами. Он получает клиентский запрос и перенапрявлеят его бекнду, после чего - возвращает ответ сервера клиенту.
func (s *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

type Config struct {
	Port                string   `json:"port"`
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []string `json:"servers"`
}

func loadConfig(file string) (Config, error) {
	var config Config

	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// health check function that runs in given interval to check health of servers
func healthCheck(s *Server, healthCheckInterval time.Duration) {
	for range time.Tick(healthCheckInterval) {
		res, err := http.Head(s.URL.String())
		s.Mutex.Lock()
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Printf("%s is down\n", s.URL)
			s.IsHealthy = false
		} else {
			s.IsHealthy = true
		}
		s.Mutex.Unlock()
	}
}

// Чтобы посмотреть как работает балансировщик - нужно запустить несколько бекендов на портах из конфигов:python3 -m http.server 500(1-5)
// После - можно убедиться в том, что балансер корректно распредляет нагрузку, сделав клиентский запрос:curl -s -i http://localhost:8080 | grep -i "X-Forwarded-Server"
func main() {
	config, err := loadConfig("/Users/slava/GolandProjects/GolangTheory/SytemDesign/theory/balun/balancing/config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	var servers []*Server
	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		server := &Server{URL: u, IsHealthy: true}
		servers = append(servers, server)
		go healthCheck(server, healthCheckInterval)
	}

	lb := LoadBalancer{Current: 0}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := lb.getNextServer(servers)
		if server == nil {
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		// adding this header just for checking from which server the request is being handled.
		// this is not recommended from security perspective as we don't want to let the client know which server is handling the request.
		w.Header().Add("X-Forwarded-Server", server.URL.String())
		server.ReverseProxy().ServeHTTP(w, r)
	})

	log.Println("Starting load balancer on port", config.Port)
	err = http.ListenAndServe(config.Port, nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}
