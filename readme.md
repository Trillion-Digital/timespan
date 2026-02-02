Idea:

Timespan as an interface. The interface will be named Window and each implementor describes 
how to find its start and end dates. A non-custom window must not be created from both a
start and end date. For this especialized functions should be used in the format:

```go
w := timespan.WindowEndingOn(timespan.Month, t)
w := timespan.WindowStartingOn(timespan.Year, t)
```
