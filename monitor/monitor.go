package monitor

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/fsnotify.v1"
)

// FileMonitor 是文件监控组件的结构体
type FileMonitor struct {
	Path         string        // 监控的文件/目录路径
	Immediate    bool          // 是否立即执行命令
	Delay        time.Duration // 执行命令的延时
	OnFileChange func()        // 文件变化时的回调函数
}

// Start 启动文件监控服务
func (fm *FileMonitor) Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = filepath.Walk(fm.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Println("Error watching file:", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("File modified:", event.Name)

				if fm.Immediate {
					go fm.OnFileChange()
				} else {
					go func() {
						time.Sleep(fm.Delay)
						fm.OnFileChange()
					}()
				}
			}
		case err := <-watcher.Errors:
			log.Println("Error watching files:", err)
		}
	}
}

// NewFileMonitor 创建一个新的文件监控组件
func NewFileMonitor(path string, immediate bool, delay time.Duration, onFileChange func()) *FileMonitor {
	return &FileMonitor{
		Path:         path,
		Immediate:    immediate,
		Delay:        delay,
		OnFileChange: onFileChange,
	}
}
