package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateAndFillTemp(t *testing.T, pattern string, content string) string {
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	t.Cleanup(func() { os.Remove(tempFile.Name()) })
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Ошибка при записи в временный файл: %v", err)
	}
	tempFile.Close()
	return tempFile.Name()
}

// тест при правильной конфигурации
func TestMustLoad_ValidConfigFile(t *testing.T) {
	const validConfigFileName = "config*.yml"
	const validConfig = `env: "dev"
http_server:
  address: "1.1.1.1:8080"
  timeout: "10s"
  idle_timeout: "120s"
  rate limit:
    limit: 5
    interval: "2s" 
    msg: "Тест."
`
	name := CreateAndFillTemp(t, validConfigFileName, validConfig)
	cfg := MustLoad(name)
	assert.Equal(t, "dev", cfg.Env)
	assert.Equal(t, "1.1.1.1:8080", cfg.Address)
	assert.Equal(t, 10*time.Second, cfg.Timeout)
	assert.Equal(t, 120*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 5, cfg.Limit)
	assert.Equal(t, 2*time.Second, cfg.Interval)
	assert.Equal(t, "Тест.", cfg.Msg)
}

// тест при отсуствии файла кофигурации
func TestMustLoad_MissingConfigFile(t *testing.T) {
	assert.Panics(t, func() { _ = MustLoad("non_existing_file.yml") })
}

// файл с некорректной конфигурацией
func TestMustLoad_InvalidConfigFile(t *testing.T) {
	const invalidConfigFileName = "invalid_config*.yml"
	const invalidConfig = `invalid_yaml`
	name := CreateAndFillTemp(t, invalidConfigFileName, invalidConfig)
	assert.Panics(t, func() { _ = MustLoad(name) })
}

// проверка того, что функция выставляет стандартные значения в случае отсуствия их в файле конфигурации
func TestMustLoad_ValidConfigFile_DefaultValues(t *testing.T) {
	const partialConfigFileName = "partial_config*.yml"
	const partialConfig = `env: "prod"
http_server:
  address: "1.1.1.1:8080"
`
	name := CreateAndFillTemp(t, partialConfigFileName, partialConfig)
	cfg := MustLoad(name)
	assert.Equal(t, "prod", cfg.Env)
	assert.Equal(t, "1.1.1.1:8080", cfg.Address)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, 60*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 100, cfg.Limit)
	assert.Equal(t, 60*time.Second, cfg.Interval)
	assert.Equal(t, "Слишком много запросов.", cfg.Msg)
}

func TestMustLoad_InvalidConfigFile_AbsenceOfEnv(t *testing.T) {
	const invalidConfigFileName = "invalid_config*.yml"
	const invalidConfig = `http_server:
  address: "1.1.1.1:8080"
  timeout: "10s"
  idle_timeout: "120s"
`
	name := CreateAndFillTemp(t, invalidConfigFileName, invalidConfig)
	assert.Panics(t, func() { _ = MustLoad(name) })
}

func TestMustLoad_InvalidConfigFile_AbsenceOfAddress(t *testing.T) {
	const invalidConfigFileName = "invalid_config*.yml"
	const invalidConfig = `env: "dev"
http_server:
  timeout: "10s"
  idle_timeout: "120s"
`
	name := CreateAndFillTemp(t, invalidConfigFileName, invalidConfig)
	assert.Panics(t, func() { _ = MustLoad(name) })
}
