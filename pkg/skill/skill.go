package skill

type Attribute int

const (
	STR Attribute = iota
	INT
	SPD
)

type Skill struct {
	STR     int
	INT     int
	SPD     int
	current int
}

func NewSkill() *Skill {
	return &Skill{
		STR:     1,
		INT:     1,
		SPD:     1,
		current: 1,
	}
}

func (s *Skill) LevelUp(skill Skill) {
	s.STR += skill.STR
	s.INT += skill.INT
	s.SPD += skill.SPD
	s.current += 1
}

func (s *Skill) IncreaseSTR(str int) {
	s.STR += str
}

func (s *Skill) IncreaseINT(intl int) {
	s.INT += intl
}

func (s *Skill) IncreaseSPD(spd int) {
	s.SPD += spd
}

func (s *Skill) DecreaseSTR(str int) {
	s.STR -= str
}

func (s *Skill) DecreaseINT(intl int) {
	s.INT -= intl
}

func (s *Skill) DecreaseSPD(spd int) {
	s.SPD -= spd
}

func (s *Skill) Reset() {
	s.STR = 0
	s.INT = 0
	s.SPD = 0
}

func (s *Skill) GetLevel() int {
	return s.current
}

func (s *Skill) GetTotalSkillPoints() int {
	return s.STR + s.INT + s.SPD
}
