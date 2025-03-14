??? question " Что произойдет?"
    
    === "Условие"
    ```Go
    package  main
    import  "fmt"

    func  main()  {
        ch  :=  make(chan  int)
        select  {
        case  val  :=  <-ch:
            fmt.Println(val)
        default:
            fmt.Println("no  one  will  write  to  chan")
        }
    }
    ```

    === "Решение"
        ```go
        package main

        import "fmt"

        func main() {
            ch := make(chan int)
            select {
            case val := <-ch:
                fmt.Println(val)
            default:
                fmt.Println("no  one  will  write  to  chan")
            }
        }

        // Output: no  one  will  write  to  chan
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что выведет на экран следующий код?"
    
    === "Условие"
        ```go
        func main() {
            ch := make(chan int, 1)
            for i := 0; i < 10; i++ {
                select {
                case ch <- i:
                case val := <-ch:
                println(val)
                }
            }
        }
        ```
    
    === "Решение"
        ```go
        package main

        func main() {
            ch := make(chan int, 1)
            for i := 0; i < 10; i++ {
                select {
                case ch <- i:
                case val := <-ch:
                    println(val)
                }
            }
        }

        // Output: 0,2,4,6,8
        ```
    
    === "Объяснение"

    The output of the program is `0, 2, 4, 6, 8` because of the way the `select` statement works with the buffered channel.

    Here is a step-by-step explanation:

    1. The channel `ch` is created with a buffer size of 1.
    2. The `for` loop runs 10 times, with `i` ranging from 0 to 9.
    3. In each iteration of the loop, the `select` statement has two cases:
        - `case ch <- i:` attempts to send the value of `i` to the channel.
        - `case val := <-ch:` attempts to receive a value from the channel and assign it to `val`.

    Since the channel has a buffer size of 1, it can hold only one value at a time. The `select` statement will choose one of the cases that is ready to proceed.

    - When `i` is even (0, 2, 4, 6, 8), the channel is empty, so the `ch <- i` case is selected, and the value of `i` is sent to the channel.
    - When `i` is odd (1, 3, 5, 7, 9), the channel already contains the previous even value, so the `val := <-ch` case is selected, and the value is received from the channel and printed.

    Thus, the program prints the values received from the channel, which are the even numbers `0, 2, 4, 6, 8`.

??? question "Необходимо вывести Hello"
    
    === "Условие"
        ```go
        package main
        import "fmt"


        func main() {
            ch := make(chan string)
            ch <- "Hello"

            go func() {
                fmt.Println(<-ch)
            }()
        }
        ```
    
    === "Решение"
        ```go
        import (
            "fmt"
            "time"
        )

        // В той, версии, которая есть - fatal error: all goroutines are asleep - deadlock!

        // Answer:
        func main() {
            ch := make(chan string)

            go func() {
                fmt.Println(<-ch)
            }()

            ch <- "Hello"

            time.Sleep(10 * time.Millisecond)
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что произойдет и как починить?"
    
    === "Условие"
        ```go
        package main

        import "fmt"

        func main() {
            c := make(chan string)
            
            i := 0
            for ; i < 5; i++ {
                go func() {
                    c <- fmt.Sprintf("%d", i)
                }()
            }
            
            for {
                println(<-c)
            }
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
            "sync"
        )

        // 5, 5, 5, 5, 5 fatal error: all goroutines are asleep - deadlock!
        func Task() {
            c := make(chan string)
            i := 0
            for ; i < 5; i++ {
                go func() {
                    c <- fmt.Sprintf("%d", i)
                }()
            }
            for {
                println(<-c)
            }
        }

        // Solution 1:
        func Solution1() {
            c := make(chan string)

            var wg sync.WaitGroup

            for i := 0; i < 5; i++ {
                wg.Add(1)

                i := i

                go func() {
                    defer wg.Done()
                    c <- fmt.Sprintf("%d", i)
                }()
            }

            go func() {
                wg.Wait()
                close(c)
            }()

            for v := range c {
                fmt.Print(v)
            }
        }

        // Solution 2:
        func Solution2() {
            c := make(chan string)
            done := make(chan struct{})

            for i := 0; i < 5; i++ {
                i := i

                go func() {
                    c <- fmt.Sprintf("%d", i)
                    done <- struct{}{}
                }()

            }

            go func() {
                for i := 0; i < 5; i++ {
                    <-done
                }
                close(c)
            }()

            for {
                select {
                case v, ok := <-c:
                    if !ok {
                        return
                    }
                    fmt.Print(v)
                }
            }

        }

        func run(f func()) {
            f()
        }

        func main() {
            run(Solution2)
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что будет и как изменить?"
    
    === "Условие"
        ```go
        func WaitToCloseManyChans(a, b chan bool) {
            var aclosed, bclosed bool
            for !aclosed || !bclosed {
                select {
                case <-a:
                    aclosed = true
                case <-b:
                    bclosed = true
                }
            }
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        func WaitToCloseManyChans(a, b chan bool) {
            var aclosed, bclosed bool
            for !aclosed || !bclosed {
                select {
                case <-a:
                    aclosed = true
                case <-b:
                    bclosed = true
                }
            }
        }

        func WaitToCloseManyChansSol(a, b chan bool) {
            for a != nil || b != nil {
                select {
                case <-a:
                    a = nil
                case <-b:
                    b = nil
                }
            }
        }

        func main() {
            a := make(chan bool)
            b := make(chan bool)

            go func() {
                time.Sleep(1 * time.Second)
                a <- true
                close(a)
            }()

            go func() {
                time.Sleep(2 * time.Second)
                b <- true
                close(b)
            }()

            WaitToCloseManyChansSol(a, b)
            fmt.Println("Both channels are closed")
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что будет выведено на экран?"
    
    === "Условие"
        ```go
        func main(){
            wg:=sync.WaitGroup{}
            
            for i:=0;i<5;i++{
                wg.Add(1)
                go func(){
                    defer wg.Done()
                    fmt.Println(i)
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
            "sync"
        )

        func main() {
            wg := sync.WaitGroup{}

            for i := 0; i < 5; i++ {
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    fmt.Println(i)
                }()
            }

            wg.Wait()
        }

        // 4 1 2 3 0
        ```
    
    === "Объяснение"

    Выводится в таком порядке из-за: Локальности и Work Stealing: https://www.youtube.com/watch?v=-K11rY57K7k&t=1468s

    Локальная очередь работает по принципу FIFO, но совместно с ней используется 1-элементный LIFO стек для переключения между последними двумя горутинами. Так же этот LIFO элемент не может быть украден другим процессом. Это некая оптимизация в шедулере go, что есть некая вероятность, что у горутины которая создала нашу горутину, будет общий контекст, замыкания переменных и тд

??? question "Добавить timeout, чтобы избежать длительного ожидания"
    
    === "Условие"
        ```go
        // add timeout to avoid long waiting
        func main() {
            rand.Seed(time.Now().Unix())

            chanForResp := make(chan int)
            go RPCCall(chanForResp)

            result := <-chanForResp
            fmt.Println(result)
        }

        func RPCCall(ch chan<- int) {
            // sleep 0-4 sec
            time.Sleep(time.Second * time.Duration(rand.Intn(5)))

            ch <- rand.Int()
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "context"
            "errors"
            "fmt"
            "math/rand"
            "time"
        )

        type resp struct {
            id  int
            err error
        }

        // add ctx with timeout
        func main() {
            rand.Seed(time.Now().Unix())
            ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
            defer cancel()

            chanForResp := make(chan resp)
            go RPCCall(ctx, chanForResp)

            resp := <-chanForResp
            fmt.Println(resp.id, resp.err)
        }

        func RPCCall(ctx context.Context, ch chan<- resp) {
            select {
            case <-ctx.Done():
                ch <- resp{
                    id:  0,
                    err: errors.New("request aborted due timeout"),
                }
            case <-time.After(time.Second * time.Duration(rand.Intn(5))):
                ch <- resp{
                    id: rand.Int(),
                }
            }
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить