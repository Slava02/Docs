??? question "Что будет выведено на экран?"
    
    === "Условие"
        ```go
        func main() {
            m := map[string]int{
                "foo": 1,
                "bar": 2,
                "baz": 3,
            }

            for b1, b2 := range m {
                println(b1, b2) // <- ?
            }

            a1, a2 := m["foo"]
            println(a1, a2) // <- ? 
        }
        ```
    
    === "Решение"
        ```go
        package main

        func main() {
            m := map[string]int{
                "foo": 1,
                "bar": 2,
                "baz": 3,
            }

            for b1, b2 := range m {
                println(b1, b2)
            }

            //foo 1
            //bar 2
            //baz 3

            a1, a2 := m["foo"]
            println(a1, a2)

            // 1 true
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что будет выведено на экран?"
    
    === "Условие"
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        var m = map[string]int{"a": 1}

        func main() {
            go read()
            time.Sleep(1 * time.Second)
            go write()
            time.Sleep(1 * time.Second)
        }
        func read() {
            for {
                fmt.Println(m["a"])
            }
        }
        func write() {
            for {
                m["a"] = 2
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

        var m = map[string]int{"a": 1}

        func main() {
            go read()
            time.Sleep(1 * time.Second)
            go write()
            time.Sleep(1 * time.Second)
        }
        func read() {
            for {
                fmt.Println(m["a"])
            }
        }
        func write() {
            for {
                m["a"] = 2
            }
        }

        // Куча 1
        // concurrent map read and map write
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что будет выведено на экран?"
    
    === "Условие"
        ```go
        package main

        import "log"

        type Struct {
            p *int
        }

        func main() {
            m := map[interface{}]int{
                1: 1,
                2: 2,
            }

            log.Println(m)
            m = map[interface{}]int{
                [3]int{1, 2, 3}: 1,
                []int{1, 2, 3}:  1,
            }
            
            log.Println(m)
        }
        ```
    
    === "Решение"
        ```go
        package main

        import "log"

        func main() {
            m := map[interface{}]int{
                1: 1,
                2: 2,
            }

            log.Println(m)
            m = map[interface{}]int{
                [3]int{1, 2, 3}: 1,
                //[]int{1, 2, 3}:  1,
            }

            log.Println(m)
        }

        // panic: runtime error: hash of unhashable type []int

        // After commenting
        // 2024/12/11 11:23:41 map[1:1 2:2]
        // 2024/12/11 11:23:41 map[[1 2 3]:1]
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Что произойдет?"
    
    === "Условие"
        ```go
        package main

        import "fmt"

        func main() {
            var m map[string]int
            fmt.Println(m["foo"])
            m["foo"] = 42
            fmt.Println(m["foo"])
        }
        ```
    
    === "Решение"
        ```go
        package main

        import "fmt"

        func main() {
            var m map[string]int
            fmt.Println(m["foo"])
            m["foo"] = 42
            fmt.Println(m["foo"])
        }

        // 0
        //panic: assignment to entry in nil map
        ```
    
    === "Объяснение"

    Необходимо заполнить