package channels

import "context"

// Pipeline represents a channel backed pipeline, with a given start and end channel.  This type is useful for ensuring
// that a given pipeline starts and ends with a given type, but the operations which occur in the middle of the pipeline
// (i.e. how an input is converted into the required output) are not specified.
type Pipeline[I, O any] struct {
	start <-chan I
	end   <-chan O
}

// PipelineCreationFunc is a function which takes a context and a channel of the input type and returns a channel of the
// output type. The context should be threaded through every stage the function builds so the whole pipeline can be
// cancelled as a unit.
type PipelineCreationFunc[I, O any] func(ctx context.Context, input <-chan I) <-chan O

// NewPipeline creates a new Pipeline, with the given input channel and PipelineCreationFunc.  The PipelineCreationFunc
// is used to create the end channel of the pipeline.
//
// The supplied context is passed to the PipelineCreationFunc so it can be threaded through every stage; cancelling it
// tears the whole pipeline down and reclaims its goroutines deterministically.
func NewPipeline[I, O any](ctx context.Context, input <-chan I, fn PipelineCreationFunc[I, O]) *Pipeline[I, O] {
	end := fn(ctx, input)
	return &Pipeline[I, O]{
		start: input,
		end:   end,
	}
}

// CollectAsSlice collects all elements from the end channel of the pipeline into a slice, which is returned.  This
// function will block until the end channel is closed.
func (p Pipeline[I, O]) CollectAsSlice() []O {
	return CollectAsSlice(p.end)
}
