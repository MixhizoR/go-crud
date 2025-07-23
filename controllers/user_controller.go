// controllers/user_controller.go (veya sadece controllers.go)
package controllers

import (
	"net/http"
	"strconv" // string'den int'e dönüşüm için

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/MixhizoR/go-crud/config" // db.go'nuzun bulunduğu paketi import edin
	"github.com/MixhizoR/go-crud/models" // Model tanımlarınızın bulunduğu paketi import edin
)

// users slice'ı kaldırıldı, artık veritabanı kullanılacak.

// GetUsers tüm kullanıcıları veritabanından getirir.
func GetUsers(c *gin.Context) {
	db := config.GetDB() // GORM DB bağlantısını al
	var users []models.User
	result := db.Find(&users) // Tüm kullanıcıları bul

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcılar getirilirken hata oluştu", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser belirli bir kullanıcıyı ID'sine göre getirir.
func GetUser(c *gin.Context) {
	db := config.GetDB() // GORM DB bağlantısını al
	idParam := c.Param("id")

	// GORM'un ID'si genellikle int64'tür (gorm.Model'den gelir)
	id, err := strconv.ParseUint(idParam, 10, 64) // string ID'yi uint64'e dönüştür
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si", "details": err.Error()})
		return
	}

	var user models.User
	// First: Veritabanından ID'ye göre ilk kaydı getirir.
	// If not found, returns gorm.ErrRecordNotFound error.
	result := db.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı getirilirken hata oluştu", "details": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	db := config.GetDB() // GORM DB bağlantısını al
	var user models.User

	// JSON body'yi User struct'ına bağla
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kullanıcıyı veritabanına kaydet
	result := db.Create(&user)
	if result.Error != nil {
		// Örneğin, Email benzersiz değilse burada bir hata dönebilir.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı oluşturulurken hata oluştu", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user) // GORM, oluşturulan kullanıcının ID'sini `user` nesnesine dolduracaktır.
}

func UpdateUser(c *gin.Context) {
	db := config.GetDB() // GORM DB bağlantısını al
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64) // string ID'yi uint64'e dönüştür
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si", "details": err.Error()})
		return
	}

	var existingUser models.User
	// Güncellenecek kullanıcıyı veritabanından getir
	result := db.First(&existingUser, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı getirilirken hata oluştu", "details": result.Error.Error()})
		}
		return
	}

	var updatedUserInput models.User
	// Güncelleme verilerini JSON body'den al
	if err := c.ShouldBindJSON(&updatedUserInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sadece güncellenecek alanları mevcut kullanıcıya kopyala
	// ID'yi manuel olarak atamaya gerek yok, GORM kendisi halleder.
	existingUser.Name = updatedUserInput.Name
	existingUser.Email = updatedUserInput.Email
	// Diğer alanları da buraya ekleyebilirsiniz.

	// Değişiklikleri veritabanına kaydet
	result = db.Save(&existingUser) // Save() metodunu kullan
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı güncellenirken hata oluştu", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, existingUser) // Güncellenmiş kullanıcıyı döndür
}

// DeleteUser belirli bir kullanıcıyı ID'sine göre siler (GORM ile soft delete).
func DeleteUser(c *gin.Context) {
	db := config.GetDB() // GORM DB bağlantısını al
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64) // string ID'yi uint64'e dönüştür
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID'si", "details": err.Error()})
		return
	}

	var user models.User
	// İlk olarak kullanıcıyı bulup var olduğundan emin olalım
	result := db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı getirilirken hata oluştu", "details": result.Error.Error()})
		}
		return
	}

	// Kullanıcıyı sil (gorm.Model olduğu için soft delete yapar)
	// Yani, kaydı gerçekten silmez, sadece 'deleted_at' alanını doldurur.
	result = db.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı silinirken hata oluştu", "details": result.Error.Error()})
		return
	}

	// 204 No Content döndür, çünkü silme başarılı oldu ve döndürülecek bir içerik yok.
	c.JSON(http.StatusNoContent, nil)
}
