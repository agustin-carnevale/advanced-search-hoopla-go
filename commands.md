### 🔍 Keyword Search

Classical keyword-based retrieval using an inverted index and probabilistic ranking
models such as **TF-IDF** and **BM25**.

**Use this when:**

- You need exact term matching
- You want explainable relevance scores
- You’re benchmarking classical IR approaches

### Common commands

| Command              | Description                |
| -------------------- | -------------------------- |
| `build`              | Build the inverted index   |
| `search`             | Basic keyword search       |
| `tf`, `idf`, `tfidf` | Inspect scoring components |
| `bm25search`         | Full BM25 ranking          |
| `bm25searchP`        | Parallel BM25 search       |

### Examples

Basic and advanced keyword-based retrieval using inverted indices.

```bash
# Build the inverted index from data
./hoopla keyword build

# Basic keyword search
./hoopla keyword search "Christopher Nolan"

# Advanced BM25 scoring search
./hoopla keyword bm25search "dark knight" --limit 10
```

### 🧠 Semantic Search

Uses vector embeddings to find documents based on meaning rather than just exact word matches.

```bash
# Verify embeddings exist or generate them
./hoopla semantic verifyEmbeddings

# Semantic similarity search
./hoopla semantic search "movies about space travel"

# Search using semantic chunking for long documents
./hoopla semantic searchChunked "intense psychological thriller"
```

### 🔀 Hybrid Search

Combines the precision of keyword search with the conceptual depth of semantic search.

```bash
# Reciprocal Rank Fusion (RRF) search
./hoopla hybrid rrfSearch "Batman movies with Joker"

# Weighted search with custom importance
./hoopla hybrid weightedSearch "superhero action" --limit 5

# RRF search with LLM-powered query expansion and re-ranking
./hoopla hybrid rrfSearch "Batman" --enhance expand --rerankMethod crossEncoder
```

### 🤖 RAG (Retrieval-Augmented Generation)

Connects the search results to an LLM to provide natural language answers, summaries, and citations.

```bash
# Ask a question based on retrieved context
./hoopla rag question "Who is the main antagonist in The Dark Knight?"

# Generate a summary of the search results
./hoopla rag summarize "Summarize the plots of these Batman movies"

# Answer with specific citations to the source documents
./hoopla rag citations "What are the common themes in Nolan's movies?"
```

### 🖼️ Multimodal

Leverages vision-language models for image-related tasks.

```bash
# Describe an image and use it in a RAG pipeline
./hoopla multimodal describeImage 'bear in London' --imagePath "data/paddington.jpeg"
```

### 📊 Evaluation & Debugging

Tools for measuring performance and debugging the search pipeline.

```bash
# Evaluate retrieval performance against a golden dataset
./hoopla evaluation --limit 20

# Run a search with detailed debug logging
./hoopla hybrid rrfSearch "query" --debug
```
