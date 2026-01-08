# hoopla

**A Go-based CLI demonstrating advanced search techniques, RAG principles, and LLM-assisted scoring pipelines.**

This CLI experiment showcases how to combine _multiple search strategies_, normalize and weight results, and integrate retrieval-augmented logic with language models. Although it's still a work-in-progress, it already includes meaningful implementations that demonstrate advanced search and ranking behaviors.

## Concepts covered

### Pre-Processing

- Case sensitivity
- Punctuation
- Tokenization
- Stop words
- Stemming

### Keyword Search

- Inverted Index
- Term Frequency (TF)
- Inverse Document Frequency (IDF)
- TF-IDF (combines both metrics)
- BM25 (Best Matching version 25)
- BM25-IDF / BM25-TF
- BM25search

### Semantic Search

- Embeddings
- Dot Product Similarity
- Cosine Similarity
- Chunking
- Semantic Chunking (with overlap)

### Hybrid Search

- Score Normalization
- Weighted Search
- Reciprocal Rank Fusion (RRF Search)

### LLMs

- Pre-Process/enhance query (check spell, re-write, expansion)
- Re-Ranking (individual, batch, cross-encoder)

### Evaluation

- Manual evaluation
- Golden Dataset
- Precision, recall and f1 score metrics
- Error analysis, debug/tracing, structured logs
- LLM evaluation

## Design Decisions & Architecture Overview

- **Idiomatic Project Structure**: Follows the standard Go `cmd/` and `internal/` layout to clearly separate CLI concerns from core search and ranking logic, keeping the domain code modular and testable.

- **Interface-Driven Extensibility**: Retrieval and ranking strategies (BM25, semantic search, RRF) are defined behind interfaces, making it easy to introduce new algorithms or hybrid approaches without changing the search pipeline.

- **Composable Search Pipelines**: The search workflow is modeled as a sequence of independent stages—query analysis, retrieval, fusion, and re-ranking—allowing experimentation with different scoring and combination strategies.

- **Decoupled LLM Integration**: LLM interactions are abstracted behind a dedicated service layer, enabling model swaps for query enhancement and evaluation while keeping retrieval logic isolated from model-specific concerns.

- **Performance & Concurrency**: Uses Go’s native concurrency primitives to parallelize computationally intensive stages, including BM25 document scoring via a bounded worker pool (`internal/index/inverted_index.go`) and embedding generation during semantic indexing (`internal/methods/semantic_search.go`), improving throughput for large data sets.

<br>
<br>

<sub>
Disclaimer: This project was inspired by a Python-based bootcamp exercise.
The Go implementation, architecture, and all extensions were designed and developed independently, including additional features and improvements.
The goal is to demonstrate system design, search relevance techniques, and how Go can be an excellent choice for building performant CLIs and LLM-powered applications.
</sub>
