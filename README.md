## Go Figure

Simple configuration retrieval for Go:

```go
func main() {
	var dbconf struct {
		Development struct {
			Driver string `yaml:"driver"`
			Open   string `yaml:"open"`
		} `yaml:"development"`
	}

	gofigure.AddVariants("local")
	gofigure.Load("dbconf", &dbconf)

	fmt.Println(dbconf)
}
```
