package model

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	NIK      string `gorm:"size:50;unique" json:"nik"`
	Username string `gorm:"size:100;unique" json:"username"`
	Email    string `gorm:"size:100;unique" json:"email"`
	Password string `gorm:"size:255" json:"-"` // hidden saat JSON response
}

// Override default table name (opsional, kalau table kamu bernama "users")
func (User) TableName() string {
	return "users"
}
