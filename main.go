package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis"
)

type GeoserviceRepository interface {
	GetData() string
}
type GeoserviceRepositoryImpl struct {
	geoserviceRepository GeoserviceRepository
	cache                *redis.Client
}

func (r *GeoserviceRepositoryImpl) GetData() string {
	return "geoservice data"
}

type GeoserviceRepositoryProxy struct {
	geoserviceRepository GeoserviceRepository
	cache                *redis.Client
}

type ReverseProxy struct {
	geoserviceRepository GeoserviceRepository
	port                 string
}

func NewReverseProxy(geoserviceRepository GeoserviceRepository, port string) *ReverseProxy {
	return &ReverseProxy{
		geoserviceRepository: geoserviceRepository,
		port:                 port,
	}
}

func NewGeoserviceRepositoryProxy(geoserviceRepository GeoserviceRepository, cache *redis.Client) *GeoserviceRepositoryProxy {
	return &GeoserviceRepositoryProxy{
		geoserviceRepository: geoserviceRepository,
		cache:                cache,
	}
}

func (r *GeoserviceRepositoryProxy) GetData() string {
	data, err := r.cache.Get("geoservice_data").Result()
	if err == nil {
		return data
	}

	originalData := r.geoserviceRepository.GetData()

	err = r.cache.Set("geoservice_data", originalData, 5*time.Minute).Err()
	if err != nil {
		log.Println("Error caching geoservice data:", err)
	}

	return originalData
}

func main() {
	geoserviceRepo := &GeoserviceRepositoryImpl{}
	geoserviceProxy := NewGeoserviceRepositoryProxy(geoserviceRepo, redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))

	//os.Setenv("HOST", "127.0.0.1:1313 hugo")
	router := chi.NewRouter()

	proxy := NewReverseProxy(geoserviceProxy, "1313")
	//proxy := NewReverseProxy("hugo", "1313")
	router.Use(proxy.ReverseProxy)

	router.Get("/api", handlerRoute)

	http.ListenAndServe(":8080", router)
}

func handlerRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var text = fmt.Sprintf("<!DOCTYPE html><html><head><title>Webserver</title></head><body>Hello API</body></html>")

	w.Write([]byte(text))

}

func (rp *ReverseProxy) ReverseProxy(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api" {
			fmt.Println("proxy start")
			targetURL, _ := url.Parse("http://hugo:1313")
			proxy := httputil.NewSingleHostReverseProxy(targetURL)
			proxy.ServeHTTP(w, r)

			return

		}
		handler.ServeHTTP(w, r)
	})
}

const content = `---
menu:
    before:
        name: tasks
        weight: 5
title: Обновление данных в реальном времени
---

# Задача: Обновление данных в реальном времени

Напишите воркер, который будет обновлять данные в реальном времени, на текущей странице.
Текст данной задачи менять нельзя, только время и счетчик.

Файл данной страницы: /app/static/tasks/_index.md

Должен меняться счетчик и время:

Текущее время:%s

Счетчик: %d



## Критерии приемки:
- [ ] Воркер должен обновлять данные каждые 5 секунд
- [ ] Счетчик должен увеличиваться на 1 каждые 5 секунд
- [ ] Время должно обновляться каждые 5 секунд
`

func WorkerTest3() {
	t := time.NewTicker(5 * time.Second)
	var b int
	path := "/app/static/tasks/_index.md"
	for {
		select {
		case <-t.C:
			{
				err := os.WriteFile(path, []byte(fmt.Sprintf(content, time.Now().Format("2006-01-02 15:04:05"), b)), 0644)
				if err != nil {
					log.Println(err)
				}
				b++
			}
		}
	}
}
