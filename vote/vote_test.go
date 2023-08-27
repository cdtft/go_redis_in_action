package vote

import (
	"fmt"
	"testing"
)

func TestArticle_PostArticle(t *testing.T) {
	type fields struct {
		ArticleId string
		Title     string
		Link      string
		UserId    string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			fields: fields{
				Title:  "article_1",
				Link:   "www.xxx.com",
				UserId: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := Article{
				ArticleId: tt.fields.ArticleId,
				Title:     tt.fields.Title,
				Link:      tt.fields.Link,
				UserId:    tt.fields.UserId,
			}
			if got := article.PostArticle(); got != tt.want {
				t.Errorf("PostArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticle_ArticleVote(t *testing.T) {
	type fields struct {
		ArticleId string
		Title     string
		Link      string
		UserId    string
	}
	type args struct {
		userId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			args: args{
				userId: "2",
			},
			fields: fields{
				ArticleId: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := Article{
				ArticleId: tt.fields.ArticleId,
				Title:     tt.fields.Title,
				Link:      tt.fields.Link,
				UserId:    tt.fields.UserId,
			}
			article.ArticleVote(tt.args.userId)
		})
	}
}

func TestGetArticles(t *testing.T) {

	articles := GetArticles(1, 5, "")
	fmt.Printf("内容：%v", articles)
}

func TestAddRemoveGroups(t *testing.T) {
	type args struct {
		articleId string
		toAdd     []string
		toRemove  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestAddRemoveGroups",
			args: args{
				articleId: "1",
				toRemove:  nil,
				toAdd:     []string{"programing"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddRemoveGroups(tt.args.articleId, tt.args.toAdd, tt.args.toRemove)
		})
	}
}

func TestGetGroupArticles(t *testing.T) {
	type args struct {
		group string
		page  int64
		size  int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetGroupArticles",
			args: args{
				group: "programming",
				page:  1,
				size:  5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles := GetGroupArticles(tt.args.group, tt.args.page, tt.args.size)
			fmt.Printf("article %v", articles)
		})
	}
}
