package libra

import "fmt"

type GoOptions struct {
	WTime     int  // white time remaining (ms)
	BTime     int  // black time remaining (ms)
	WInc      int  // white increment per move (ms)
	BInc      int  // black increment per move (ms)
	MovesToGo int  // moves until next time control (0 = sudden death)
	MoveTime  int  // exact time per move (ms)
	Depth     int  // search to exactly this depth
	Infinite  bool // search until "stop" command
}

const timeManagementSafetyMarginMs = 100

func ParseGoOptions(fields []string) GoOptions {
	opts := GoOptions{}
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "wtime":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.WTime)
			}
		case "btime":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.BTime)
			}
		case "winc":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.WInc)
			}
		case "binc":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.BInc)
			}
		case "movestogo":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.MovesToGo)
			}
		case "movetime":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.MoveTime)
			}
		case "depth":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &opts.Depth)
			}
		case "infinite":
			opts.Infinite = true
		}
	}
	return opts
}

// CalcTimeLimit computes optimal (soft) and maximum (hard) time limits in ms.
// optimalTime: target time per move, used to decide when to stop deepening.
// maxTime: absolute ceiling for in-flight searches.
func (opts *GoOptions) CalcTimeLimit(whiteToMove bool) (optimalTime int, maxTime int) {
	// Fixed time per move
	if opts.MoveTime > 0 {
		return opts.MoveTime, opts.MoveTime
	}

	// Infinite or depth-only: no time constraint
	if opts.Infinite || (opts.Depth > 0 && opts.WTime == 0 && opts.BTime == 0) {
		return 0, 0
	}

	remaining := opts.BTime
	increment := opts.BInc
	if whiteToMove {
		remaining = opts.WTime
		increment = opts.WInc
	}

	// No time info at all: fall back to a sensible default
	if remaining <= 0 {
		return MaxEvaluationTimeMs, MaxEvaluationTimeMs
	}

	// Safety: never use more than remaining - margin
	safeRemaining := remaining - timeManagementSafetyMarginMs
	if safeRemaining < 50 {
		safeRemaining = 50
	}

	if opts.MovesToGo > 0 {
		// Moves until next time control
		optimalTime = safeRemaining/opts.MovesToGo + increment*3/4
		maxTime = safeRemaining / 2
	} else {
		// Sudden death: estimate ~30 moves remaining
		optimalTime = safeRemaining/30 + increment*3/4
		maxTime = safeRemaining / 5
	}

	// Clamp optimal to not exceed max
	if optimalTime > maxTime {
		optimalTime = maxTime
	}

	// Never exceed safe remaining time
	if maxTime > safeRemaining {
		maxTime = safeRemaining
	}
	if optimalTime > safeRemaining {
		optimalTime = safeRemaining
	}

	return optimalTime, maxTime
}
