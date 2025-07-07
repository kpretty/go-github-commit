package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const (
	githubCommitUrl = "https://github.com/users/%s/contributions"
)

type githubCommitInfo struct {
	Data  string `json:"data"`
	Level string `json:"level"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Referer")

	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "缺少user参数", http.StatusBadRequest)
		return
	}

	if r.URL.Path == "/api" {
		githubCommit, err := getGithubCommit(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(githubCommit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

// ------------------------------------------------------------------------------
// 通过爬虫获取 github commit 信息
// 需要 github username 和 year
// ------------------------------------------------------------------------------
// userName github 用户名
func getGithubCommit(userName string) ([]*githubCommitInfo, error) {
	// 拼接请求参数
	url := fmt.Sprintf(githubCommitUrl, userName)
	log.Printf("GET %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	// 使用 goquery 解析 html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make([]*githubCommitInfo, 0) // 申请可以存储一年的空间，尽可能避免扩容逻辑
	doc.Find("table.js-calendar-graph-table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				date, exists := s.Attr("data-date")
				if !exists {
					return
				}
				level, _ := s.Attr("data-level")
				result = append(result, &githubCommitInfo{Data: date, Level: level})
			})
		})
	})

	return result, nil
}
