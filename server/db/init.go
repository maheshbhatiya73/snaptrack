package db

func Init() {
	DB.AutoMigrate(&Backup{}, &User{}, &Log{})
}
