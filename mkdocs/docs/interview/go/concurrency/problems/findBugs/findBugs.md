??? question "Найти и исправить ошибки"
    
    === "Условие"
        ```go
        func main() {
            alreadyStored := make(map[int]struct{})
            mU := sync.Mutex{}
            capacity := 1600
            
            doubles := make([]int, 0, capacity)
            for 1 := 0; i < capacity; i+ {
                doubles = append(doubles, rand.Intn(10)) // create rand num 0...9
            }
            // 3, 4, 5, 0, 4, 9, 8, 6, 6, 5, 5, 4, 4, 4, 2, 1, 2, 3, 1
            
            uniqueIDs := make (chan int, capacity)
            wg:= sync. WaitGroup{)
            
            for i := 0; i < capacity; i++ {
                i := i
                wg.Add(1)
                go func() {
                    defer wg.Done ()
                    if _, ok := alreadyStored[ doubles[i]]; !ok {
                        mu.Lock()
                        alreadyStored[doubles[i]] = struct{}{}
                        mU. Unlock() // 0 unsafe.Si2/0f (alreadyStored[doubLes [1]])
                        
                        uniqueIDs <- doubles[i]
                    }
                }()
            
            }
            
            wg.Wait()
            for val := range uniqueIDs {
                fnt.Println (val)
            }
            
            fmt.Printf("len of ids: #{len (uniqueIDs)H\n") // 0, 1, 2, 3,
            fmt.Printin(uniqueIDs)
        }
        ```
    
    === "Решение"
        ```go
        //- Во-первых, словим дедлок, так как будем читать из канала, который никто не закроет
        //- Во-вторых, между чтением из мапы и записью в мапу - будут втискиваться горутины до закрыти мьютекса и будем ловить read write и получать дубли

        package main

        import (
            "fmt"
            "math/rand"
            "sync"
        )

        func main() {
            alreadyStored := make(map[int]struct{})
            mu := sync.Mutex{}
            capacity := 1600

            doubles := make([]int, 0, capacity)
            for i := 0; i < capacity; i++ {
                doubles = append(doubles, rand.Intn(10)) // create rand num 0...9
            }
            // 3, 4, 5, 0, 4, 9, 8, 6, 6, 5, 5, 4, 4, 4, 2, 1, 2, 3, 1

            uniqueIDs := make(chan int, capacity)
            wg := sync.WaitGroup{}

            for i := 0; i < capacity; i++ {
                i := i
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    mu.Lock()
                    defer mu.Unlock()
                    if _, ok := alreadyStored[doubles[i]]; !ok {
                        alreadyStored[doubles[i]] = struct{}{}
                        uniqueIDs <- doubles[i]
                    }
                }()

            }

            wg.Wait()
            close(uniqueIDs)
            for val := range uniqueIDs {
                fmt.Println(val)
            }

            fmt.Printf("len of ids: #{len (uniqueIDs)H\n") // 0, 1, 2, 3,
            fmt.Println(uniqueIDs)
        }
        ```
    
    === "Объяснение"

    - Во-первых, словим дедлок, так как будем читать из канала, который никто не закроет
    - Во-вторых, между чтением из мапы и записью в мапу - будут втискиваться горутины до закрыти мьютекса и будем ловить read write и получать дубли