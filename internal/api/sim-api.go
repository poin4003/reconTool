package api

import (
	"archive/zip"
	"fmt"
	"net/http"
	"reconTool/global"
	"reconTool/internal/dto"
	"reconTool/internal/model"
	"reconTool/internal/service"
	"strconv"
	"time"

	"github.com/szyhf/go-excel"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func InitSimRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/sim")
	{
		r.POST("/import", ImportSimExcel)
		r.GET("/export", ExportSimExcel)
		r.POST("/sync-sim/list", SyncSimsHandler)
		r.POST("/sync-sim", SyncSimApi)
		r.POST("", CreateSim)
		r.GET("", GetSimList)
		r.DELETE("", DeleteSims)
	}
}

func CreateSim(c *gin.Context) {
	var sim model.SimModel

	if err := c.ShouldBindJSON(&sim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	if sim.Isdn == "" || sim.Serial == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Isdn và Serial không được để trống"})
		return
	}

	err := service.CheckAndInsertSim(sim)
	if err != nil {
		global.Logger.Warn("Tạo SIM thất bại", zap.String("isdn", sim.Isdn), zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	global.Logger.Info("Tạo SIM thành công", zap.String("isdn", sim.Isdn))
	c.JSON(http.StatusCreated, gin.H{"message": "Tạo SIM thành công", "data": sim})
}

func GetSimList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var filter dto.SimFilter
	c.ShouldBindQuery(&filter)

	result, err := service.GetSims(page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func ImportSimExcel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn file"})
		return
	}
	defer file.Close()

	zipReader, err := zip.NewReader(file, header.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không hợp lệ"})
		return
	}

	conn := excel.NewConnecter()
	if err := conn.OpenReader(zipReader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	var sims []model.SimModel
	rd, _ := conn.NewReader("Sheet1")
	rd.ReadAll(&sims)

	for _, sim := range sims {
		global.ImportSimJobChan <- sim
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xử lý hoàn tất"})
}

func ExportSimExcel(c *gin.Context) {
	var filter dto.SimFilter
	c.ShouldBindQuery(&filter)

	sims, err := service.GetSimsForExport(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy dữ liệu xuất file"})
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	sheetName := "Danh Sách SIM"
	f.SetSheetName("Sheet1", sheetName)

	headers := map[string]string{"A1": "Isdn", "B1": "Serial", "C1": "Status", "D1": "Note", "E1": "CreatedAt"}

	for cell, val := range headers {
		f.SetCellValue(sheetName, cell, val)
	}

	for i, sim := range sims {
		noteVal := ""
		if sim.Note != nil {
			noteVal = *sim.Note
		}
		rowIdx := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), sim.Isdn)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), sim.Serial)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), sim.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), noteVal)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), sim.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=sim_export_%s.xlsx", time.Now().Format("20060102")))
	f.Write(c.Writer)
}

func DeleteSims(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Danh sách ID không hợp lệ hoặc trống"})
		return
	}

	if err := service.DeleteSims(req.IDs); err != nil {
		global.Logger.Error("Lỗi xóa SIM", zap.Strings("ids", req.IDs), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xử lý xóa thất bại"})
		return
	}

	global.Logger.Info("Đã xóa vĩnh viễn danh sách SIM", zap.Strings("ids", req.IDs))
	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa thành công", "deleted_count": len(req.IDs)})
}

func SyncSimApi(c *gin.Context) {
	var sims []model.SimModel
	if err := global.DB.Find(&sims).Error; err != nil {
		global.Logger.Error("Không thể lấy danh sách SIM", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn DB"})
		return
	}

	total := len(sims)
	global.Logger.Info("Bắt đầu đẩy SIM vào hàng đợi đối soát", zap.Int("total", total))

	go func() {
		for _, sim := range sims {
			global.SyncSimJobChan <- sim
		}
		global.Logger.Info("Đã đẩy xong toàn bộ SIM vào hàng đợi đối soát", zap.Int("total", total))
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã bắt đầu tiến trình đối soát ngầm",
		"total":   total,
		"status":  "processing",
	})
}

func SyncSimsHandler(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Danh sách ID trống"})
		return
	}

	var sims []model.SimModel
	if err := global.DB.Where("id IN ?", req.IDs).Find(&sims).Error; err != nil {
		global.Logger.Error("Lỗi truy vấn SIM", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn cơ sở dữ liệu"})
		return
	}

	if len(sims) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy SIM nào trong danh sách ID cung cấp"})
		return
	}

	go func() {
		count := 0
		for _, sim := range sims {
			global.SyncSimJobChan <- sim
			count++
		}
		global.Logger.Info("Đã đẩy thủ công SIM vào hàng đợi", zap.Int("count", count))
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã tiếp nhận yêu cầu sync cho %d SIM", len(sims)),
		"count":   len(sims),
		"status":  "processing",
	})
}
