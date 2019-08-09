package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Scanf(title string) string {
	fmt.Println(title)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatalln("Scanf() Error:", err)
	}
	return New().ReplaceAll(scanner.Text(), " ", "+")
}
