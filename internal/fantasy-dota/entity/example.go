package entity

type ExampleModel struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func (m *ExampleModel) TableName() string {
	return "example"
}
