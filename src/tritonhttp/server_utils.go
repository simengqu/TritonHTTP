package tritonhttp

import (
	"bufio"
	"os"
	"log"
	"strings"
)
/** 
	Load and parse the mime.types file 
**/
func ParseMIME(MIMEPath string) (MIMEMap map[string]string, err error) {
	// panic("todo - ParseMIME")
	MIMEMap = make(map[string]string)
	f, err := os.Open(MIMEPath)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		paris := strings.Fields(s.Text())
		ext := paris[0]
		contentType := paris[1]
		MIMEMap[ext] = contentType
	}
	// log.Println(MIMEMap)
	return MIMEMap, err
}

// func main(){
// 	_, err := ParseMIME("/Users/simengqu/Documents/UCSD/CSE224/module-2-project-224-squ/src/mime.types")
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// }