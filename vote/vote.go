package vote

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"redis/client"
	"strconv"
	"time"
)

const OneWeekSeconds = 7 * 86400
const Score float64 = 432

var rdb = client.Rdb

type Article struct {
	ArticleId string
	Title     string
	Link      string
	UserId    string
	Time      string
}

func (article Article) ArticleVote(userId string) {
	//获取当前的unix时间
	cutoff := time.Now().Unix() - OneWeekSeconds
	var articleValue = "article:" + article.ArticleId

	score, err := rdb.ZScore(context.Background(), "time", articleValue).Result()
	if err != nil {
		log.Fatalf("error zscore, %s", err)
	}
	//时间超过一周不在评分
	if score < float64(cutoff) {
		return
	}
	result, err := rdb.SAdd(context.Background(), "voted:"+article.ArticleId, userId).Result()
	if err != nil {
		log.Fatalf("error sadd user vote, %s", err)
	}
	if result > 0 {
		rdb.ZIncrBy(context.Background(), "score", Score, articleValue)
		rdb.HIncrBy(context.Background(), articleValue, "votes", 1)
	}
}

func (article Article) PostArticle() int64 {
	articleId, err := rdb.Incr(context.Background(), "article").Result()
	if err != nil {
		log.Fatalf("gennerate aritcle id error, %s", err)
	}
	now := time.Now().Unix()

	articleKey := "article:" + strconv.FormatInt(articleId, 10)
	data := map[string]interface{}{
		"title": article.Title,
		"link":  article.Link,
		"time":  now,
		"voted": article.UserId,
	}
	rdb.HMSet(context.Background(), articleKey, data)
	scoreZ := redis.Z{
		Score:  float64(now) + Score,
		Member: articleKey,
	}
	timeZ := redis.Z{
		Score:  float64(now),
		Member: articleKey,
	}

	rdb.ZAdd(context.Background(), "score", scoreZ)
	rdb.ZAdd(context.Background(), "time", timeZ)
	return articleId
}

// GetArticles 获取所有的文章
func GetArticles(page int64, size int64, key string) []Article {
	if key == "" {
		key = "score"
	}
	start := (page - 1) * size
	end := start + size - 1
	result, err := rdb.ZRevRange(context.Background(), key, int64(start), int64(end)).Result()
	if err != nil {
		log.Fatalf("score range error, %s", err)
	}
	var articles []Article
	for _, value := range result {
		articleMap, err := rdb.HGetAll(context.Background(), value).Result()
		if err != nil {
			log.Fatalf("get article error:%s", err)
		}
		if err != nil {
			log.Fatalf("error parse int, %s", err)
		}
		article := Article{
			ArticleId: articleMap["articleId"],
			Title:     articleMap["title"],
			Link:      articleMap["link"],
			UserId:    articleMap["voted"],
			Time:      articleMap["time"],
		}
		articles = append(articles, article)
	}
	return articles
}

// AddRemoveGroups 文章分组
func AddRemoveGroups(articleId string, toAdd []string, toRemove []string) {
	value := "article:" + articleId
	for _, group := range toAdd {
		rdb.SAdd(context.Background(), "groups:"+group, value)
	}
	for _, group := range toRemove {
		rdb.SRem(context.Background(), "groups"+group, value)
	}
}

// GetGroupArticles 获取分组中的文章，使用score排序
func GetGroupArticles(group string, page int64, size int64) []Article {
	//检查是否有对应分类的zset
	key := "score:" + group
	result, err := rdb.Exists(context.Background(), key).Result()
	if err != nil {
		log.Fatalf("error exist:%s", err)
	}
	//不存在键，合并
	if result == 0 {
		rdb.ZInterStore(context.Background(), key, &redis.ZStore{
			Keys:      []string{"groups:" + group, "score"},
			Aggregate: "max",
		})
		rdb.Expire(context.Background(), key, 60*time.Second)
	}
	return GetArticles(page, size, key)
}
