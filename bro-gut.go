package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func bro_cut(convert_times bool, ofs string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// You may check here if err == io.EOF
			break
		}

		log.Println(line)
	}
}

func main() {
	var convert_times = flag.Bool("d", false, "Convert time values into human-readable format")
	var ofs = flag.String("ofs", "\t", "Sets a different output field separator.")
	flag.Parse()

	bro_cut(*convert_times, *ofs)
}
