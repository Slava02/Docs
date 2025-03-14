## Основные свойства систем

### Производительность
**latency** - времени время обслуживания одного клиента
**throughput** - Число обработанных клиентов в единицу 

Обычно нам дан latency и нам нужно выжать максимальный throughput из него

### Масштабируемость

Способность системы расти без потери производительности и других характеристик, а также без необходимости менять программную реализацию.

Виды масштабирования:

Вертикальное
Вертикальным называется масштабирование при котором для увеличения производительности в уже имеющееся оборудование добавляют новые процессоры, диски, память. Такой подход применяется в случаях, когда лимит производительности элементов инфраструктуры исчерпан.

Горизонтальное
Суть горизонтального масштабирования — в добавлении новых узлов в IT-инфраструктуру. Вместо того, чтобы увеличивать мощность отдельных компонентов узла, компания добавляет новые серверы. С каждым дополнительным узлом нагрузка перераспределяться между всеми узлами.

Параметры нагрузки:
- число запросов в секунду (RPS)
- трафик (KB/MB/GB/TB per second)
- число одновременных соединений (С10к, с100k)

## Балансировка нагрузки

Чаще всего - задача стоит примерно так: горизонтально масштабировались, есть 5 серверов, нужно балансировать нагрузку между ними

### Клиентская балансировка

В таком случае клиентское приложение само будет ответственно за распределение запросов между серверами. Это будет работать быстро, но ему прийдется знать о всех серверах и, более того, учитывать ситуации, когда кол-во инстансов меняется.

![img.png](img.png)

### Серверная балансировка

Вся логика балансировки переносится на балансер. Более популярный вариант. Но как балансировать нагрузку?

### Round Robin

По кругу распределяем задачи. Пример реализации можно посмотреть в файле roundRobin.go
Также, бывает "взвешенный round robin". Он позволяет задать некие веса серверам, чтобы балансировка происходила с приоритетами

### Sticky sessions

Представим, что какой-то пользователь пришел на какой-то инстанс и там закешировались о нем какие-то данные. Не выгодно было бы отправлять его на другой при следующем запросе. 
Хотелось бы чтобы пользователь закрепялся за одним инстансом. Мы могли бы просто сделать на балансере таблицу и мапить пользователей, но это лишние ресурсы. Лучше - получать номер сервера по хешу.

ЕСТЬ ПРОБЛЕМА: что делать когда сервера умирают/добавлются? - ответ будет в разделе про консистентное хеширование

- https://habr.com/ru/companies/domclick/articles/548610/ - настройка липких сессий на nginx
- https://habr.com/ru/companies/domclick/articles/551332/#1 -настройка липких сессий на кубере

### Least connections 

Подходит для websocket-серверов, когда на балансере мониторится кол-во открытых соединений и выбирается тот, на котором их меньше всего. Может работать лучше, чем round-robin, когда соед-ия разной длительности


## L4 и L7 балансировка

- L4 балансировка работает на транспортном уровне модели OSI. Она принимает решения о маршрутизации трафика на основе информации, содержащейся в заголовках TCP/UDP, таких как IP-адреса и порты.
- L7 балансировка работает на прикладном уровне модели OSI. Она анализирует содержимое HTTP-запросов и может принимать решения о маршрутизации на основе таких факторов, как URL, заголовки и куки.

## DNS балансировка

Что произойдет, если упадет балансер? - всему капут

Тут на помощь приходит dns балансировка, когда dns резолвит разные ip балансировщиков, тем самым - обеспечивая отказоустойчивость.

![img_1.png](img_1.png)

Часто ее используют для **geoDNS балансировки** - когда роутим в дата-центр по географическому принципу

![img_2.png](img_2.png)

## Service discovery

Функционал получения информации о новых инстансах часто выносится в отдельный сервис, но может быть реализован и на балансировщике. 

Более подробно - https://cloud.vk.com/blog/service-discovery/

Example:
```Go
/ В поле current будем хранить номер текущего сервера, на который распределяем нагрузку и обязательно обложимся мьютексом
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
```