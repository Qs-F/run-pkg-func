package main

import (
	"flag"
	"github.com/Qs-F/gondom"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	process := flag.String("p", "", "(option) processes must be given through option")
	pkg := flag.String("pkg", "", "(option) pkg name")
	fn := flag.String("f", "", "(must) Function name you want to run")
	args := flag.String("a", "", "(option) Arguments of function you want to run. Each argument must be separated with `,`")
	flag.Parse()
	s := ""
	var o []byte
	if *pkg == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}
		env := os.Getenv("GOPATH")
		s = strings.TrimPrefix(wd+"/", env+"/src/")
		s = strings.TrimSuffix(s, "/")
		o, err = exec.Command("go", "test").Output()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		s = *pkg
	}
	log.Println("Test for " + s)
	rand := gondom.Make(12, time.Now().UnixNano())
	pkgName := strings.Split(s, "/")[len(strings.Split(s, "/"))-1]
	content := `// tmp go file generated by github.com/Qs-F/run-pkg-func/generate.go
package main

import (
	"fmt"
	"` + s + `"
)

func main() {
	` + *process + `
	fmt.Printf("Value is %v\n", ` + pkgName + `.` + *fn + `(` + *args + `))
}`
	err := ioutil.WriteFile("/tmp/"+rand+".go", []byte(content), 0700)
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove("/tmp/" + rand + ".go")
	log.Println("Run code is: \n", content)
	log.Println("File is: ", "/tmp/"+rand+".go")
	o, err = exec.Command("go", "run", "/tmp/"+rand+".go").Output()
	if err != nil {
		log.Println("Fatal:", err)
		return
	}
	log.Println(string(o))
}
