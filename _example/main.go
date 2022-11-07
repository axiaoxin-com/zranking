package main

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/zranking"
	"github.com/go-redis/redis/v8"
)

func main() {
	// 使用到的redis client
	rds := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 排行榜redis key名称
	key := "zranking:example_1"

	// 排行活动开始-结束时间戳（秒）
	var startTs int64 = 1667232000 // 2022-11-01 00:00:00
	var endTs int64 = 1669823999   // 2022-11-30 23:59:59

	// 排行榜数据保存时长 (活动结束后60天自动过期)
	expiration := time.Hour * 24 * 60

	// 创建zranking实例
	ranking, err := zranking.New(rds, key, startTs, endTs, expiration)
	if err != nil {
		panic(err)
	}

	// 模拟更新榜单：
	ctx := context.TODO()
	// 运行完后删除key
	defer rds.Del(ctx, key)

	// 更新排行榜，用户uid=1 +100分
	_, err = ranking.Update(ctx, 1, 100)
	if err != nil {
		panic(err)
	}
	// 更新排行榜，用户uid=2 +200分
	_, err = ranking.Update(ctx, 2, 200)
	if err != nil {
		panic(err)
	}
	// 更新排行榜，用户uid=3 +150分
	_, err = ranking.Update(ctx, 3, 150)
	if err != nil {
		panic(err)
	}
	// 更新排行榜，用户uid=4 1秒后 +200分
	time.Sleep(1 * time.Second)
	_, err = ranking.Update(ctx, 4, 200)
	if err != nil {
		panic(err)
	}

	// 以分数从大到小排，获取排行榜top3
	var topN int64 = 3   // 取top3，传0取全量
	var desc bool = true // true:降序, false:升序
	list, err := ranking.GetRankingList(ctx, topN, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
	// Output:
	// [{200.2013085 2} {200.2013084 4} {150.2013085 3}]

	// 获取指定用户的排名（下标从0开始）
	// 用户id=1的降序排名
	rank, err := ranking.GetUserRank(ctx, 1, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(rank)
	// Output:
	// 3

	// 获取指定用户的排名分
	val, err := ranking.GetUserVal(ctx, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output:
	// 100

	// 获取排行榜总人数
	total := ranking.GetTotalCount(ctx)
	fmt.Println(total)
	// Output:
	// 4

	// 用户id=1再加100，3个总分为200，用户1最后达到，因此排第三
	_, err = ranking.Update(ctx, 1, 100)
	if err != nil {
		panic(err)
	}

	rank, err = ranking.GetUserRank(ctx, 1, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(rank)
	// Output: 下标从0开始，2表示排第三
	// 2

}
