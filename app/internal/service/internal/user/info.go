package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	g "juejin/app/global"
	"juejin/app/internal/model/user"
	"strconv"
)

type SInfo struct{}

var insInfo = SInfo{}

func (s *SInfo) GetUserInfo(ctx context.Context, userBasic *user.Basic, userCounter *user.Counter, id any) error {
	sqlStr := "select digg_article_count,digg_shortmsg_count,followee_count,follower_count,got_digg_count,got_view_count,post_article_count,post_shortmsg_count,select_online_course_count,collection_set_count from user_counter where user_id = ?"
	err := g.MysqlDB.QueryRow(sqlStr, id).Scan(
		&userCounter.DiggArticleCount,
		&userCounter.DiggShortmsgCount,
		&userCounter.FolloweeCount,
		&userCounter.FollowerCount,
		&userCounter.GotDiggCount,
		&userCounter.GotViewCount,
		&userCounter.PostArticleCount,
		&userCounter.PostShortmsgCount,
		&userCounter.SelectOnlineCourseCount,
		&userCounter.CollectionSetCount)
	if err != nil {
		g.Logger.Error("get user counter error", zap.Error(err))
		return err
	}

	err = getUserCounterCache(ctx, userCounter, id)
	if err != nil {
		return err
	}

	sqlStr = "select username,description,avatar,company,job_title from user_basic where user_id=?"
	err = g.MysqlDB.QueryRow(sqlStr, id).Scan(&userBasic.Username, &userBasic.Description, &userBasic.Avatar, &userBasic.Company, &userBasic.JobTitle)
	if err != nil {
		g.Logger.Error("get user basic error", zap.Error(err))
		return err
	}
	return nil
}

func (s *SInfo) UpdateUserInfo(userBasic *user.Basic, id any) error {
	sqlStr := "update user_basic set avatar=?,description=?,company=?,job_title=?  where user_id=?"
	_, err := g.MysqlDB.Exec(sqlStr, userBasic.Avatar, userBasic.Description, userBasic.Company, userBasic.JobTitle, id)
	if err != nil {
		g.Logger.Error("update user info error", zap.Error(err))
		return err
	}
	return nil
}

func (s *SInfo) GetUserBasic(userBasic *user.Basic, id any) error {
	sqlStr := "select description,avatar,company,job_title from user_basic where user_id=?"
	err := g.MysqlDB.QueryRow(sqlStr, id).Scan(&userBasic.Description, &userBasic.Avatar, &userBasic.Company, &userBasic.JobTitle)
	if err != nil {
		g.Logger.Error("get user basic error", zap.Error(err))
		return err
	}
	return nil
}

func getUserCounterCache(ctx context.Context, u *user.Counter, id any) error {
	key := "user_counter"
	field1 := fmt.Sprintf("{%d:digg_article_count}", id)
	field2 := fmt.Sprintf("{%d:got_digg_count}", id)
	field3 := fmt.Sprintf("{%d:got_view_count}", id)
	ok, err := g.Rdb.HExists(ctx, key, field1).Result()
	if err != nil {
		g.Logger.Error("'check exist error", zap.Error(err))
		return err
	}
	if !ok {
		g.Rdb.HSet(ctx, key, field1, u.DiggArticleCount)
	}

	ok, err = g.Rdb.HExists(ctx, key, field2).Result()
	if err != nil {
		g.Logger.Error("'check exist error", zap.Error(err))
		return err
	}
	if !ok {
		g.Rdb.HSet(ctx, key, field2, u.GotDiggCount)
	}

	ok, err = g.Rdb.HExists(ctx, key, field3).Result()
	if err != nil {
		g.Logger.Error("'check exist error", zap.Error(err))
		return err
	}
	if !ok {
		g.Rdb.HSet(ctx, key, field3, u.GotDiggCount)
	}

	res, err := g.Rdb.HMGet(ctx, key, field1, field2, field3).Result()
	if err != nil {
		g.Logger.Error("get user counter cache error", zap.Error(err))
		return err
	}
	var cache = make([]int, 3)
	for k, v := range res {
		cache[k], _ = strconv.Atoi(v.(string))
	}
	u.DiggArticleCount = cache[0]
	u.GotDiggCount = cache[1]
	u.GotViewCount = cache[2]
	return nil
}
