package level

type Attribute int

const (
	STR Attribute = iota
	INT
	SPD
)

type Skills struct {
	STR int
	INT int
	SPD int
}

type LevelData struct {
	Attributes Skills
	XP         int
	Current    int
}

func NewLevelData() *LevelData {
	return &LevelData{
		Attributes: Skills{
			STR: 1,
			INT: 1,
			SPD: 1,
		},
		Current: 1,
		XP:      0,
	}
}

func (l *LevelData) GainXP(xp int) {
	l.XP += xp
}

func (l *LevelData) LevelUp(skill Skills) {
	l.Attributes.STR += skill.STR
	l.Attributes.INT += skill.INT
	l.Attributes.SPD += skill.SPD
	l.Current += 1
}

func (l *LevelData) IncreaseSTR(str int) {
	l.Attributes.STR += str
}

func (l *LevelData) IncreaseINT(intl int) {
	l.Attributes.INT += intl
}

func (l *LevelData) IncreaseSPD(spd int) {
	l.Attributes.SPD += spd
}

func (l *LevelData) DecreaseSTR(str int) {
	l.Attributes.STR -= str
}

func (l *LevelData) DecreaseINT(intl int) {
	l.Attributes.INT -= intl
}

func (l *LevelData) DecreaseSPD(spd int) {
	l.Attributes.SPD -= spd
}

func (l *LevelData) Reset() {
	l.Attributes = Skills{}
	l.Current = 0
	l.XP = 0
}

func (l *LevelData) GetTotalSkillPoints() int {
	return l.Attributes.STR + l.Attributes.INT + l.Attributes.SPD
}
