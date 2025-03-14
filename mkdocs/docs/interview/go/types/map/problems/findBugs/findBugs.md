??? question "Найти и исправить 2 бага"
    
    === "Условие"
        ```go
        // find and fix 2 bugs
        func main() {
            writes := 1000
            var storage map[int]int
            
            wg := sync.WaitGroup {}
            
            wg.Add(writes)
            for i := 0; i < writes; i++ {
                i:=i
                go func() {
                    defer wg. Done ()
                    storage[i] = i I
                }()
            }
            
            wg. Wait ()
            fmt.Println(storage)
        }
        ```
    
    === "Решение"
        ```go
        //- Assignment to nil map: мы не можем записывать по ключам в мапу, которую созали таким образом, также как и в слайс
        //- Мапу не конкурентно безопасна, нужно обложить ее мьютексом

        package main

        import (
            "fmt"
            "sync"
        )

        func main() {
            writes := 1000
            //var storage map[int]int
            storage := make(map[int]int, writes)

            wg := sync.WaitGroup{}
            mu := sync.Mutex{}

            wg.Add(writes)
            for i := 0; i < writes; i++ {
                i := i
                go func() {
                    defer wg.Done()

                    mu.Lock()
                    defer mu.Unlock()
                    storage[i] = i
                }()
            }

            wg.Wait()
            fmt.Println(storage)
        }
        ```
    
    === "Объяснение"

    - Assignment to nil map: мы не можем записывать по ключам в мапу, которую созали таким образом, также как и в слайс
    - Мапу не конкурентно безопасна, нужно обложить ее мьютексом

??? question "Найти и исправить ошибку"
    
    === "Условие"
        ```go
        func main() {
            storage := make (map[int]int, 1000)
            
            wg := sync. WaitGroup{}
            ops := 1000
            mU := sync.RWMutex{}
            
            wg.Add (ops)
            for i := 0; i < ops; i++ {
                i: 1
                go func() {
                    defer wg. Done ()
                    
                    mu. Lock()
                    defer mu. UnLock()
                    storage[i] = i
                }()
            }
            
            // wg.Wait()
            wQ. Add (ops)
            for 1 := 0; 1 < ops; i++ {
                i: 1
                go func() {
                    defer wg. Done ()
                    
                    -, - = storage [i]
                }()
            }
            
            wg.Wait ()
            fmt.Printin(storage)
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
            storage := make(map[int]int, 1000)

            wg := sync.WaitGroup{}
            ops := 1000
            mu := sync.RWMutex{}

            wg.Add(ops)
            for i := 0; i < ops; i++ {
                i := 1
                go func() {
                    defer wg.Done()

                    mu.Lock()
                    defer mu.Unlock()
                    storage[i] = i
                }()
            }

            // wg.Wait()
            wg.Add(ops)
            for i := 0; i < ops; i++ {
                i := 1
                go func() {
                    defer wg.Done()

                    mu.RLock()
                    defer mu.RUnlock()
                    _, _ = storage[i]
                }()

                wg.Wait()
                fmt.Println(storage)
            }
        }
        ```
    
    === "Объяснение"

    Чтение из мапы - потокобезопасная операция, но, в данном случае, чтение будет пересекаться с записью и выйдет ошибка concurrent map read write, но мы не хотим блокироваться на мьютексе при параллельном чтении - поэтому используем RLock/RUnlock