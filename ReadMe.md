# rag-course
## Set up
### Ollama
Install as per instructions on: https://ollama.com/

Tested running models <b>locally</b> with:
- gemma3 - <code>ollama pull gemma3</code>
- llama3.1 - <code>ollama pull llama3.1</code>

For embedding:
- nomic-embed-text - <code>ollama pull nomic-embed-text</code>

For image interpretation:
- mistral-small3.1 - <code>ollama pull mistral-small3.1</code>


### Docker
Used for Postgres with vector support

Create containers:

<code>docker compose up</code>

To recreate containers:

<code>docker compose up --force-recreate</code>

Or

<code>docker compose down</code>

<code>docker compose up</code>

### Configuration
See <code>.env</code> 

This rag model almost exclusively uses the RAG and not the general knowledge
of the LLM. This is on purpose, just to show how it works.

Basic LLM/RAG manipulation can be done by editing:
<code>prompts/system-custom.md</code>
<code>rag/prompt.go</code>

Have fun exploring!