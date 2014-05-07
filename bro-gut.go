package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func grab_value(line string) string {
	val := strings.Split(line, " ")[1]
	return val
}

func extract_sep(line string) []byte {
	sep := grab_value(line)
	sepchar, err := hex.DecodeString(sep[2:])
	if err != nil {
		log.Panic(err)
	}
	return sepchar
}

func bro_cut(convert_times bool, ofs string) {
	field_mapping := make(map[string]int)
	reader := bufio.NewReader(os.Stdin)
	for {
		line_, err := reader.ReadString('\n')
		line := strings.Trim(line_, "\n")

		if err != nil {
			// You may check here if err == io.EOF
			break
		}

		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#separator") {
				sep := extract_sep(line)
				fmt.Printf("%q\n", sep)
			} else if strings.HasPrefix(line, "#fields") {
				fields := strings.Split(line, "\t")[1:]
				fmt.Printf("%q\n", fields)
				for idx, field := range fields {
					field_mapping[field] = idx
				}
				fmt.Printf("%q\n", field_mapping)
				fmt.Printf("field query is index %d\n", field_mapping["query"])
			} else if strings.HasPrefix(line, "#types") {
				types := strings.Split(line, "\t")[1:]
				fmt.Printf("%q\n", types)
			}
		}
	}
}

func main() {
	var convert_times = flag.Bool("d", false, "Convert time values into human-readable format")
	var ofs = flag.String("ofs", "\t", "Sets a different output field separator.")
	flag.Parse()

	bro_cut(*convert_times, *ofs)
}
