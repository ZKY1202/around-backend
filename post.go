package main

import (
	"mime/multipart"
	"reflect"

	"github.com/olivere/elastic/v7"
)

const (
	POST_INDEX = "post"
)

// entity -> 相当于 java 里的 POJO （plain old java object）
// public field (大写开头字母 public， 小写开头字母 private)
type Post struct {
	Id      string `json:"id"` // `` （反引号）-> 表示raw string
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
	query := elastic.NewMatchQuery("user", user)
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
	var ptype Post
	var posts []Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(Post)
		posts = append(posts, p)
	}
	return posts
}

func savePost(post *Post, file multipart.File) error {
	mediaLink, err := saveToGCS(file, post.Id)
	if err != nil {
		return err
	}

	post.Url = mediaLink
	return saveToES(post, POST_INDEX, post.Id)
}
