package systems

type Health struct {
	max     int
	current int
}

func InitHealth(max int) *Health {
	health := &Health{
		max:     max,
		current: max,
	}

	return health
}

func (h *Health) TakeDamage(base int) {
	h.current -= base
}

func (h *Health) GetCurrent() int {
	return h.current
}
