package items

type ItemType int

const (
	Lock ItemType = iota
	Key
)

type Item struct {
	Name string
	ID   int
	Type ItemType
	Link int
}
