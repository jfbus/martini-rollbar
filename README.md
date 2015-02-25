# martini-rollbar
A martini middleware for rollbar

The middleware :

* forwards all panics to rollbar.com,
* only triggers in the production environment.

```
import "github.com/jfbus/martini-rollbar"

func main() {
	m := martini.Classic()
	m.Use(rollbar.Recovery(rollbar.Config{Token: "YOUR SERVER TOKEN"}))
}
```

rollbar.Recover recovers panics, the default Recovery handler does nothing.
