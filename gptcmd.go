/*
 gptcmd.go
 command-line access for Gpt Chat default format

 Requires 2 Env Variables:
	 GPTKEY="your OpenAI key" (required)
	 GPTMOD="engine model" (required)
	 GPTWRAP="line wrap length" (optional)
 Type your prompt on the command-line.

 log of requests is kept in file $HOME/gptcmd.log
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"strconv"
	"log"
	"bufio"
	"time"
	"runtime"
)

const (
	YEL  = "\033[33;1m"
	GRN  = "\033[0;32m"
	BLU  = "\033[34;1m"    // bright: blue
	ORG  = "\033[0;33m"   // kind of brown
	DFT  = "\033[0m\n"   // reset to default color
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Data struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

/* 	This is the struct that corresponds to the JSON
	return object. The JSON is "unmarshaled" into it.
*/
type Completion struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func wrap(input string, limit int) string {
	/* 	Create and return a new string by
		limiting substrings at word boundaries*/
    fields := strings.Fields(input)
    if limit <= 0 {
        return strings.Join(fields, " ")
    }
    format, count, result := "", 0, make([]string, len(fields))
    for index, word := range fields {
         if count + len(word) > limit {
             format = "\n"
             count = 0
         }
         count += len(word) + 1
         result[index] = format + word
         format = " "
    }
    return strings.Join(result, "")
}

func logRequest(prompt string, text string) {
	// Open file in append mode.
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = os.Getenv("USERPROFILE") + "\\"
	} else {
		filepath = os.Getenv("HOME") + "/"  // linux
	}

	if filepath == "" {
		return  // don't write to log file
	}

	file, err := os.OpenFile(filepath+"gptcmd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	// get the datetime string
	t := time.Now()
	layout := "Mon 01/02/2006 03:04 pm"
	tstr := t.Format(layout)

	// Write new data to file
	newData := "\n"+tstr+"\n> "+prompt+"\n>> "+text+"\n"
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(newData)

	// Save the changes
	err = writer.Flush()
	if err != nil {
		log.Fatalf("failed saving file: %s", err)
	}
}


/***
 *      __  __           _
 *     |  \/  |   __ _  (_)  _ __
 *     | |\/| |  / _` | | | | '_ \
 *     | |  | | | (_| | | | | | | |
 *     |_|  |_|  \__,_| |_| |_| |_|
 *
 */


func main() {

	helpstr := `
gptcmd v1.0 2023
Requires 2 Env Variables:
  GPTKEY="your OpenAI key" (required)
  GPTMOD="engine model" (required)
  GPTWRAP="line wrap length" (optional)
Type your prompt on the command-line.
`
	if len(os.Args) < 2 {
		fmt.Printf("%s%s%s",GRN, helpstr, DFT)
		os.Exit(0)
	}

	// Join all arguments with space as separator
	args := os.Args[1:]
	userprompt := strings.Join(args, " ")

	url := "https://api.openai.com/v1/chat/completions"
	// collect the Environment values
	openAPIKey := os.Getenv("GPTKEY")
	openAPIModel := os.Getenv("GPTMOD")
	wraplength := os.Getenv("GPTWRAP")

	// JSON for the Chat POST request data
	data := Data{
		Model: openAPIModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: userprompt,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// build the request object

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Handle the response ...

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// bodyBytes is in Byte format
	// fmt.Println(string(bodyBytes))

	var completion Completion
	err = json.Unmarshal(bodyBytes, &completion)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n%s%s says:\n", YEL, openAPIModel)
	fmt.Println("Model:", completion.Model)
	fmt.Println("Total Tokens:", completion.Usage.TotalTokens)
	respstr := completion.Choices[0].Message.Content
	if wraplength != "" {
	  	num, err := strconv.Atoi(wraplength)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Content:\n", wrap(respstr, num))
	} else {
		fmt.Println("Content:\n", respstr)
	}
	fmt.Println(DFT)
	logRequest(userprompt, respstr)
}
