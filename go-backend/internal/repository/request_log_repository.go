package repository

import (
	"strings"
	"time"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestLogRepository struct {
	db *gorm.DB
}

func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{db: db}
}

func (r *RequestLogRepository) Create(log *entity.RequestLog) error {
	return r.db.
		Select(
			"traceId",
			"userId",
			"apiKeyId",
			"modelId",
			"requestModel",
			"modelName",
			"requestType",
			"source",
			"promptTokens",
			"completionTokens",
			"totalTokens",
			"cost",
			"duration",
			"status",
			"errorMessage",
			"errorCode",
			"routingStrategy",
			"isFallback",
			"clientIp",
			"userAgent",
			"createTime",
			"updateTime",
		).
		Create(log).Error
}

func (r *RequestLogRepository) ListUserLogs(userID int64, limit int) ([]entity.RequestLog, error) {
	list := make([]entity.RequestLog, 0)
	err := r.db.Model(&entity.RequestLog{}).
		Where("userId = ?", userID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true}).
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *RequestLogRepository) CountUserTokens(userID int64) (int64, error) {
	type tokenSum struct {
		Total int64 `gorm:"column:total"`
	}
	var result tokenSum
	err := r.db.Model(&entity.RequestLog{}).
		Select("COALESCE(SUM(totalTokens), 0) AS total").
		Where("userId = ? AND status = ?", userID, "success").
		Take(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}

func (r *RequestLogRepository) CountUserRequests(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RequestLog{}).
		Where("userId = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RequestLogRepository) CountUserSuccessRequests(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RequestLog{}).
		Where("userId = ? AND status = ?", userID, "success").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RequestLogRepository) SumUserCost(userID int64, startTime, endTime *time.Time) (float64, error) {
	type costSum struct {
		Total float64 `gorm:"column:total"`
	}
	var result costSum
	query := r.db.Model(&entity.RequestLog{}).
		Select("COALESCE(SUM(cost), 0) AS total").
		Where("userId = ? AND status = ?", userID, "success")
	if startTime != nil {
		query = query.Where("createTime >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("createTime <= ?", *endTime)
	}
	if err := query.Take(&result).Error; err != nil {
		return 0, err
	}
	return result.Total, nil
}

func (r *RequestLogRepository) GetByID(id int64) (*entity.RequestLog, error) {
	var record entity.RequestLog
	err := r.db.Model(&entity.RequestLog{}).Where("id = ?", id).Take(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *RequestLogRepository) PageByQuery(req dto.RequestLogQueryRequest) (common.PageResponse[entity.RequestLog], error) {
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	query := r.db.Model(&entity.RequestLog{})
	if req.UserID != nil && req.UserID.Int64() > 0 {
		query = query.Where("userId = ?", req.UserID.Int64())
	}
	if strings.TrimSpace(req.RequestModel) != "" {
		query = query.Where("requestModel LIKE ?", "%"+req.RequestModel+"%")
	}
	if strings.TrimSpace(req.RequestType) != "" {
		query = query.Where("requestType = ?", req.RequestType)
	}
	if strings.TrimSpace(req.Source) != "" {
		query = query.Where("source = ?", req.Source)
	}
	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", req.Status)
	}
	if strings.TrimSpace(req.StartDate) != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err == nil {
			query = query.Where("createTime >= ?", startTime)
		}
	}
	if strings.TrimSpace(req.EndDate) != "" {
		endDate, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err == nil {
			endTime := endDate.Add(24*time.Hour - time.Nanosecond)
			query = query.Where("createTime <= ?", endTime)
		}
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return common.PageResponse[entity.RequestLog]{}, err
	}

	query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: "createTime"}, Desc: true})
	offset := (pageNum - 1) * pageSize
	list := make([]entity.RequestLog, 0)
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return common.PageResponse[entity.RequestLog]{}, err
	}
	return common.BuildPageResponse(list, pageNum, pageSize, total), nil
}

func (r *RequestLogRepository) GetUserDailyStats(userID int64, start, end time.Time) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	current := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	last := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	for !current.After(last) {
		dayStart := current
		dayEnd := current.Add(24*time.Hour - time.Nanosecond)
		logs := make([]entity.RequestLog, 0)
		if err := r.db.Model(&entity.RequestLog{}).
			Where("userId = ?", userID).
			Where("createTime >= ?", dayStart).
			Where("createTime <= ?", dayEnd).
			Find(&logs).Error; err != nil {
			return nil, err
		}

		var totalTokens int64
		var requestCount int64
		var successCount int64
		totalCost := 0.0
		for _, item := range logs {
			requestCount++
			if item.Status == "success" {
				successCount++
				totalTokens += int64(item.TotalTokens)
				totalCost += item.Cost
			}
		}
		result = append(result, map[string]any{
			"date":         dayStart.Format("2006-01-02"),
			"totalTokens":  totalTokens,
			"requestCount": requestCount,
			"successCount": successCount,
			"totalCost":    totalCost,
		})
		current = current.AddDate(0, 0, 1)
	}
	return result, nil
}

type ModelStatsRow struct {
	ModelName    string  `gorm:"column:modelName"`
	AvgLatency   int     `gorm:"column:avgLatency"`
	SuccessRate  float64 `gorm:"column:successRate"`
	TotalRequest int64   `gorm:"column:totalRequest"`
}

func (r *RequestLogRepository) QueryModelStatsSince(startTime time.Time) ([]ModelStatsRow, error) {
	rows := make([]ModelStatsRow, 0)
	err := r.db.Model(&entity.RequestLog{}).
		Select(
			"modelName",
			"CAST(AVG(CASE WHEN status = 'success' THEN duration ELSE NULL END) AS SIGNED) AS avgLatency",
			"IFNULL(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) * 100.0 / NULLIF(COUNT(*), 0), 100) AS successRate",
			"COUNT(*) AS totalRequest",
		).
		Where("createTime >= ?", startTime).
		Where("modelName IS NOT NULL AND modelName <> ''").
		Group("modelName").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
