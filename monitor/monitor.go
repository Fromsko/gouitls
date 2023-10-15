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
	Path         string                // 监控的文件/目录路径
	Delay        time.Duration         // 执行命令的延时
	OnFileChange func(fm *FileMonitor) // 文件变化时的回调函数
	*ProcessInfo
}

type ProcessInfo struct {
	PreviousProcess *os.Process    // 上一次启动的进程
	Event           fsnotify.Event // 监听事件
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

Tab:
	select {
	case event := <-watcher.Events:
		if event.Op&fsnotify.Write == fsnotify.Write {
			fm.Event = event
			// 在每次文件变化时终止上一个进程
			if fm.PreviousProcess != nil {
				_, err := fm.PreviousProcess.Wait()
				// err := fm.PreviousProcess.Kill()
				if err != nil {
					log.Println("Error killing previous process:", err)
				}
			}
			go func() {
				time.Sleep(fm.Delay)
				fm.OnFileChange(fm)
			}()
			goto Tab
		}
	case err := <-watcher.Errors:
		log.Println("Error watching files:", err)
	}
}

// NewFileMonitor 创建一个新的文件监控组件
func NewFileMonitor(path string, onFileChange func(fm *FileMonitor)) *FileMonitor {
	return &FileMonitor{
		Path:         path,
		Delay:        5 * time.Second,
		OnFileChange: onFileChange,
		ProcessInfo: &ProcessInfo{
			Event:           fsnotify.Event{},
			PreviousProcess: nil,
		},
	}
}
