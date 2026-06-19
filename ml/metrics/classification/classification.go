package classification

// Averaging selects how a per-class metric (precision, recall, F1) is combined
// into a single number across the classes of a multiclass problem.
type Averaging int

const (
	// Macro averages the per-class scores with equal weight, so every class
	// counts the same regardless of how many samples it has. Use it when rare
	// classes matter as much as common ones.
	Macro Averaging = iota
	// Micro pools the true positives, false positives and false negatives
	// across all classes before computing a single score. For single-label
	// classification micro precision, recall and F1 all equal the accuracy.
	Micro
	// Weighted averages the per-class scores weighted by each class's support
	// (its number of true samples), so larger classes count more.
	Weighted
)

// String renders an Averaging as its lower-case name ("macro", "micro",
// "weighted"), or "unknown" for an out-of-range value.
func (a Averaging) String() string {
	switch a {
	case Macro:
		return "macro"
	case Micro:
		return "micro"
	case Weighted:
		return "weighted"
	default:
		return "unknown"
	}
}

// Matrix is a multiclass confusion matrix: the counts of every (true,
// predicted) label pairing observed in a single labelling. Labels are held in
// first-seen order (scanning yTrue then yPred at each position), which keeps the
// matrix deterministic for a given input. The zero Matrix has no labels.
type Matrix[T comparable] struct {
	labels []T
	index  map[T]int
	counts [][]int // counts[t][p]: samples with true == labels[t], pred == labels[p]
}

// Labels returns a copy of the matrix's labels in first-seen order.
func (m Matrix[T]) Labels() []T {
	return append([]T(nil), m.labels...)
}

// Count returns the number of samples whose true label is yTrue and whose
// predicted label is yPred. Labels absent from the matrix count as zero.
func (m Matrix[T]) Count(yTrue, yPred T) int {
	t, okT := m.index[yTrue]
	p, okP := m.index[yPred]
	if !okT || !okP {
		return 0
	}
	return m.counts[t][p]
}

// ConfusionMatrix tallies predicted labels against true labels, returning the
// multiclass confusion Matrix together with an ok flag. Labels may be of any
// comparable type (compared with ==); the label set is the union of those seen
// in yTrue and yPred.
//
// ok is false (and the matrix is the zero Matrix) when yTrue is empty or
// len(yTrue) != len(yPred).
func ConfusionMatrix[T comparable](yTrue, yPred []T) (Matrix[T], bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return Matrix[T]{}, false
	}

	index := make(map[T]int)
	var labels []T
	add := func(v T) {
		_, ok := index[v]
		if !ok {
			index[v] = len(labels)
			labels = append(labels, v)
		}
	}
	for i := range yTrue {
		add(yTrue[i])
		add(yPred[i])
	}

	counts := make([][]int, len(labels))
	for i := range counts {
		counts[i] = make([]int, len(labels))
	}
	for i := range yTrue {
		counts[index[yTrue[i]]][index[yPred[i]]]++
	}
	return Matrix[T]{labels: labels, index: index, counts: counts}, true
}

// Accuracy returns the fraction of samples whose predicted label equals the
// true label, together with an ok flag. ok is false (and the result is 0) when
// yTrue is empty or len(yTrue) != len(yPred).
func Accuracy[T comparable](yTrue, yPred []T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return 0, false
	}

	var correct int
	for i := range yTrue {
		if yTrue[i] == yPred[i] {
			correct++
		}
	}
	return float64(correct) / float64(len(yTrue)), true
}

// classStats holds one class's confusion counts: true positives, false
// positives, false negatives and support (the number of true samples).
type classStats struct {
	tp, fp, fn, support int
}

func precisionOf(s classStats) float64 {
	denom := s.tp + s.fp
	if denom == 0 {
		return 0
	}
	return float64(s.tp) / float64(denom)
}

func recallOf(s classStats) float64 {
	denom := s.tp + s.fn
	if denom == 0 {
		return 0
	}
	return float64(s.tp) / float64(denom)
}

func f1Of(s classStats) float64 {
	p := precisionOf(s)
	r := recallOf(s)
	if p+r == 0 {
		return 0
	}
	return 2 * p * r / (p + r)
}

// perClass derives each class's confusion counts from the matrix.
func perClass[T comparable](m Matrix[T]) []classStats {
	n := len(m.labels)
	out := make([]classStats, n)
	for t := 0; t < n; t++ {
		for p := 0; p < n; p++ {
			c := m.counts[t][p]
			out[t].support += c
			if t == p {
				out[t].tp += c
			} else {
				out[t].fn += c // true t, predicted something else
				out[p].fp += c // predicted p, was actually something else
			}
		}
	}
	return out
}

