package db

func Init() {
	err := DB.AutoMigrate(&Server{}, &Backup{}, &Log{})
	if err != nil {
		panic("failed to migrate database schema: " + err.Error())
	}
}
