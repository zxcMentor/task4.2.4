package main

import (
	"testing"

	"github.com/redis/go-redis"
)

func TestNewGeoserviceRepositoryProxy(t *testing.T) {
	mockGeoserviceRepo := &GeoserviceRepositoryImpl{}

	mockRedisClient := &redis.Client{}

	proxy := NewGeoserviceRepositoryProxy(mockGeoserviceRepo, mockRedisClient)

	if proxy.geoserviceRepository != mockGeoserviceRepo {
		t.Errorf("Expected geoserviceRepository to be set to mockGeoserviceRepo, but got %v", proxy.geoserviceRepository)
	}

	if proxy.cache != mockRedisClient {
		t.Errorf("Expected cache to be set to mockRedisClient, but got %v", proxy.cache)
	}
}

func TestNewReverseProxy(t *testing.T) {
	mockGeoserviceRepo := &GeoserviceRepositoryImpl{}

	port := "8080"
	proxy := NewReverseProxy(mockGeoserviceRepo, port)

	if proxy.geoserviceRepository != mockGeoserviceRepo {
		t.Errorf("Expected geoserviceRepository to be set to mockGeoserviceRepo, but got %v", proxy.geoserviceRepository)
	}

	if proxy.port != port {
		t.Errorf("Expected port to be set to %s, but got %s", port, proxy.port)
	}
}
