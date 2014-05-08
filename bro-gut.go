package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func contains_string(haystack []string, needle string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

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
func find_output_indexes(fields []string, cols []string, negate bool) []int {
	var out_indexes []int

	if len(cols) == 0 {
		out_indexes = make([]int, len(fields))
		for idx := range fields {
			out_indexes[idx] = idx
		}
		return out_indexes
	}

	field_mapping := make(map[string]int)
	for idx, field := range fields {
		field_mapping[field] = idx
	}
	if !negate {
		out_indexes = make([]int, len(cols))
		for idx, field := range cols {
			if field_index, ok := field_mapping[field]; ok {
				out_indexes[idx] = field_index
			} else {
				out_indexes[idx] = -1
			}
		}
	} else {
		out_indexes = make([]int, len(fields)-len(cols))
		cur_idx := 0
		for field_index, field := range fields {
			if !contains_string(cols, field) {
				out_indexes[cur_idx] = field_index
				cur_idx++
			}
		}

	}
	return out_indexes
}

func convert_time(ts string) string {
	seconds, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return ts
	}
	ms := int64(1000 * seconds)
	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Format(time.RFC3339Nano)
}

func bro_cut(cols []string, convert_times bool, ofs string, negate bool) {
	var out string
	var sep string
	var out_indexes []int
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
				out_indexes = find_output_indexes(fields, cols, negate)
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
		outparts := make([]string, len(out_indexes))
		for idx, val := range out_indexes {
			if val != -1 {
				if convert_times && time_fields[idx] {
					outparts[idx] = convert_time(parts[val])
				} else {
					outparts[idx] = parts[val]
				}
			}
		}
		out = strings.Join(outparts, ofs)
		fmt.Println(out)
	}
}

func main() {
	var convert_times = flag.Bool("d", false, "Convert time values into human-readable format")
	var ofs = flag.String("F", "\t", "Sets a different output field separator.")
	var negate = flag.Bool("n", false, "Print all fields *except* those specified.")
	flag.Parse()

	cols := flag.Args()

	bro_cut(cols, *convert_times, *ofs, *negate)
}
