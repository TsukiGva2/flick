## Example


```go
    package main

    import (
        "fmt"
        "log"
        //"time"
    )

    func InitFunctions(forth SerialForth) {

        forth.Run(": DRW 0 m $ d a ;")

        /*
            Expected output:

                PORTAL   xxxxxx
                REGIST.  xxxxxx
                UNICAS   xxxxxx
                COMUNICANDO WEB
        */

        forth.Run("VAR comm VAR all_tag VAR tag_unique VAR dev")

        /*
            : SCX ( l1 l2 l3 l4 -- )

                3 FOR
                    I DRW
                    50 DLY
                NXT

                0 DRW
            ;
        */

        forth.Run(fmt.Sprintf(": SCX 3 FOR I DRW NXT 0 DRW ;"))
    }

    func Screen1(forth SerialForth, device, tag_set, tag_cont, comunicando string) {

        forth.Run(
            fmt.Sprintf("%s %s %s %s SCX",
                forth.getBytes(device),
                forth.getBytes(tag_set),
                forth.getBytes(tag_cont),
                forth.getBytes(comunicando),
            ),
        )
    }

    func main() {

        forth, err := NewSerialForth()

        if err != nil {

            log.Fatalf("Error opening arduino: %v", err)
        }

        defer forth.Close()

        InitFunctions(forth)

        var tags_unicas, registros uint8 = 0, 0

        for range 300 {

            tags_unicas++
            registros++

            Screen1(forth,
                "PORTAL   701",
                fmt.Sprintf("UNICAS   %d", tags_unicas),
                fmt.Sprintf("REGIST.  %d", registros),
                "COMUNICANDO WEB",
            )
        }
    }
```
