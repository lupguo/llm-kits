package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/webdav"
)

// CustomLockSystem 自定义锁系统（支持持久化锁令牌）
type CustomLockSystem struct {
	webdav.LockSystem
	tokens map[string]bool // 存储活跃锁令牌
}

func NewCustomLockSystem() *CustomLockSystem {
	return &CustomLockSystem{
		LockSystem: webdav.NewMemLS(),
		tokens:     make(map[string]bool),
	}
}

func (cls *CustomLockSystem) GenerateToken() string {
	token := "urn:uuid:" + strings.ReplaceAll(uuid.New().String(), "-", "")
	cls.tokens[token] = true
	return token
}

func (cls *CustomLockSystem) IsValidToken(token string) bool {
	return cls.tokens[token]
}

// VersionedFileSystem 支持版本控制的文件系统
type VersionedFileSystem struct {
	webdav.Dir
	versionsDir string // 版本存储目录
}

func (vfs *VersionedFileSystem) MkdirAll(versionsDir string) error {
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		return err
	}
	vfs.versionsDir = versionsDir
	return nil
}

func (vfs *VersionedFileSystem) SaveVersion(path string) error {
	src := filepath.Join(string(vfs.Dir), path)
	dst := filepath.Join(vfs.versionsDir, path+"."+time.Now().Format("20060102-150405"))
	return os.Link(src, dst) // 硬链接节省空间
}

// 主服务
func main() {
	// 初始化目录
	dataDir := "/tmp/webdata/data"
	versionsDir := "/tmp/webdata/versions"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		log.Fatal(err)
	}

	// 初始化锁系统和文件系统
	lockSystem := NewCustomLockSystem()
	vfs := &VersionedFileSystem{Dir: webdav.Dir(dataDir)}
	if err := vfs.MkdirAll(versionsDir); err != nil {
		log.Fatal(err)
	}

	// 创建WebDAV处理器
	handler := &webdav.Handler{
		FileSystem: vfs,
		LockSystem: lockSystem,
		Logger: func(r *http.Request, err error) {
			op := r.Method
			if op == "LOCK" || op == "UNLOCK" {
				log.Printf("WebDAV %s %s | Token: %v", op, r.URL.Path, r.Header.Get("Lock-Token"))
			}
			if err != nil {
				log.Printf("Error: %v", err)
			}
		},
	}

	// 添加中间件（认证+版本控制）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. 基础认证
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "123456" {
			w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. 在PUT操作前保存旧版本
		if r.Method == "PUT" {
			if _, err := os.Stat(filepath.Join(dataDir, r.URL.Path)); err == nil {
				if err := vfs.SaveVersion(r.URL.Path); err != nil {
					log.Printf("Version save failed: %v", err)
				}
			}
		}

		// 3. 处理WebDAV请求
		handler.ServeHTTP(w, r)
	})

	// 启动服务器
	log.Println("WebDAV server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
