package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	objectSlice := map[string]string{
		"https://upload.wikimedia.org/wikipedia/commons/f/ff/Pizigani_1367_Chart_10MB.jpg":                                           "image.jpg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/%C3%89douard_Mendy_2021.jpg/220px-%C3%89douard_Mendy_2021.jpg":    "image2.jpg",
		"https://upload.wikimedia.org/wikipedia/en/thumb/b/b4/Sun_Records_%28TV_series%29.jpg/250px-Sun_Records_%28TV_series%29.jpg": "image3.jpg",
	}
	const workers = 10
	wg := new(sync.WaitGroup)
	in := make(chan string, 2*workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for char := range in {
				sliceChar := strings.Split(char, ",")
				downloadObject(sliceChar[0], sliceChar[1])
			}
		}()
	}

	for keyString, sliceValue := range objectSlice {
		if keyString != "" {
			in <- keyString + "," + sliceValue
		}
	}
	close(in)
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Program took %s to run", elapsed)

}

func downloadObject(url string, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	object, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	ioutil.WriteFile(filename, object, 0666)

	log.Printf("Object %s saved", filename)
}
