package file_reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reader/modules/domain"
	"reader/modules/utils"
)

type Entry struct {
	Error   error
	FooData domain.FooData
}

type Stream struct {
	stream chan Entry
}

func NewJsonStream() Stream {
	return Stream{
		stream: make(chan Entry), // try to remove "stream: "
	}
}

func ReadFileWithStream(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	fmt.Println("func readFile with token: ", filename)
	d := json.NewDecoder(file)

	if _, err := d.Token(); err != nil {
		utils.Check(err)
		return err
	}

	for d.More() {
		// s, _ := d.Token()
		var fd domain.FooData
		if err := d.Decode(&fd); err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("read %q\n", fd)
	}

	return nil
}

func ReadFullFile(filename string) ([]domain.FooData, error) {
	var data []domain.FooData

	dataString, err := ioutil.ReadFile(filename)
	utils.Check(err)

	err = json.Unmarshal(dataString, &data)
	utils.Check(err)

	fmt.Println("Len of file dataset is: ", len(data))

	return data, nil
}

func ReadFileWithBufTocken(filename string) error {
	var data domain.Dataset

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()
	fmt.Println("func readFile: ", filename)

	scanner := bufio.NewScanner(file)
	dataString := scanner.Text()
	fmt.Println("DataString is: ", string(dataString))

	json.Unmarshal([]byte(dataString), &data)

	fmt.Println("Len of dataset is: ", len(data.Dataset))
	fmt.Println("Data dataset is: ", data)

	return nil
}

func ReadFileWithBuf(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	fmt.Println("func readFile: ", filename)

	reader := bufio.NewReader(file)

	buf := make([]byte, 16)

	i := 1
	for {
		n, err := reader.Read(buf)

		if err != nil {

			if err != io.EOF {

				log.Fatal(err)
			}

			break
		}

		fmt.Print(string(buf[0:n]))
		i++
	}

	var data domain.Dataset
	fmt.Println(data)
	return nil
}

func ReadFileWithTokens(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	fmt.Println("func readFile with token: ", filename)
	d := json.NewDecoder(file)

	if _, err := d.Token(); err != nil {
		utils.Check(err)
		return err
	}

	for d.More() {
		// s, _ := d.Token()
		var fd domain.FooData
		if err := d.Decode(&fd); err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("read %q\n", fd)
	}

	return nil
}

func ReadFileByLines(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	fmt.Println("func readFileWithTokens with token: ", filename)
	d := json.NewDecoder(file)

	for {
		var fd domain.FooData
		if err := d.Decode(&fd); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s: %s\n", fd.ID, fd.Data, fd.Signal)
	}
	return nil
}
