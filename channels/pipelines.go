package channels

type Pipeline[I, O any] struct {
	start <-chan I
	end   <-chan O
}

type PipelineCreationFunc[I, O any] func(input <-chan I) <-chan O

func NewPipeline[I, O any](input <-chan I, fn PipelineCreationFunc[I, O]) *Pipeline[I, O] {
	end := fn(input)
	return &Pipeline[I, O]{
		start: input,
		end:   end,
	}
}

func (p Pipeline[I, O]) CollectAsSlice() []O {
	return CollectAsSlice(p.end)
}
