// esplus is a helper cli for espanso
// run : esplus <command> <args>
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
)

var version string = "dev"

// cmdTemplate : esplus template <template> <args> : execute a template with args.
// If only one arg is given, execute the template with it as string value.
// If more than one arg is given, execute the template with them as []string value.
// The [sprig](github.com/Masterminds/sprig) functionsare available.
// First check if the arg[0] is a file, if so, read it and use it as template (with {{ and }} as delimiters).
// If not, consider it as template string. The delimiters are set to [[ and ]], so that they don't conflict with espanso's delimiters.
func cmdTemplate(args []string) (err error) {
	// args = os.Args[2:] = <file|template> <args>
	if len(args) == 0 {
		return nil
	}
	var tmpl *template.Template
	// check if args[0] is a file
	if _, err := os.Stat(args[0]); err == nil {
		// args[0] is a file, read it
		content, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		tmpl, err = template.New("tmpl").Funcs(sprig.FuncMap()).Parse(string(content))
	} else {
		// compile the template
		tmpl, err = template.New("tmpl").Delims("[[", "]]").Funcs(sprig.FuncMap()).Parse(args[0])
	}
	if err != nil {
		return err
	}
	if len(args) == 2 {
		return tmpl.Execute(os.Stdout, args[1])
	}
	return tmpl.Execute(os.Stdout, args[1:])
}

// cmdWait is called by 'esplus wait <milliseconds> <cmd> <args>'
// which is called by 'esplus run <milliseconds> <cmd> <args>'.
// It waits for milliseconds and then runs <cmd> <args>.
func cmdWait(args []string) (err error) {
	// args = os.Args[2:] = <milliseconds> <cmd> <args>
	if len(args) == 0 {
		return nil
	}
	wait, e := strconv.Atoi(args[0])
	if e != nil {
		return e
	}
	args = args[1:]
	if wait > 0 {
		// sleep for milliseconds
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
	if len(args) == 0 {
		return nil
	}
	// run the command (without release)
	c := exec.Command(args[0], args[1:]...)
	return c.Start()
}

// cmdRun is called by 'esplus run [milliseconds] <cmd> <args>'.
// If milliseconds is not given, execute <cmd> <args> without waiting for it to finish.
// If milliseconds is a number, execute 'esplus wait <milliseconds> <cmd> <args>' without waiting for it to finish.
func cmdRun(args []string) error {
	// args = os.Args = <path to esplus> run [milliseconds] <cmd> <args>
	if len(args) == 0 {
		return nil
	}
	var c *exec.Cmd
	_, err := strconv.Atoi(args[2])
	if err != nil { // no wait, execute : <cmd> <args>
		c = exec.Command(args[2], args[3:]...)
	} else { // execute esplus wait <milliseconds> <cmd> <args>
		args[1] = "wait"
		c = exec.Command(args[0], args[1:]...)
	}
	if err := c.Start(); err != nil {
		return err
	}
	c.Process.Release()
	return nil
}

var help string = `esplus is a helper cli for espanso. 
Version: ` + version + `
Usage: esplus <command> <args>

Commands:
	template <file> <args> : if file exists, use it as template with args (using {{ and }} as delimiters)
	template <template string> <args> : execute a template with args (using [[ and ]] as delimiters)
	run [milliseconds] <cmd> <args> : run a command (with delay) without waiting for it to finish

Examples:
	esplus template 'Hello [[.|upper]]' 'World'
	esplus template 'Hello [[range .]][[.|upper|printf "%s\n"]][[end]]' 'World' 'and' 'Earth'
	esplus template 'file.template.txt' 'World'
	esplus run 200 code .

Project repository:
	https://github.com/kpym/esplus
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	var err error
	switch os.Args[1] {
	case "template":
		err = cmdTemplate(os.Args[2:])
	case "wait":
		err = cmdWait(os.Args[2:])
	case "run":
		err = cmdRun(os.Args)
	default:
		err = fmt.Errorf("unknown command %s", os.Args[1])
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
