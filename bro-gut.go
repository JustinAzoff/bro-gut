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

func extract_sep(line string) string {
	sep := grab_value(line)
	sepchar, err := hex.DecodeString(sep[2:])
	if err != nil {
		log.Panic(err)
	}
	return string(sepchar)
}

func bro_cut(cols []string, convert_times bool, ofs string) {
	col_size := len(cols)
	var out string
	var sep string
	field_mapping := make(map[string]int)
	time_fields := make(map[int]bool)
	var fields []string
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
				sep = extract_sep(line)
			} else if strings.HasPrefix(line, "#fields") {
				fields = strings.Split(line, "\t")[1:]
				for idx, field := range fields {
					field_mapping[field] = idx
				}
			} else if strings.HasPrefix(line, "#types") {
				types := strings.Split(line, "\t")[1:]
				for idx, typ := range types {
					if typ == "time" {
						time_fields[idx] = true
					}
				}
			}
			continue
		}
		parts := strings.Split(line, sep)
		outparts := make([]string, col_size)
		for idx, field := range cols {
			if field_index, ok := field_mapping[field]; ok {
				outparts[idx] = parts[field_index]
			}
		}
		out = strings.Join(outparts, ofs)
		fmt.Println(out)
	}
}

func main() {
	var convert_times = flag.Bool("d", false, "Convert time values into human-readable format")
	var ofs = flag.String("F", "\t", "Sets a different output field separator.")
	flag.Parse()

	cols := flag.Args()

	bro_cut(cols, *convert_times, *ofs)
}
