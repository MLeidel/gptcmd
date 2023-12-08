
# gptcmd
**A command line AI chat client for windows  
and linux written in Go using the OpenAI API**


command-line access for Gpt Chat limited API format

 Requires 2 Env Variables:  
>   GPTKEY="your OpenAI key" (required)  
    GPTMOD="engine model" (required)  
    GPTWRAP="line wrap length" (optional)  

Type your prompt on the command-line.

log of requests is kept in file $HOME/gptcmd.log

**gptcmd** and **gptcom** *fancier/compiled ways* to do the following:

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


