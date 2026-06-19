package preprocessing

// Transformer is the shared contract implemented by every fitted estimator in
// the package. An estimator is first fitted from training data (via its
// concrete, chainable Fit method) and then applies the captured parameters to
// any input via Transform.
//
// Transform is the common surface; Fit is deliberately NOT part of this
// interface because its signature varies by family (a scaler fits on
// []float64, an encoder on []C, a target encoder on categories plus a target),
// and because the concrete Fit methods return the receiver for chaining rather
// than a bool. Code that only needs to apply an already-fitted estimator can
// depend on Transformer; code that fits depends on the concrete type.
//
// In is the element type of the input slice and Out is the whole output type
// (for example []float64 for a scaler, [][]float64 for a one-hot encoder, or
// []int for a label encoder or binner). Transform returns ok == false when the
// estimator has not been fitted, or when the input is otherwise unprocessable
// per the concrete type's documented policy.
type Transformer[In any, Out any] interface {
	Transform(input []In) (Out, bool)
}
