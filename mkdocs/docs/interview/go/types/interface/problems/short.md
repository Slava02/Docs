??? question "Сработает ли приведение и вызов ф-ии?"
    
    === "Условие"
        ```go
        package main

        import (
            "fmt"
        )

        type errorString struct {
            s string
        }

        func (e errorString) Error() string {
            return e.s
        }

        func checkErr(err error) {
            fmt.Println(err == nil)
        }

        func main() {
            var e1 error
            checkErr(e1) // Что выведет?

            var e *errorString
            checkErr(e) // Что выведет?

            e2 := &errorString{}
            checkErr(e2) // Что выведет?

            fmt.Println(e == e2) // Что выведет?

            e = nil
            checkErr(e) // Что выведет?
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
        )

        type errorString struct {
            s string
        }

        func (e errorString) Error() string {
            return e.s
        }

        func checkErr(err error) {
            fmt.Println(err == nil)
        }

        func main() {
            var e1 error
            checkErr(e1) // true

            var e *errorString
            checkErr(e) // false

            e2 := &errorString{}
            checkErr(e2) // false

            fmt.Println(e == e2) // false (во втором инициализируется значение)

            e = nil
            checkErr(e) // true
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Сработает ли приведение и вызов ф-ии?"
    
    === "Условие"
        ```go
        package main

        type Bar struct{}

        func (b *Bar) X() {}
        func (b *Bar) Y() {}
        func (b *Bar) Z() {}

        type XY interface {
            X()
            Y()
        }

        type YZ interface {
            Y()
            Z()
        }

        func main() {
            var b XY = &Bar{}
            z := b.(YZ) // Сработает ли это приведение?
            z.X()       // А этот вызов?
            _ = z
        }
        ```
    
    === "Решение"
        ```go
        package main

        type Bar struct{}

        func (b *Bar) X() {}
        func (b *Bar) Y() {}
        func (b *Bar) Z() {}

        type XY interface {
            X()
            Y()
        }

        type YZ interface {
            Y()
            Z()
        }

        func main() {
            var b XY = &Bar{}
            z := b.(YZ) // Сработает
            z.X()       // Не сработает, так как YZ не реализует метод X()
            _ = z
        }
        ```
    
    === "Объяснение"

    Необходимо заполнить

??? question "Какой будет вывод?"
    
    === "Условие"
        ```go
        import (
            "fmt"
        )

        type Example struct {
            Value string
        }

        type MyInterface interface{}

        func example1() MyInterface {
            var e *Example 
            return e      
        }

        func example2() MyInterface {
            return nil 
        } 

        func main() {
            fmt.Println(example1()) 
            fmt.Println(example2()) 
            fmt.Println(example1() == example2()) 
        }
        ```
    
    === "Решение"
        ```go
        package main

        import (
            "fmt"
        )

        type Example struct {
            Value string
        }

        type MyInterface interface{}

        func example1() MyInterface {
            var e *Example // nil
            return e       // nil
        }

        func example2() MyInterface {
            return nil // nil
        }

        func main() {
            fmt.Println(example1()) // nil
            fmt.Println(example2()) // nil
            fmt.Println(example1() == example2()) // false
        } 
        ```
    
    === "Объяснение"

    Необходимо заполнить