package database

func (u User) Board(text string) {
	if len(text) == 0 {
		return
	}
	board := Board{
		Text:   text,
		UserId: u.Id,
		User:   u,
	}
	DB.Create(&board)
}
