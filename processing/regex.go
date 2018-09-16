package processing

import (
	"fmt"
	"regexp"
)

func replaceCommands(src string, regEx string) string {

	//MustCompile simplifies safe initialization of global variables holding compiled regular expressions
	r, _ := regexp.Compile(regEx)

	//replace values with what we specify using the regex above

	tmp := r.ReplaceAllString(src, "*cmd*")

	fmt.Println(tmp)

	return tmp

}
