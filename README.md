# ISO8601 duration strings for Go

A package for working with ISO8601 duration strings. See <https://en.wikipedia.org/wiki/ISO_8601#Durations> for more info.

Some duration string examples:

* `P3Y6M4DT12H30M5.345S` - 3 years, 6 months, 4 days, 12 hours, 30 minutes, 5 seconds, and 345 milliseconds.
* `P3W` - 3 weeks
* `P1D3H10S` - 1 day, 3 hours, and 10 seconds

The package also supports a format extension to include signed durations and/or duration components. Some examples here:

* `-P1M10D` - negative 1 month and 10 days
* `P1Y-10D` - 355 days (1 year less 10 days)

## Usage

```go
import (
    "os"

    "github.com/bbetov/isoduration"
)

func main() {
    duration, err := isoduration.Parse("P4W")
    if err != nil {
        println(err)
        os.Exit(1)
    }
    println(duration.String())

    println(duration.StringWeeks())
}
```

## Limitations

The ISO 8601 allows the durations to be specified using a decimal notation.

Months are considered to be 30 days.

Years are considered to be 365 days.

## TODO

- [ ] Support `PYYYYMMDDThhmmss` format
- [ ] Support `P[YYYY]-[MM]-[DD]T[hh]:[mm]:[ss]` format
- [ ] Support decimal durations.
- [ ] Support different signs for the individual components (via a separate function)
