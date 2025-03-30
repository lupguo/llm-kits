package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func main() {
	// 设置WebDAV根目录
	dir := "./webdav_data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	// 创建WebDAV处理器
	handler := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(), // 内存锁系统
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WebDAV %s %s: %v", r.Method, r.URL.Path, err)
			}
		},
	}

	// 添加基础认证中间件
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "123456" {
			w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})

	// 启动服务器
	log.Println("WebDAV server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", authHandler))
}