// aggregate combines per-class scores under the chosen averaging strategy.
// score maps one class's counts to its metric (precision, recall or F1).
func aggregate(classes []classStats, avg Averaging, score func(classStats) float64) (float64, bool) {
	switch avg {
	case Macro:
		var sum float64
		for _, s := range classes {
			sum += score(s)
		}
		return sum / float64(len(classes)), true
	case Weighted:
		// The supports sum to the (non-empty) sample count, so total is always
		// positive here — no divide-by-zero guard is needed.
		var sum, total float64
		for _, s := range classes {
			sum += score(s) * float64(s.support)
			total += float64(s.support)
		}
		return sum / total, true
	case Micro:
		var pooled classStats
		for _, s := range classes {
			pooled.tp += s.tp
			pooled.fp += s.fp
			pooled.fn += s.fn
			pooled.support += s.support
		}
		return score(pooled), true
	default:
		return 0, false
	}
}

// Precision returns the precision of a multiclass labelling — TP / (TP + FP)
// per class, combined under avg — together with an ok flag. Precision answers
// "of the samples predicted to be in a class, what fraction really were?".
//
// A class with no predictions contributes a precision of 0 (rather than a
// divide-by-zero). ok is false (and the result is 0) when yTrue is empty,
// len(yTrue) != len(yPred), or avg is not a recognised strategy.
func Precision[T comparable](yTrue, yPred []T, avg Averaging) (float64, bool) {
	m, ok := ConfusionMatrix(yTrue, yPred)
	if !ok {
		return 0, false
	}
	return aggregate(perClass(m), avg, precisionOf)
}

// Recall returns the recall (sensitivity) of a multiclass labelling — TP /
// (TP + FN) per class, combined under avg — together with an ok flag. Recall
// answers "of the samples that really were in a class, what fraction did we
// find?".
//
// A class with no true samples contributes a recall of 0. ok is false (and the
// result is 0) when yTrue is empty, len(yTrue) != len(yPred), or avg is not a
// recognised strategy.
func Recall[T comparable](yTrue, yPred []T, avg Averaging) (float64, bool) {
	m, ok := ConfusionMatrix(yTrue, yPred)
	if !ok {
		return 0, false
	}
	return aggregate(perClass(m), avg, recallOf)
}

// F1 returns the F1 score of a multiclass labelling — the harmonic mean of
// precision and recall per class, 2·P·R / (P + R), combined under avg —
// together with an ok flag. F1 balances precision and recall in a single
// number.
//
// A class with precision + recall of 0 contributes an F1 of 0. ok is false (and
// the result is 0) when yTrue is empty, len(yTrue) != len(yPred), or avg is not
// a recognised strategy.
func F1[T comparable](yTrue, yPred []T, avg Averaging) (float64, bool) {
	m, ok := ConfusionMatrix(yTrue, yPred)
	if !ok {
		return 0, false
	}
	return aggregate(perClass(m), avg, f1Of)
}

// binaryStats tallies the confusion counts for a single positive label, the
// rest of the labels being treated as negative.
func binaryStats[T comparable](yTrue, yPred []T, positive T) (classStats, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return classStats{}, false
	}

	var s classStats
	for i := range yTrue {
		actualPos := yTrue[i] == positive
		predPos := yPred[i] == positive
		switch {
		case actualPos && predPos:
			s.tp++
		case !actualPos && predPos:
			s.fp++
		case actualPos && !predPos:
			s.fn++
		}
		if actualPos {
			s.support++
		}
	}
	return s, true
}

// PrecisionBinary returns the precision for the designated positive label, all
// other labels being treated as negative, together with an ok flag. It is the
// binary special case of Precision: TP / (TP + FP) for the positive class.
//
// With no positive predictions the precision is 0. ok is false (and the result
// is 0) when yTrue is empty or len(yTrue) != len(yPred).
func PrecisionBinary[T comparable](yTrue, yPred []T, positive T) (float64, bool) {
	s, ok := binaryStats(yTrue, yPred, positive)
	if !ok {
		return 0, false
	}
	return precisionOf(s), true
}

// RecallBinary returns the recall for the designated positive label, all other
// labels being treated as negative, together with an ok flag: TP / (TP + FN)
// for the positive class.
//
// With no positive samples the recall is 0. ok is false (and the result is 0)
// when yTrue is empty or len(yTrue) != len(yPred).
func RecallBinary[T comparable](yTrue, yPred []T, positive T) (float64, bool) {
	s, ok := binaryStats(yTrue, yPred, positive)
	if !ok {
		return 0, false
	}
	return recallOf(s), true
}

// F1Binary returns the F1 score for the designated positive label, all other
// labels being treated as negative, together with an ok flag: the harmonic mean
// of the positive class's precision and recall.
//
// With precision + recall of 0 the F1 is 0. ok is false (and the result is 0)
// when yTrue is empty or len(yTrue) != len(yPred).
func F1Binary[T comparable](yTrue, yPred []T, positive T) (float64, bool) {
	s, ok := binaryStats(yTrue, yPred, positive)
	if !ok {
		return 0, false
	}
	return f1Of(s), true
}
