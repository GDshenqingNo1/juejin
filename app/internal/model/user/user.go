package user

import "time"

type InfoPack struct {
	Counter `json:"user_counter"`
	Basic   `json:"user_basic"`
}

type Auth struct {
	Id         int64     `json:"id" form:"id" db:"id"`
	Username   string    `json:"username" form:"username" db:"username"`
	Password   string    `json:"password" form:"password" db:"password"`
	Email      string    `json:"email" form:"email" db:"email"`
	Phone      string    `json:"phone" form:"phone" db:"phone" `
	Code       string    `json:"code" form:"code" db:"-"`
	GithubId   int64     `json:"github_id" form:"github_id" db:"github_id"`
	CreateTime time.Time `json:"create_time" form:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" form:"update_time" db:"update_time"`
}

type Basic struct {
	Id          int64  `json:"-" form:"id" db:"id"`
	Username    string `json:"username" form:"username" db:"username"`
	Avatar      string `json:"avatar" form:"avatar" db:"avatar"`
	Company     string `json:"company" form:"company" db:"company"`
	Description string `json:"description" form:"description" db:"description"`
	JobTitle    string `json:"job_title" form:"job_title" db:"job_title"`
	IsFollow    bool   `json:"is_follow"`
}

type Counter struct {
	Id                      int64 `json:"-" form:"id" db:"id"`
	DiggArticleCount        int   `json:"digg_article_count" form:"digg_article_count" db:"digg_article_count"`
	DiggShortmsgCount       int   `json:"digg_shortmsg_count" form:"digg_shortmsg_count" db:"digg_shortmsg_count"`
	FolloweeCount           int   `json:"followee_count" form:"followee_count" db:"followee_count"`
	FollowerCount           int   `json:"follower_count" form:"follower_count" db:"follower_count"`
	GotDiggCount            int   `json:"got_digg_count" form:"got_digg_count" db:"got_digg_count"`
	GotViewCount            int   `json:"got_view_count" form:"got_view_count" db:"got_view_count"`
	PostArticleCount        int   `json:"post_article_count" form:"post_article_count" db:"post_article_count"`
	PostShortmsgCount       int   `json:"post_shortmsg_count" form:"post_shortmsg_count" db:"post_shortmsg_count"`
	SelectOnlineCourseCount int   `json:"select_online_course_count" form:"select_online_course_count" db:"select_online_course_count"`
	CollectionSetCount      int   `json:"collection_set_count" form:"collection_set_count"`
}
