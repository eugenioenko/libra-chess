[![Go Report Card](https://goreportcard.com/badge/github.com/eugenioenko/libra-chess)](https://goreportcard.com/report/github.com/eugenioenko/libra-chess)

# ⚖️ Libra Chess Engine

> A UCI-Compliant Chess Engine in Go

---

## 1. 📝 Overview

Libra Chess is a UCI (Universal Chess Interface) compliant chess engine written in Go. The primary goal of Libra is to achieve a balance between high performance, modern software architecture, and clarity of design. This project serves as an exploration of chess engine development leveraging Go's unique strengths in concurrency, tooling, and efficient compilation.

This engine is designed for chess enthusiasts, developers looking to understand chess engine internals, and as a demonstration of software engineering principles applied to a complex domain.

---

## 2. ✨ Key Features

- **UCI Protocol Compliant:** Seamless integration with popular UCI-compatible GUIs (e.g., CuteChess, CoreChess, PyChess).
- **Alpha-Beta Search:** Optimized Alpha-Beta pruning forms the core of the search algorithm.
- **Iterative Deepening:** Allows for flexible time management and progressive deepening of the search.
- **Transposition Tables:** Utilizes Zobrist hashing to store and retrieve previously evaluated positions, significantly speeding up search by avoiding redundant computations.
- **Piece-Square Tables (PSTs):** Employs PSTs for nuanced positional evaluation, guiding the engine's understanding of piece placement.
- **Material Evaluation:** Core evaluation component based on standard piece values.
- **Endgame Evaluation Heuristics:** Includes logic to encourage king activity and mating sequences in the endgame (e.g., incentivizing the stronger side's king to approach the opponent's king).
- **Move Generation:** Optimized and validated pseudo-legal move generation with legality checks.
- **Comprehensive Testing Suite:**
  - Unit tests for core logic (`go test`).
  - Perft testing for move generation correctness.
  - Automated match play against itself and Stockfish using `cutechess-cli` for strength benchmarking and regression testing.

---

## 3. 🛠️ Installation & Building

### 3.1. Prerequisites

- **Go:** Version 1.23 or higher.
- **Make:** Standard `make` utility.
- **`golangci-lint` (Optional):** For running linters (`make lint`). Install via their official instructions.
- **`cutechess-cli` (Optional but Recommended):** For running match tests. Download from its official repository. Ensure it's in your PATH or adjust `Makefile` paths.

### 3.2. Building the Engine

```bash
go build
# or
make build
```

This will produce the main executable `libra-chess`

### 3.3. Running the Engine

- **Directly (for UCI interaction):**
  ```bash
  ./libra
  ```
  The engine will start and await UCI commands.
- **With a UCI GUI:**
  Configure your favorite UCI GUI (e.g., CuteChess, CoreChess, PyChess) to use the `./libra-chess` executable.
- **Using `main.go` (if it contains a simple CLI or test loop):**
  ```bash
  go run main.go
  ```

### 3.4. Running Tests

- **Unit Tests:**
  ```bash
  make test
  ```
- **Linter:**
  ```bash
  make lint
  ```
- **Match Play vs. Itself (MainLibra vs PullLibra):**
  To run a versus match between two versions of the engine using `cutechess-cli`, you need to have two binaries: `libra-chess` and `libra-main`. If you want to test the current version against itself, simply copy the `libra-chess` binary to `libra-main`:
  ```bash
  cp libra-chess libra-main
  ```
  Then run:
  ```bash
  make test-cutechess
  ```
  _(Note: The Makefile refers to `./libra-chess` and `./libra-main`. Ensure these binaries are correctly built and named. `libra-chess` might be the current development version and `libra-main` a stable baseline or a copy for self-play.)_
- **Match Play vs. Stockfish:**
  ```bash
  make test-stockfish
  ```
  _(Requires Stockfish CLI to be at `./stockfish/stockfish-cli`)_
- **Debug Match:**
  ```bash
  make test-debug
  ```

---

## 4. 🏛️ Architectural Overview & Design Philosophy

Libra Chess is architected with modularity, simplicity, and maintainability as primary considerations. The choice of Go as the implementation language was deliberate, aiming to harness its excellent concurrency primitives, straightforward syntax, and robust standard library for an extra boost in performance.

### 4.1. Language Choice: Go

- **Advantages:**
  - **Concurrency:** Go's goroutines and channels offer a powerful yet simple model for concurrent programming, which is pivotal for future enhancements like parallel search.
  - **Performance:** While not C/C++, Go offers impressive performance, especially with its efficient garbage collector (GC) and direct compilation to machine code. Careful memory management is still crucial.
  - **Simplicity & Readability:** Go's clean syntax and established conventions promote maintainable and understandable code.
  - **Tooling:** Rich ecosystem including `gofmt` for automated formatting, `go test` for testing, and `golangci-lint` for static analysis.
- **Trade-offs & Considerations:**
  - **Garbage Collection:** While Go's GC is highly optimized, in performance-critical sections like deep search, GC pauses can be a concern. This necessitates careful memory allocation patterns (Move/UndoMove instead of cloning).
  - **Ecosystem for Chess Engines:** The C++ ecosystem for chess engines is more mature with a larger pool of shared libraries and knowledge. Libra aims to contribute to the growing Go presence in this domain.

### 4.2. Modularity

The engine's core logic is organized within the `pkg/` directory, with separation of concerns:

- `board.go`: Board representation, piece management, and core game state.
- `evaluate.go`: Static evaluation function, including material, PSTs, and endgame heuristics.
- `generate.go`: Move generation logic.
- `search.go`: Search algorithms (Alpha-Beta, iterative deepening).
- `tt.go`: Transposition table implementation.
- `zobrist.go`: Zobrist hashing for position keys.
- `move.go`: Move/UndoMove for calculations,
- `piece.go`: Pieces definitions,
- `utils.go`: Utility and data structure definitions.

This modular design facilitates easier testing, debugging, and future feature development.

### 4.3. Evaluation Function Design (`evaluate.go`)

The current evaluation function is a classical handcrafted one, balancing speed and accuracy.

- **Components:**
  - **Material:** Standard piece values (Pawn:100, Knight:300, Bishop:300, Rook:500, Queen:900).
  - **Piece-Square Tables (PSTs):** Static tables that assign positional bonuses or penalties to pieces based on their square. These are currently simplified and offer a good baseline.
  - **Endgame Heuristics:** A specific heuristic encourages the king with a material advantage to move towards the opponent's king when total material on the board is low (<= 14 units, excluding kings). This promotes checkmates in won endgames.
- **Trade-offs:**
  - **Speed vs. Accuracy:** Handcrafted evaluation functions are generally fast. The current PSTs are relatively simple, which makes evaluation quick but potentially less nuanced than more complex schemes or ML-based models.
  - **Complexity of Terms:** Adding more evaluation terms (e.g., detailed pawn structure, king safety, mobility beyond basic checks, passed pawns) can improve strength but increases computational cost and tuning complexity. This is a key area for future refinement.

### 4.4. Search Algorithm (`search.go`)

- **Minimax with Alpha-Beta Pruning:** A standard and effective framework for chess search.
- **Iterative Deepening:** Allows the engine to search to a certain depth, then use the information from that search (e.g., principal variation) to order moves for the next, deeper iteration. This is essential for effective time management.
- **Quiescence Search:** (Assumed, or a high-priority addition) To mitigate the horizon effect, a quiescence search is typically implemented to evaluate "quiet" positions by extending the search for captures and other tactical moves.
- **Trade-offs:**
  - **Search Depth vs. Time:** The deeper the search, generally the stronger the play, but time is a finite resource. Effective move ordering, pruning, and extensions/reductions are key to searching deeper within the allocated time.
  - **Selectivity:** Deciding which branches of the search tree to explore deeply (extensions) and which to prune or search shallowly (reductions) is a complex balancing act.

### 4.5. Testing and Benchmarking Philosophy

A robust testing strategy is paramount for engine development.

- **Correctness:** `perft` tests validate move generation exhaustively. Unit tests cover individual functions and modules.
- **Strength:** Regular match play against baseline versions of Libra, other engines like Stockfish (at controlled ELO levels), and itself (`test-cutechess`) provides a measure of playing strength and helps identify regressions or improvements from changes. The `Makefile` provides targets for these tests.
- **Debugging:** The `test-debug` target in the `Makefile` facilitates focused debugging sessions with `cutechess-cli`.

---

## 5. 🔬 Technical Deep Dive & Trade-offs

### 5.1. Board Representation (`board.go`)

- **Bitboards:** The board state is represented using 64-bit integers (bitboards) for each piece type and color (e.g., `WhitePawns`, `BlackKnights`). Each bit corresponds to a square, enabling fast piece lookup and efficient move generation.
- **Castling Rights:** Castling availability is managed with a dedicated struct, storing four booleans for each possible castling right (white/black, king/queen side).
- **Turn, En Passant, and Move Counters:**
  - `WhiteToMove` (bool): Indicates which side is to move.
  - `OnPassant` (byte): Stores the en passant target square index (0 if not available).
  - `HalfMoveClock` and `FullMoveCounter`: Track the 50-move rule and move number.
- **FEN Support:**
  The board can be initialized from and exported to FEN, supporting all standard fields (piece placement, turn, castling, en passant, clocks).
- **Utility Methods:**
  Includes methods for piece lookup, move application, cloning, and board printing.
- **Design Rationale:**
  This design leverages bitboards for performance and clarity, and is extensible for future features.

### 5.2. Move Generation (`generate.go`)

- **Pseudo-Legal to Legal:** Moves are typically generated as pseudo-legal (ignoring checks to the king) and then validated.
- **Efficiency:** Techniques like pre-calculated attack tables for sliding pieces and knight moves are common. Libra's current approach is direct computation, which can be optimized further.

### 5.3. Transposition Table (`tt.go` & `zobrist.go`)

- **Zobrist Hashing:** Each position is mapped to a unique (with high probability) 64-bit hash key. Keys are updated incrementally when moves are made/unmade.
- **Table Structure:** Typically a hash map or a large array with a simple indexing scheme (e.g., `hash % table_size`).
- **Stored Information:** Each entry stores the hash key (for collision detection), depth of the search, score, score type (exact, lower bound, upper bound), and best move.
  - _Trade-off:_ The amount of information stored per entry affects TT size and the utility of hits. More info is better but costs memory.

### 5.4. Concurrency Strategy

- **Parallel Move Evaluation:** The core search algorithm leverages Go's concurrency by evaluating top-level moves in parallel. The `Search` function distributes each legal move at the root to a worker goroutine, allowing multiple positions to be searched simultaneously and efficiently utilizing all available CPU cores.
- **Worker Pool:** The number of worker goroutines is determined by the number of logical CPUs (`runtime.GOMAXPROCS(0)`), ensuring optimal parallelism on the host system.
- **Result Aggregation:** Results from all workers are collected and the best move is selected according to the search score, preserving move ordering for tie-breaking.
- **Thread Safety:** Each worker operates on a cloned board state, ensuring thread safety and correctness.
- **UCI Communication:** UCI command handling can still be performed in a separate goroutine to keep the engine responsive during search.

This approach provides significant speedup for the root search and is a foundation for further parallelism in deeper search layers in the future.

---

## 6. 🛣️ Roadmap & Future Enhancements

Libra Chess is an actively evolving project. Key areas for future development include:

- **Search Enhancements:**
  - **Principal Variation Search (PVS):** Implement PVS for more efficient search.
  - **Late Move Reductions (LMR):** Reduce search depth for moves ordered later.
  - **Futility Pruning & Razoring:** More aggressive pruning techniques.
  - **Null Move Pruning:** A powerful pruning technique.
  - **Improved Quiescence Search:** More robust handling of tactical positions.
- **Evaluation Refinements:**
  - **Advanced Positional Terms:** Incorporate pawn structure analysis, king safety, mobility scores, passed pawn evaluation.
  - **Tapered Evaluation:** Smoothly transition PSTs and other eval terms from middlegame to endgame.
  - **Automated Tuning:** Explore techniques like CLOP (Chess ELO Optimizer) for tuning evaluation parameters.
- **Time Management:** More sophisticated algorithms to allocate time effectively across moves.
- **Opening Book:** Develop or integrate a more comprehensive internal opening book format.
- **Endgame Tablebases:** Integrate support for Syzygy or Gaviota tablebases for perfect play in endgames.
- **Concurrency in Search:** Leverage Go's goroutines to implement parallel search algorithms (e.g., Lazy SMP or ABDADA).
- **UCI Options:** Expose more internal parameters (e.g., Hash size, contempt factor) via UCI options.
- **Continuous Integration/Delivery:** Enhance CI pipeline for automated testing and releases.

---

## 7. 🤝 Contributing

We welcome contributions of all kinds! Whether you're a chess enthusiast, Go developer, or just curious, your ideas, bug reports, code improvements, and questions are valuable to us.

- **How to contribute:**
  - Open an [Issue](https://github.com/eugenioenko/libra-chess/issues) for bugs, feature requests, or questions.
  - Submit a Pull Request for code, documentation, or test improvements.
  - Suggest enhancements or discuss design ideas.
  - Help with testing, benchmarking, or documentation.

If you're unsure where to start, feel free to ask!

---

## 8. 📄 License

This project is licensed under the [MIT License](./LICENSE) (or specify your chosen license).
