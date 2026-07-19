# rag-course
## Ollama
Install as per instructions on: https://ollama.com/

Tested running models <b>locally</b> with:
- gemma3 - <code>ollama pull gemma3</code>
- llama3.1 - <code>ollama pull llama3.1</code>

For embedding:
- nomic-embed-text - <code>ollama pull nomic-embed-text</code>


## Docker
Create containers:

<code>docker compose up</code>

To recreate containers:

<code>docker compose up --force-recreate</code>

Or

<code>docker compose down</code>

<code>docker compose up</code>

## Configuration
See <code>.env</code> 