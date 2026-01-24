package health

type Health struct {
	max     int
	current int
}

func Init(max int) *Health {
	health := &Health{
		max:     max,
		current: max,
	}

	return health
}

func (h *Health) ChangeMax(max int) {
	h.max = max
}

func (h *Health) Reset() {
	h.current = h.max
}

func (h *Health) TakeDamage(base int) {
	h.current -= base
}

func (h *Health) GetCurrent() int {
	return h.current
}

func (h *Health) GetMax() int {
	return h.max
}
