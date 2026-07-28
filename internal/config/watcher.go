package config

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Watch(path string, onChange func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(path); err != nil {
		return err
	}

	var debounce *time.Timer

	go func() {
		defer func() {
			_ = watcher.Close()
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if !event.Has(fsnotify.Write) &&
					!event.Has(fsnotify.Create) &&
					!event.Has(fsnotify.Rename) {
					continue
				}

				// Reset the debounce timer if another event arrives.
				if debounce != nil {
					debounce.Stop()
					debounce = nil
				}

				debounce = time.AfterFunc(200*time.Millisecond, func() {
					log.Println("Config file changed, reloading...")
					onChange()
				})

			case err := <-watcher.Errors:
				log.Println("Error watching config file:", err)
			}
		}
	}()

	return nil
}
