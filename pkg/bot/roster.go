package bot

import "time"

// User represents a twitch user
type User struct {
	Nick      string `json:"nick"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (u *User) FirstSeen() time.Time {
	return time.Unix(u.CreatedAt, 0)
}

func (u *User) LastSeen() time.Time {
	return time.Unix(u.UpdatedAt, 0)
}
