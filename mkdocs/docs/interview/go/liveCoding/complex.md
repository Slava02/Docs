??? question "Необходимо сделать 2 запроса параллельно, получить код ответа и распечатать его на экран. Если будет ошибка - просто return. Сайты: avito.ru google.com"
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
            "net/http"
            "os"
            "sync"
        )

        func main() {
            sites := []string{"avito.ru", "google.com"}

            var wg sync.WaitGroup
            wg.Add(2)

            for _, v := range sites {
                go func() {
                    defer wg.Done()
                    resp, err := http.Get(v)
                    if err != nil {
                        os.Exit(2)
                    }
                    fmt.Println(resp)
                }()
            }

            wg.Wait()
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "## Теперь у нас добавилось 999999998 urls и мы хотим ограничить кол-во выполняемых горутин, как это сделать?"
    
    === "Условие"
        ```go
        package main

        import (
            "fmt"
            "net/http"
            "os"
            "sync"
        )

        func main() {
            sites := []string{"avito.ru", "google.com"..."site.ru"} // 1 000 0000 

            var wg sync.WaitGroup
            wg.Add(2)

            for _, v := range sites {
                go func() {
                    defer wg.Done()
                    resp, err := http.Get(v)
                    if err != nil {
                        os.Exit(2)
                    }
                    fmt.Println(resp)
                }()
            }

            wg.Wait()
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
            "golang.org/x/sync/errgroup"
            "net/http"
            "os"
        )
        // Вариант 1
        func main() {
            sites := []string{"avito.ru", "google.com"}

            ch := make(chan struct{}, 2)

            for _, v := range sites {
                ch <- struct{}{}
                go func() {
                    defer func() {
                        <-ch
                    }()
                    resp, err := http.Get(v)
                    if err != nil {
                        os.Exit(2)
                    }
                    fmt.Println(resp)
                }()
            }

            close(ch)
        }

        // Вариант 2
        func main() {
            sites := []string{"avito.ru", "google.com"}

            eg := errgroup.Group{}
            eg.SetLimit(2)

            for _, v := range sites {
                eg.Go(func() error {
                    resp, err := http.Get(v)
                    if err != nil {
                        os.Exit(2)
                    }
                    fmt.Println(resp)
                    return nil
                })
            }
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Теперь надо оганичить время выполнения программы до 10 секунд максимум"
    
    === "Условие"
        ```go
        // TODO
        ```
    
    === "Решение"
        ```go
        // TODO
        ```
    
    === "Объяснение"

    Необходимо заполнить


??? question "Сколько пройдет времени? Как оптимизировать до 3 секунд?"
    
    === "Условие"
        ```go
        package main

        import "time"

        func worker() chan int {
            ch := make(chan int)

            go func() {
                time.Sleep(3 * time.Second)
                ch <- 50
            }()

            return ch
        }

        func main() {
            timeStart := time.Now()
            _, _ = <-worker(), <-worker()
            println(int(time.Since(timeStart).Seconds()))
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "sync"
            "time"
        )

        func worker() chan int {
            ch := make(chan int)

            go func() {
                time.Sleep(3 * time.Second)
                ch <- 50
            }()

            return ch
        }

        func Task() {
            timeStart := time.Now()
            _, _ = <-worker(), <-worker()
            println(int(time.Since(timeStart).Seconds())) // 6 секунд
        }

        // Solution:
        func Solution() {
            timeStart := time.Now()
            <-merge(worker(), worker())
            println(int(time.Since(timeStart).Seconds()))
        }

        // Merger:
        func merge[T any](chans ...chan T) chan T {
            result := make(chan T, 1024)

            go func() {
                wg := &sync.WaitGroup{}
                wg.Add(len(chans))

                for _, ch := range chans {
                    ch := ch
                    go func() {
                        defer wg.Done()
                        for val := range ch {
                            result <- val
                        }
                    }()
                }

                wg.Wait()
                close(result)
            }()

            return result
        }

        func run(f func()) {
            f()
        }

        func main() {
            run(Solution)
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Необходимо вывести 10 раз FooBar"
    
    === "Условие"
        ```go
        package main

        import (
            "fmt"
        )

        func printFoo() {
            for i := 0; i < 10; i++ {
                fmt.Print("Foo")
            }
        }

        func printBar() {
            for i := 0; i < 10; i++ {
                fmt.Print("Bar\n")
            }
        }

        func main() {
            printFoo()
            printBar()
        }

        // Output (10 times):
        // FooBar
        // FooBar
        // ...
        // FooBar
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
            "sync"
        )

        // Solution 1
        func printFoo(ch1 chan struct{}, ch2 chan struct{}) {
            for i := 0; i < 10; i++ {
                <-ch1
                fmt.Print("Foo")
                ch2 <- struct{}{}
            }
            close(ch2)
        }

        func printBar(ch1 chan struct{}, ch2 chan struct{}) {
            for i := 0; i < 10; i++ {
                ch1 <- struct{}{}
                <-ch2
                fmt.Print("Bar\n")
            }
            close(ch1)
        }

        Solution 2
        func printFoo(ch chan struct{}) {
            for i := 0; i < 10; i++ {
                <-ch // block and Wait for signal
                fmt.Print("Foo")
            }
        }

        func printBar(ch chan struct{}) {
            for i := 0; i < 10; i++ {
                ch <- struct{}{}
                runtime.Gosched() // !!!
                fmt.Print("Bar\n")
            }
        }

        func main() {
            wg := sync.WaitGroup{}
            wg.Add(2)

            ch1, ch2 := make(chan struct{}), make(chan struct{})

            go func() {
                defer wg.Done()
                printFoo(ch1, ch2)
            }()

            go func() {
                defer wg.Done()
                printBar(ch1, ch2)
            }()

            wg.Wait()
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить


??? question "Необходимо написать свою реализацию кеша"
    
    === "Условие"
        ## Необходимо написать свою реализацию кеша

        🔫 Есть условие: дан интерфейс get и put, просят реализовать in-memory key/value хранилище.

        ```Go
        package Cache

        import "errors"

        var (
            ErrNotFound = errors.New("not found")
            ErrExpired  = errors.New("expired")
        )

        type Cache[K comparable, V any] interface {
            Get(K) (V, error)
            Set(K, V) error
        }
        ```


    
    === "Решение 1"
        ```go
        package Cache

        import (
            "sync"
        )

        type MyCache1[K comparable, V any] struct {
            m  map[K]V
            mu sync.RWMutex
        }

        func NewMyCache1[K comparable, V any]() *MyCache1[K, V] {
            return &MyCache1[K, V]{
                m:  make(map[K]V),
                mu: sync.RWMutex{},
            }
        }

        func (c *MyCache1[K, V]) Get(key K) (V, error) {
            c.mu.RLock()
            defer c.mu.RUnlock()

            v, ok := c.m[key]
            if !ok {
                return *(new(V)), ErrNotFound
            }

            return v, nil
        }

        func (c *MyCache1[K, V]) Set(key K, value V) error {
            c.mu.Lock()
            defer c.mu.Unlock()

            c.m[key] = value

            return nil
        }
        ```
    === "Решение 2"
        ```go
        package Cache

        import (
            "sync"
            "time"
        )

        type Item2[V any] struct {
            val V
            exp time.Time
        }

        type MyCache2[K comparable, V any] struct {
            m   map[K]Item2[V]
            mu  sync.RWMutex
            ttl time.Duration
        }

        func NewMyCache2[K comparable, V any](ttl time.Duration) *MyCache2[K, V] {
            c := &MyCache2[K, V]{
                m:   make(map[K]Item2[V]),
                mu:  sync.RWMutex{},
                ttl: ttl,
            }

            go func() {
                t := time.NewTicker(time.Second * 5)
                for _ = range t.C {
                    c.mu.Lock()
                    for key, v := range c.m {
                        if v.exp.Before(time.Now()) {
                            delete(c.m, key)
                        }
                    }
                    c.mu.Unlock()
                }
            }()

            return c
        }

        func (c *MyCache2[K, V]) Get(key K) (V, error) {
            c.mu.RLock()
            defer c.mu.RUnlock()

            v, ok := c.m[key]

            if !ok {
                return *(new(V)), ErrNotFound
            }

            if v.exp.Before(time.Now()) {
                return *(new(V)), ErrExpired
            }

            return v.val, nil
        }

        func (c *MyCache2[K, V]) Set(key K, val V) error {
            c.mu.Lock()
            defer c.mu.Unlock()

            c.m[key] = Item2[V]{
                val: val,
                exp: time.Now().Add(c.ttl),
            }

            return nil
        }

        ```
    
    === "Решение 3"
        ```go
        package Cache

        import (
            "sync"
            "time"
        )

        type Item3[V any] struct {
            val V
            exp time.Time
        }

        type Shard[K comparable, V any] struct {
            m  map[K]Item3[V]
            mu sync.RWMutex
        }

        type MyCache3[K comparable, V any] struct {
            shards []*Shard[K, V]
            ttl    time.Duration
        }

        func NewMyCache3[K comparable, V any](ttl time.Duration, shardsAmount int) *MyCache3[K, V] {
            cache := &MyCache3[K, V]{
                shards: make([]*Shard[K, V], shardsAmount),
                ttl:    ttl,
            }

            go func() {
                t := time.NewTicker(ttl)
                for _ = range t.C {
                    for _, shard := range cache.shards {
                        shard.mu.Lock()
                        for k, v := range shard.m {
                            if time.Now().After(v.exp) {
                                delete(shard.m, k)
                            }
                        }
                        shard.mu.Unlock()
                    }
                }
            }()

            return cache
        }

        func (c *MyCache3[K, V]) Get(key K) (V, error) {
            hashedKey := hash(key)
            shard := c.shards[hashedKey%len(c.shards)]

            shard.mu.RLock()
            defer shard.mu.RUnlock()

            item, ok := shard.m[key]

            if !ok {
                return *(new(V)), ErrNotFound
            }

            if time.Now().After(item.exp) {
                return *(new(V)), ErrExpired
            }

            return item.val, nil
        }

        func (c *MyCache3[K, V]) Set(key K, val V) error {
            hashedKey := hash(key)
            shard := c.shards[hashedKey%len(c.shards)]

            shard.mu.RLock()
            defer shard.mu.RUnlock()

            shard.m[key] = Item3[V]{
                val: val,
                exp: time.Now().Add(c.ttl),
            }

            return nil

        }

        func hash[K comparable](key K) int {
            return 0
        }

        ```

    === "Решение LRU"
        ```go
       package Cache

        type Node struct {
            Key, Val   int
            Next, Prev *Node
        }

        func newNode(key, val int) *Node {
            return &Node{
                Key: key,
                Val: val,
            }
        }

        type LRUCache struct {
            Head, Tail *Node
            Mp         map[int]*Node
            Capacity   int
        }

        func newLRUCache(head, tail *Node, cap int) LRUCache {
            return LRUCache{
                Head:     head,
                Tail:     tail,
                Mp:       make(map[int]*Node),
                Capacity: cap,
            }
        }

        func Constructor(capacity int) LRUCache {
            head, tail := newNode(0, 0), newNode(0, 0)

            head.Next = tail
            tail.Prev = head
            return newLRUCache(head, tail, capacity)
        }

        func (this *LRUCache) Get(key int) int {
            if node, ok := this.Mp[key]; ok {
                this.remove(node)
                this.insert(node)
                return node.Val
            }

            return -1
        }

        func (this *LRUCache) Put(key int, value int) {
            if node, ok := this.Mp[key]; ok {
                this.remove(node)
            }

            if len(this.Mp) >= this.Capacity {
                this.remove(this.Tail.Prev)
            }

            this.insert(&Node{Key: key, Val: value})
        }

        func (this *LRUCache) remove(node *Node) {
            delete(this.Mp, node.Key)
            node.Next.Prev = node.Prev
            node.Prev.Next = node.Next
        }

        func (this *LRUCache) insert(node *Node) {
            this.Mp[node.Key] = node
            next := this.Head.Next
            this.Head.Next = node
            next.Prev = node
            node.Prev = this.Head
            node.Next = next
        }

        ```
    
    === "Объяснение"

    🏓 Первое решение: написать структуру, в которой будет мапа и мьютекс. Таким образом просто и очень быстро запилить первое рабочее решение. Потом немного подумать и понять, что мьютекс можно замнить на RWMutex, чтобы ускорить чтение.

    🥊 Описанное выше решение может хорошо подойти джуну+/мидлу-, но мидлы-сеньеры здесь не остановятся. И на следующем вопросе от интервьюера про то, как ускорить чтение, начнут рассуждать.

    ❓На какие вопросы какие темы стоит поднимать:
    - Как ускорить чтение? И может даже запись? Можно порассуждать про шардирование, подумать, каким алгоритмам решардинга тут можно воспользоваться - углубляться как только можно
    - А что если хочет сохранять на диск, знаешь уже такие кейсы в других in-memory базах? Да, рассказываешь, чем плоха и хороша персистентность редиса
    - Какие проблемы могут возникнуть? В случае in-memory key/value известная проблема - это переполнение кэша, поэтому нужно думать о вытеснении. Тут есть смысл рассказать про LRU, LFU, TTL и другие алгоритмы. Один инстанс мапы упадет, что делать - еще одна проблема. Тут есть смысл рассказать про репликацию, копая в самую глубь. А как с шардирование соединяться?
    - далекий вопрос в сторону дизайна и опыта: Какие юзкейсы применения такого in-memory key/value видишь? На этот вопрос будет в пост в скором времени
    - совсем далекий: А что если бы ты разворачивал на двух и более дц, как избегать network partitionов, как сделать так, чтобы если один дц упадет - все было все равно доступно?