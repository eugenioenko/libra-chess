[![Go Report Card](https://goreportcard.com/badge/github.com/eugenioenko/libra-chess)](https://goreportcard.com/report/github.com/eugenioenko/libra-chess)

# ⚖️ Libra Chess Engine

> A UCI-Compliant Chess Engine in Go
>
> ## [Play against Libra Chess live!](https://eugenioenko.github.io/libra-chess-ui)

---

## 1. 📝 Overview

Libra Chess is a UCI (Universal Chess Interface) compliant chess engine written in Go. The primary goal of Libra is to achieve a balance between high performance, modern software architecture, and clarity of design. This project serves as an exploration of chess engine development leveraging Go's unique strengths in concurrency, tooling, and efficient compilation.

This engine is designed for chess enthusiasts, developers looking to understand chess engine internals, and as a demonstration of software engineering principles applied to a complex domain.

The engine also compiles to WebAssembly (WASM), allowing it to run entirely in the browser. You can [play against Libra Chess live](https://eugenioenko.github.io/libra-chess-ui) in the web interface.

**Estimated strength: ~1800 ELO** (tested against Stockfish at various ELO levels using cutechess-cli).

---

## 2. ✨ Key Features

- **UCI Protocol Compliant:** Seamless integration with popular UCI-compatible GUIs (e.g., CuteChess, CoreChess, PyChess). Supports `wtime`, `btime`, `winc`, `binc`, `movestogo`, `movetime`, `depth`, `infinite`, and `stop`.
- **Alpha-Beta Search with Quiescence:** Alpha-Beta pruning with quiescence search at leaf nodes to resolve tactical sequences and avoid the horizon effect.
- **Iterative Deepening:** Progressive deepening with soft/hard time limits for flexible time management.
- **Transposition Table:** Zobrist hashing with bound types (exact, lower, upper) for effective position caching and search cutoffs.
- **Tapered Evaluation:** PeSTO piece-square tables with middlegame/endgame interpolation based on game phase, providing phase-aware positional understanding.
- **Move Ordering:** TT move, MVV-LVA captures, killer moves, and history heuristic for efficient alpha-beta pruning.
- **Parallel Root Search:** Distributes root moves across worker goroutines using all available CPU cores.
- **Endgame Heuristics:** King proximity bonus in endgames to encourage mating with material advantage.
- **WASM Build:** Compiles to WebAssembly, enabling the engine to run entirely in the browser. Powers the [live web interface](https://eugenioenko.github.io/libra-chess-ui).
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
make build        # Native binary: ./libra-chess
make build-wasm   # WASM binary: wasm/libra.wasm
make build-release # Cross-compile for all platforms
```

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
  _(Requires Stockfish CLI to be at `./dist/stockfish/stockfish-cli`)_
- **Fast ELO Estimate:**
  ```bash
  make test-elo
  ```
  _(40 games vs Stockfish 1500, 10+1)_
- **Debug Match:**
  ```bash
  make test-debug
  ```

---

## 4. 🏛️ Architectural Overview & Design Philosophy

Libra Chess is architected with modularity, simplicity, and maintainability as primary considerations. The choice of Go as the implementation language was deliberate, aiming to harness its excellent concurrency primitives, straightforward syntax, and robust standard library for an extra boost in performance.

### 4.1. Language Choice: Go

- **Advantages:**
  - **Concurrency:** Go's goroutines and channels offer a powerful yet simple model for concurrent programming, used for parallel root search and search cancellation.
  - **Performance:** While not C/C++, Go offers impressive performance, especially with its efficient garbage collector (GC) and direct compilation to machine code. Careful memory management is still crucial.
  - **Simplicity & Readability:** Go's clean syntax and established conventions promote maintainable and understandable code.
  - **Tooling:** Rich ecosystem including `gofmt` for automated formatting, `go test` for testing, and `golangci-lint` for static analysis.
- **Trade-offs & Considerations:**
  - **Garbage Collection:** While Go's GC is highly optimized, in performance-critical sections like deep search, GC pauses can be a concern. This necessitates careful memory allocation patterns (Move/UndoMove instead of cloning).
  - **Ecosystem for Chess Engines:** The C++ ecosystem for chess engines is more mature with a larger pool of shared libraries and knowledge. Libra aims to contribute to the growing Go presence in this domain.

### 4.2. Modularity

The engine's core logic is organized within the `pkg/` directory, with separation of concerns:

- `board.go`: Board representation, piece management, and core game state.
- `evaluate.go`: Static evaluation function, including tapered PeSTO evaluation and endgame heuristics.
- `generate.go`: Move generation logic (legal moves and capture-only generation for quiescence).
- `search.go`: Search algorithms (Alpha-Beta, quiescence search, iterative deepening, parallel root search).
- `sort.go`: Move ordering (TT move, MVV-LVA, killer moves, history heuristic).
- `tt.go`: Transposition table implementation with bound types.
- `zobrist.go`: Zobrist hashing for position keys.
- `move.go`: Move/UndoMove for calculations.
- `piece.go`: Pieces definitions.
- `const.go`: Constants, piece-square tables, and phase values.
- `uci.go`: UCI time management and command parsing.
- `utils.go`: Utility and data structure definitions.

This modular design facilitates easier testing, debugging, and future feature development.

### 4.3. Evaluation Function Design (`evaluate.go`)

The evaluation uses a tapered PeSTO approach — chosen over simpler material-only evaluation because PeSTO tables are well-studied, require no tuning infrastructure, and provide both material and positional scoring in a single lookup. Alternatives like Texel-tuned custom tables were deferred since they require a large labeled game dataset and tuning framework that aren't justified at this stage.

- **Components:**
  - **Tapered PeSTO Evaluation:** Separate middlegame and endgame piece-square tables for all piece types, with material values baked in at startup. The game phase is computed from remaining pieces (knights, bishops, rooks, queens) and used to interpolate between MG and EG scores.
  - **Endgame Heuristics:** King proximity bonus encourages the stronger side's king to approach the opponent's king when material is low, promoting checkmates in won endgames.
- **Trade-offs:**
  - PeSTO tables provide a strong positional baseline with zero tuning cost, but lack awareness of pawn structure, king safety, and piece mobility — terms that require per-position computation and would slow evaluation.
  - Adding evaluation complexity has diminishing returns without search improvements to reach the positions where it matters. Search depth was prioritized first.

### 4.4. Search Algorithm (`search.go`)

- **Alpha-Beta with Iterative Deepening:** Searches to progressively deeper depths, using soft time limits (stop deepening) and hard time limits (abort in-flight search). Move ordering from previous iterations improves pruning at each new depth.
- **Quiescence Search:** At leaf nodes, extends the search for all captures and promotions until the position is quiet, using stand-pat evaluation and MVV-LVA ordering. Without quiescence, the engine would evaluate positions mid-exchange and make severe tactical blunders.
- **Move Ordering:** TT move first, then MVV-LVA captures, killer moves, and history heuristic. Good move ordering is the single biggest factor in alpha-beta efficiency — the difference between searching 10x more or fewer nodes in the same time.
- **Trade-offs:**
  - Parallel root search clones the board for each worker, trading memory for thread safety. This avoids lock contention entirely but means interior nodes can't share pruning information across threads — a known limitation that Lazy SMP would address.
  - The TT uses a Go `map` with `sync.RWMutex`, which is simple and correct but has GC pressure and cache-unfriendly access patterns compared to a fixed-size array. This is a deliberate simplicity-first choice; profiling shows it's not yet the bottleneck.

### 4.5. Testing and Validation

- **Correctness:** `perft` tests validate move generation against known node counts at each depth. Unit tests cover evaluation, search, and move generation.
- **Strength Regression:** Every change is validated through head-to-head matches against the previous version (`make test-cutechess`) and against Stockfish at controlled ELO levels (`make test-elo`). A change that doesn't win more than it loses doesn't ship.
- **Methodology:** See [Design Decisions & Measured Impact](#55-design-decisions--measured-impact) for detailed results.

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

- **Pseudo-Legal to Legal:** Moves are generated as pseudo-legal (ignoring pins and checks) then filtered through `IsMoveLegal()` which applies the move and checks if the king is attacked. This is simpler than maintaining attack maps but means illegal moves are generated and discarded.
- **Precomputed Tables:** `RookRays`, `BishopRays`, `KnightOffsets`, `KingOffsets`, and `SquaresToEdge` are computed once at startup, avoiding repeated calculation during search.
- **Capture Generation:** `GenerateLegalCaptures()` generates only captures and promotions for quiescence search, avoiding the cost of generating quiet moves at leaf nodes.

### 5.3. Transposition Table (`tt.go` & `zobrist.go`)

- **Zobrist Hashing:** Each position is mapped to a 64-bit hash key. Currently computed from scratch each lookup rather than incrementally updated — simpler to implement correctly but slower. Incremental updates are a future optimization.
- **Table Structure:** Go `map[uint64]TTEntry` with `sync.RWMutex`. A fixed-size array (`hash % size`) would be more cache-friendly and avoid GC pressure, but the map approach was chosen for correctness-first development.
- **Bound Types:** Each entry stores the score's relationship to the search window — exact (PV node), lower bound (beta cutoff), or upper bound (failed low). This allows the TT to produce cutoffs even when the stored score isn't exact, which dramatically increases hit rates.
- **Replacement Policy:** Depth-preferred — only overwrites if the new search depth >= stored depth. This preserves the most valuable (deepest) results but can cause the table to fill with stale deep entries over time. Age-based replacement would address this.

### 5.4. Concurrency Strategy

- **Parallel Root Search:** Root moves are distributed to a worker pool sized to `runtime.GOMAXPROCS(0)`. Each worker clones the board and searches independently — no shared mutable state, no lock contention during search.
- **Trade-off:** Cloning per worker means interior nodes can't share alpha-beta bounds across threads. This is less efficient than Lazy SMP (where threads share the TT and occasionally duplicate work but benefit from different move orderings). However, the clone approach is simpler, correct by construction, and avoids subtle concurrency bugs.
- **Cancellation:** Search goroutines listen on a `Done` channel for timeouts and UCI `stop` commands. The UCI loop runs in a separate goroutine so the engine remains responsive during search.

### 5.5. Design Decisions & Measured Impact

Each major feature was validated through head-to-head matches (10 games, 30+0 time control) against the previous version using `cutechess-cli` before merging.

| Change | Result vs. Previous | Impact |
| ------ | ------------------- | ------ |
| **Tapered PeSTO Evaluation** (replacing simple material + basic PSTs) | 8-0-4 | Largest single ELO gain. PeSTO tables provide both material and positional scoring with phase-aware interpolation, replacing ~10 lines of material counting with a well-studied lookup scheme. |
| **Quiescence Search** (capture search at leaf nodes) | 6-0-4 | Eliminated horizon effect — the engine no longer evaluates positions mid-exchange. Stand-pat evaluation with MVV-LVA ordered captures. |
| **TT Bound Types** (exact/lower/upper instead of exact-only) | 8-1-1 | Most TT entries are bounds, not exact scores. Without bound types, the vast majority of TT hits were wasted. Search test suite dropped from 9.1s to 6.2s. |
| **UCI Time Management** (soft/hard limits, increment, movestogo) | — | No strength change in self-play, but critical for tournament play — the engine previously used a hardcoded 1s/move regardless of time control. |

**Methodology:** Each change is tested in isolation against the immediately prior version. Self-play results are directional (small sample), so major changes are also validated against Stockfish at fixed ELO levels for absolute strength estimation.

### 5.6. Move Generation Benchmarks

| perft(N) | Time per op (μs) | Legal moves calculated |
| -------- | ----------------- | ---------------------- |
| 1        | 2.81              | 20                     |
| 2        | 37.85             | 400                    |
| 3        | 274.88            | 8,902                  |
| 4        | 6512.77           | 197,281                |
| 5        | 142827.80         | 4,865,609              |
| 6        | 3865875.06        | 119,060,324            |

> cpu: AMD Ryzen 7 6800H — approximately **30 million moves per second**.

### 5.6. Playing Strength

| Opponent          | Games | Score      | Win Rate | Est. ELO Diff |
| ----------------- | ----- | ---------- | -------- | ------------- |
| Stockfish 1500    | 40    | 31-9-0     | 77.5%    | +215          |
| Stockfish 1800    | 10    | 7-2-1      | 75.0%    | +191          |
| Stockfish 2000    | 10    | 2-8-0      | 20.0%    | -241          |

Estimated playing strength: **~1800 ELO**.

---

## 6. 🛣️ Roadmap

Prioritized by expected ELO impact relative to implementation complexity. Items higher on the list have better strength-to-effort ratios based on results from other engines at similar rating ranges.

### Phase 1: Search Depth (highest impact)

These techniques let the engine search deeper in the same time by pruning more of the tree. At ~1800 ELO, search depth is the primary bottleneck.

- **Null Move Pruning (~50-100 ELO):** In non-zugzwang positions, skip a move and search with reduced depth. If the opponent still can't beat beta, the position is so good that a full search is unnecessary. Cheap to implement, large pruning gains.
- **Late Move Reductions (~50-80 ELO):** Moves ordered late by the move ordering heuristic are unlikely to be best. Search them at reduced depth first and only re-search at full depth if they surprise. Synergizes with good move ordering — which is already in place.
- **Principal Variation Search (~20-40 ELO):** Search the first move (expected best from TT/move ordering) with a full window and all remaining moves with a zero window. Re-search on fail-high. Effective when move ordering is good.

### Phase 2: Search Efficiency

- **Aspiration Windows:** Start each iterative deepening iteration with a narrow window around the previous score. Most iterations confirm the score, saving work. Re-search with a wider window on fail.
- **Check Extensions:** Extend search by one ply when in check, since check positions are tactically sharp and shouldn't be cut short by depth limits.
- **SEE (Static Exchange Evaluation):** Evaluate capture sequences without actually searching them. Replaces MVV-LVA for capture ordering and enables pruning of clearly losing captures in quiescence.

### Phase 3: Evaluation Refinement

Search improvements plateau without better evaluation to guide the search.

- **Pawn Structure:** Penalize doubled and isolated pawns, bonus for passed pawns. High impact in endgames where pawn structure determines the outcome.
- **Positional Terms:** Bishop pair bonus, rook on open file, king safety, mobility.
- **Automated Tuning:** Once enough evaluation terms exist, use Texel tuning or similar to optimize weights against a corpus of games.

### Phase 4: Infrastructure & Correctness

- **Draw Detection:** Repetition and 50-move rule detection. Currently the engine can't detect draws, which causes it to shuffle pieces in drawn endgames instead of seeking other plans.
- **Fixed-Size TT Array:** Replace `map[uint64]TTEntry` with a fixed-size slice indexed by `hash % size`. Eliminates GC pressure, improves cache locality, and allows memory budget control via UCI `Hash` option.
- **Incremental Zobrist Updates:** Update the hash incrementally on Move/UndoMove instead of recomputing from scratch. Reduces hashing cost from O(pieces) to O(1) per move.
- **Endgame Tablebases:** Syzygy tablebase support for perfect play in positions with ≤6 pieces.

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

This project is licensed under the [MIT License](./LICENSE)
