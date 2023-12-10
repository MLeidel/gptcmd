
# gptcmd
**A command line AI chat client for windows  
and linux written in Go using the OpenAI API**

 Requires 2 Environment Variables:  
>   GPTKEY="your OpenAI key" (required)  
    GPTMOD="engine model" (required)  
    GPTWRAP="line wrap length" (optional)  
    GPTTMP="temperature" (optional)  

Type your prompt on the command-line.
```bash
  $> gptcmd Hello my prompt is everything I type after gptcmd ...

```

A log of requests is kept in file _HOME_/gptcmd.log  
for Windows _USERPROFILE_/gptcmd.log

Compiled versions are offered here to use at your own risk.  
The source is here as well. Install Golang and compile yourself:  
```bash
  $> go build gptcmd.go   # linux
  $> GOOS=windows GOARCH=386 go build gptcmd.go   # windows
  $> GOOS=linux GOARCH=arm GOARM=5 go build gptcmd.go   # raspberry pi

```


---
**gptcmd**, **GptCLI**, and **gptcom** are basically  
go, c, and python versions of the following bash script:

```bash
#!/bin/bash

MODEL="gpt-4"

read -p "Enter prompt: " PROMPT

RESPONSE=$(curl -s https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${GPTKEY}" \
  -d '{
    "model": "'"$MODEL"'",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "'"$PROMPT"'"
      }
    ]
  }')

  date
  echo "You said: ${PROMPT}"
  echo "-------"
  echo $RESPONSE | jq -r '.choices[0].message.content'
```


