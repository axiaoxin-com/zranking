package zranking

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestZRanking(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.TODO()
	key := "zranking:1001_1031"
	defer rds.Del(ctx, key)

	var uid int64 = 54088
	var startTs int64 = 1667549955
	var endTs int64 = 1670141955
	zr, err := New(rds, key, startTs, endTs, time.Hour*24*30)
	require.Nil(t, err)

	score, err := zr.Update(ctx, uid, 100)
	require.Nil(t, err)

	time.Sleep(1 * time.Second)
	score1, err := zr.Update(ctx, uid-1, 200)
	require.Nil(t, err)

	time.Sleep(1 * time.Second)
	score2, err := zr.Update(ctx, uid+1, 100)
	require.Nil(t, err)
	t.Logf("score:%v, score1:%v, score2:%v", score, score1, score2)

	total := zr.GetTotalCount(ctx)
	require.Equal(t, int64(3), total)

	val, err := zr.score2val(ctx, score)
	require.Nil(t, err)
	require.Equal(t, int64(100), val)

	val1, err := zr.score2val(ctx, score1)
	require.Nil(t, err)
	require.Equal(t, int64(200), val1)

	val2, err := zr.score2val(ctx, score2)
	require.Nil(t, err)
	require.Equal(t, int64(100), val2)

	rank, err := zr.GetRankingList(ctx, 3, true)
	require.Nil(t, err)
	require.Len(t, rank, 3)
	require.Equal(t, uid-1, rank[0].UID)
	require.Equal(t, uid, rank[1].UID)
	require.Equal(t, uid+1, rank[2].UID)

	rank, err = zr.GetRankingList(ctx, 2, true)
	require.Nil(t, err)
	require.Len(t, rank, 2)
	require.Equal(t, uid-1, rank[0].UID)
	require.Equal(t, uid, rank[1].UID)

	rank, err = zr.GetRankingList(ctx, 4, true)
	require.Nil(t, err)
	require.Len(t, rank, 3)

	rank, err = zr.GetRankingList(ctx, 0, true)
	require.Nil(t, err)
	require.Len(t, rank, 3)

	rank, err = zr.GetRankingList(ctx, 0, false)
	require.Nil(t, err)
	require.Len(t, rank, 3)
	require.Equal(t, uid-1, rank[2].UID)
	require.Equal(t, uid, rank[1].UID)
	require.Equal(t, uid+1, rank[0].UID)

	ur, err := zr.GetUserRank(ctx, uid, true)
	require.Nil(t, err)
	require.Equal(t, int64(1), ur)
	uv, err := zr.GetUserVal(ctx, uid)
	require.Nil(t, err)
	require.Equal(t, int64(100), uv)

	ur, err = zr.GetUserRank(ctx, uid-1, true)
	require.Nil(t, err)
	require.Equal(t, int64(0), ur)
	uv, err = zr.GetUserVal(ctx, uid-1)
	require.Nil(t, err)
	require.Equal(t, int64(200), uv)

	ur, err = zr.GetUserRank(ctx, uid+1, true)
	require.Nil(t, err)
	require.Equal(t, int64(2), ur)
	uv, err = zr.GetUserVal(ctx, uid+1)
	require.Nil(t, err)
	require.Equal(t, int64(100), uv)

	ur, err = zr.GetUserRank(ctx, uid+1, false)
	require.Nil(t, err)
	require.Equal(t, int64(0), ur)
	uv, err = zr.GetUserVal(ctx, uid+1)
	require.Nil(t, err)
	require.Equal(t, int64(100), uv)

	ur, err = zr.GetUserRank(ctx, uid-1, false)
	require.Nil(t, err)
	require.Equal(t, int64(2), ur)
	uv, err = zr.GetUserVal(ctx, uid-1)
	require.Nil(t, err)
	require.Equal(t, int64(200), uv)

	score3, err := zr.Update(ctx, uid+1, 100)
	require.Nil(t, err)
	val3, err := zr.score2val(ctx, score3)
	require.Nil(t, err)
	require.Equal(t, int64(200), val3)
	time.Sleep(1 * time.Second)

	score4, err := zr.Update(ctx, uid, 100)
	val4, err := zr.score2val(ctx, score4)
	require.Nil(t, err)
	require.Equal(t, int64(200), val4)
	require.Nil(t, err)

	rank, err = zr.GetRankingList(ctx, 0, true)
	require.Nil(t, err)
	require.Len(t, rank, 3)
	require.Equal(t, uid-1, rank[0].UID)
	require.Equal(t, uid+1, rank[1].UID)
	require.Equal(t, uid, rank[2].UID)

	ur, err = zr.GetUserRank(ctx, uid-1, true)
	require.Nil(t, err)
	require.Equal(t, int64(0), ur)
	uv, err = zr.GetUserVal(ctx, uid-1)
	require.Nil(t, err)
	require.Equal(t, int64(200), uv)

	ur, err = zr.GetUserRank(ctx, uid+1, true)
	require.Nil(t, err)
	require.Equal(t, int64(1), ur)
	uv, err = zr.GetUserVal(ctx, uid+1)
	require.Nil(t, err)
	require.Equal(t, int64(200), uv)

	ur, err = zr.GetUserRank(ctx, uid, true)
	require.Nil(t, err)
	require.Equal(t, int64(2), ur)
	uv, err = zr.GetUserVal(ctx, uid)
	require.Nil(t, err)
	require.Equal(t, int64(200), uv)
}
