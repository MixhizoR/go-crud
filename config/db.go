package config // Veya projenizin ana paketinin dışındaysa 'database' veya 'internal/database' gibi bir paket adı kullanın

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MixhizoR/go-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB, GORM veritabanı bağlantı nesnesidir. Uygulamanızın her yerinden erişilebilir olmalı.
var DB *gorm.DB

// InitDB veritabanı bağlantısını başlatır ve AutoMigrate işlemini yapar.
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL ortam değişkeni ayarlanmamış. Lütfen docker-compose.yml dosyanızda veya ortamınızda yapılandırın.")
	}

	var err error
	const maxRetries = 10              // Maksimum deneme sayısı
	const retryDelay = 2 * time.Second // Her deneme arası bekleme süresi

	// Veritabanı bağlantısı için bir retry döngüsü
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // GORM'un SQL sorgularını loglaması için
		})
		if err == nil {
			// GORM bağlantısı açıldı, şimdi alttaki SQL bağlantısını ping ile kontrol et
			sqlDB, pingErr := DB.DB()
			if pingErr == nil {
				pingErr = sqlDB.Ping()
				if pingErr == nil {
					fmt.Println("🎉 PostgreSQL veritabanına GORM ile başarıyla bağlanıldı!")
					break // Başarılı bağlantı, döngüden çık
				}
			}
			// If pingErr is not nil, it means the database is not fully ready (e.g., "database system is starting up")
			log.Printf("Veritabanı ping hatası (deneme %d/%d): %v", i+1, maxRetries, pingErr)
		} else {
			// If err is not nil, it means GORM couldn't even open a connection
			log.Printf("Veritabanına bağlanılamadı (deneme %d/%d): %v", i+1, maxRetries, err)
		}

		if i < maxRetries-1 { // Son deneme değilse bekle
			log.Printf("Yeniden deniyor... %v bekleniyor.", retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Fatalf("Maksimum deneme sayısına ulaşıldı. Veritabanına bağlanılamadı: %v", err)
		}
	}

	// Bağlantı havuzu ayarları (GORM bağlantısı başarılıysa devam eder)
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Uyarı: Bağlantı havuzu yapılandırması için alttaki *sql.DB nesnesine erişilemedi: %v", err)
	} else {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	// Modelleri otomatik migrate et
	err = DB.AutoMigrate(&models.User{}) // 'models' paketinizin doğru import edildiğinden emin olun
	if err != nil {
		log.Fatalf("Veritabanı şeması otomatik geçişi başarısız oldu: %v", err)
	}
	fmt.Println("Veritabanı otomatik geçişi başarıyla tamamlandı!")
}

// CloseDB veritabanı bağlantısını kapatır.
// Genellikle uygulamanın kapanışında veya testlerde kullanılır.
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Hata: Veritabanı bağlantısı kapatılırken alttaki *sql.DB nesnesine erişilemedi: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("Veritabanı bağlantısı kapatılırken hata oluştu: %v", err)
		} else {
			fmt.Println("Veritabanı bağlantısı kapatıldı.")
		}
	}
}

// GetDB mevcut GORM DB bağlantısını döndürür.
// Controller'larınızda veya servislerinizde bu fonksiyonu kullanarak DB nesnesine erişin.
func GetDB() *gorm.DB {
	return DB
}
