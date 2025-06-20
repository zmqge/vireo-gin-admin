package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NoticesService Notices服务接口
type NoticesService interface {
	GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) // 使用uint类型
	ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error)
	CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	DeleteNotices(ctx *gin.Context, id uint) error // 使用uint类型
	PageNotices(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	GetNoticesForm(ctx *gin.Context, id uint) (*models.NoticesModel, error) // 使用uint类型
	RevokeNotice(ctx *gin.Context, id uint) error
	PublishNotice(ctx *gin.Context, id uint) error
	PublishNoticeWithReceivers(ctx *gin.Context, id uint) error
	MyPageNotices(ctx *gin.Context, userID uint, keywords string, isRead uint, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	GetMyNoticesByID(ctx *gin.Context, userID uint, id uint) (*models.NoticesModel, error)
	MarkAllAsRead(ctx *gin.Context, userID uint) error
}

// NoticesServiceImpl Notices服务实现
type NoticesServiceImpl struct {
	repo         repositories.NoticesRepository
	receiverRepo repositories.NoticeReceiverRepository
	userRepo     repositories.UserRepository
}

// NewNoticesService 创建Notices服务
func NewNoticesService(
	repo repositories.NoticesRepository,
	userRepo repositories.UserRepository,
	receiverRepo repositories.NoticeReceiverRepository,
) NoticesService {
	return &NoticesServiceImpl{
		repo:         repo,
		userRepo:     userRepo,
		receiverRepo: receiverRepo,
	}
}

// GetNoticesByID 根据ID获取Notices
func (s *NoticesServiceImpl) GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	entity, err := s.repo.GetNoticesByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("Notices not found")
	}
	return entity, nil
}

// GetMyNoticesByID 根据用户id和NoticesID获取用户的Notices
func (s *NoticesServiceImpl) GetMyNoticesByID(ctx *gin.Context, userID uint, id uint) (*models.NoticesModel, error) {
	entity, err := s.repo.GetMyNoticesByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("Notices not found")
	}
	// 2. 标记为已读（无论当前是否已读都更新，确保最后查看时间准确）
	if err := s.repo.MarkNoticeAsRead(ctx, userID, id); err != nil {
		// 可以根据业务需求决定是否返回错误
		log.Printf("警告：标记通知为已读失败（通知ID：%d，用户ID：%d），错误：%v", id, userID, err)
	}
	return entity, nil
}

// ListNoticess 获取Notices列表
func (s *NoticesServiceImpl) ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error) {
	return s.repo.ListNoticess(ctx)
}

// CreateNotices 创建Notices
func (s *NoticesServiceImpl) CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	return s.repo.CreateNotices(ctx, entity)
}

// UpdateNotices 更新Notices
func (s *NoticesServiceImpl) UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	if entity.Title == "" {
		return errors.New("name is required")
	}
	return s.repo.UpdateNotices(ctx, entity)
}

// DeleteNotices 删除Notices
func (s *NoticesServiceImpl) DeleteNotices(ctx *gin.Context, id uint) error {
	return s.repo.DeleteNotices(ctx, id)
}

// PageNotices 分页获取Notices列表
func (s *NoticesServiceImpl) PageNotices(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
	return s.repo.PageNotices(ctx, keywords, publishStatus, pageNum, pageSize)
}

// GetNoticesForm 表单
func (s *NoticesServiceImpl) GetNoticesForm(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	return s.repo.GetNoticesByID(ctx, id)
}

// RevokeNotice 撤回通知
func (s *NoticesServiceImpl) RevokeNotice(ctx *gin.Context, id uint) error {
	return s.repo.RevokeNotice(ctx, id)
}

// PublishNotice 发布通知
func (s *NoticesServiceImpl) PublishNotice(ctx *gin.Context, id uint) error {
	return s.repo.PublishNotice(ctx, id)
}

func (s *NoticesServiceImpl) PublishNoticeWithReceivers(ctx *gin.Context, id uint) error {
	// 1. 参数校验
	if ctx == nil {
		return errors.New("gin context cannot be nil")
	}

	// 获取通知详情（包括接收者关系）
	entity, err := s.repo.GetNoticeWithReceivers(ctx, id)
	if err != nil {
		return fmt.Errorf("获取通知详情失败: %w", err)
	}
	if entity == nil {
		return errors.New("通知不存在")
	}

	// 2. 校验通知状态
	if entity.Status == models.NoticeStatusPublished {
		return errors.New("通知已发布，无需重复操作")
	}

	// 3. 校验基本字段
	if entity.Title == "" {
		return errors.New("标题不能为空")
	}
	if entity.TargetType < 1 || entity.TargetType > 4 {
		return errors.New("目标类型无效")
	}
	if entity.TargetType == 4 && len(entity.TargetIDs) == 0 {
		return errors.New("指定用户发布时，目标用户ID不能为空")
	}

	// 4. 开启事务
	tx := s.repo.BeginTx(ctx)
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] 事务回滚: %v", r)
		}
	}()

	// 5. 更新通知状态为已发布
	if err := s.repo.UpdateNoticeStatusTx(tx, id, models.NoticeStatusPublished); err != nil {
		return fmt.Errorf("更新通知状态失败: %w", err)
	}

	// 6. 处理接收者（支持重新发布）
	if err := s.processNoticeReceiversForPublish(ctx, tx, entity); err != nil {
		return fmt.Errorf("处理接收者失败: %w", err)
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// 专门处理发布时的接收者逻辑
func (s *NoticesServiceImpl) processNoticeReceiversForPublish(ctx *gin.Context, tx *gorm.DB, notice *models.NoticesModel) error {
	// 1. 获取目标用户ID
	targetUserIDs, err := s.getTargetUserIDs(ctx, notice)
	if err != nil {
		return err
	}

	// 2. 批量处理接收者关系
	receivers := make([]models.NoticeReceiver, 0, len(targetUserIDs))
	for _, userID := range targetUserIDs {
		receivers = append(receivers, models.NoticeReceiver{
			NoticeID:  notice.ID,
			UserID:    userID,
			IsRead:    0,
			CreatedAt: time.Now(),
		})
	}

	// 3. 使用事务批量插入（存在则更新）
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "notice_id"},
			{Name: "user_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"is_read":    0,
			"read_time":  nil,
			"updated_at": time.Now(),
		}),
	}).CreateInBatches(receivers, 500).Error; err != nil {
		return fmt.Errorf("批量处理接收者失败: %w", err)
	}

	return nil
}

// 获取目标用户ID（根据通知类型）
func (s *NoticesServiceImpl) getTargetUserIDs(ctx *gin.Context, notice *models.NoticesModel) ([]uint, error) {
	switch notice.TargetType {
	case 1: // 全部用户
		return s.getAllActiveUserIDs(ctx)
	case 2: // 指定用户
		return notice.TargetIDs, nil
	case 3: // 指定部门
		if len(notice.TargetIDs) == 0 {
			return nil, errors.New("部门ID不能为空")
		}
		return s.getUserIDsByDepartments(ctx, notice.TargetIDs)
	case 4: // 指定角色
		if len(notice.TargetIDs) == 0 {
			return nil, errors.New("角色ID不能为空")
		}
		return s.getUserIDsByRoles(ctx, notice.TargetIDs)
	default:
		return nil, errors.New("不支持的目标类型")
	}
}

// 获取所有活跃用户ID
func (s *NoticesServiceImpl) getAllActiveUserIDs(ctx *gin.Context) ([]uint, error) {
	// 验证数据库连接是否存在
	if s.userRepo.GetDB() == nil {
		return nil, fmt.Errorf("userRepo数据库连接为空")
	}

	var userIDs []uint

	// 使用gin上下文的请求ID作为日志标识
	requestID, exists := ctx.Get("request_id")
	if !exists {
		requestID = "unknown"
	}

	// 计算30天前的时间点
	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)

	// 构建查询
	query := s.userRepo.GetDB().WithContext(ctx.Request.Context()).
		Model(&models.User{}).
		Where("last_login_time > ?", thirtyDaysAgo).
		Pluck("id", &userIDs)

	// 执行查询并检查错误
	if query.Error != nil {
		log.Printf("[ERROR] [RequestID: %v] 查询活跃用户ID失败: %v", requestID, query.Error)
		return nil, query.Error
	}

	// 记录查询结果
	log.Printf("[DEBUG] [RequestID: %v] 查询到的活跃用户ID数量: %d", requestID, len(userIDs))

	return userIDs, nil
}

// 根据部门获取用户ID
func (s *NoticesServiceImpl) getUserIDsByDepartments(ctx *gin.Context, deptIDs []uint) ([]uint, error) {
	var userIDs []uint
	err := s.userRepo.GetDB().WithContext(ctx).
		Model(&models.UserDept{}).
		Where("dept_id IN ?", deptIDs).
		Pluck("DISTINCT user_id", &userIDs).Error
	return userIDs, err
}

// 根据角色获取用户ID
func (s *NoticesServiceImpl) getUserIDsByRoles(ctx *gin.Context, roleIDs []uint) ([]uint, error) {
	var userIDs []uint
	err := s.userRepo.GetDB().WithContext(ctx).
		Model(&models.UserRole{}).
		Where("role_id IN ?", roleIDs).
		Pluck("DISTINCT user_id", &userIDs).Error
	return userIDs, err
}

// 分批插入接收者记录
func (s *NoticesServiceImpl) batchInsertReceivers(tx *gorm.DB, noticeID uint, userIDs []uint) error {
	const batchSize = 500
	total := len(userIDs)

	if total == 0 {
		log.Printf("[WARN] 没有可插入的接收者用户ID\n")
		return nil
	}

	// 去重处理
	uniqueUserIDs := make(map[uint]bool)
	for _, id := range userIDs {
		uniqueUserIDs[id] = true
	}

	// 转换为切片
	uniqueUsers := make([]uint, 0, len(uniqueUserIDs))
	for id := range uniqueUserIDs {
		uniqueUsers = append(uniqueUsers, id)
	}

	log.Printf("[DEBUG] 开始批量插入接收者，去重后总数: %d，将分 %d 批处理\n", len(uniqueUsers), (len(uniqueUsers)+batchSize-1)/batchSize)

	for i := 0; i < len(uniqueUsers); i += batchSize {
		end := i + batchSize
		if end > len(uniqueUsers) {
			end = len(uniqueUsers)
		}

		batch := uniqueUsers[i:end]
		receivers := make([]models.NoticeReceiver, len(batch))

		for j, userID := range batch {
			receivers[j] = models.NoticeReceiver{
				NoticeID: noticeID,
				UserID:   userID,
				IsRead:   0,
			}
		}

		// 使用ON CONFLICT DO NOTHING处理重复插入
		// 使用正确的 clause.OnConflict
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "notice_id"},
				{Name: "user_id"},
			},
			DoNothing: true,
		}).CreateInBatches(receivers, batchSize).Error; err != nil {
			return fmt.Errorf("batch insert receivers failed at batch %d: %w", i/batchSize, err)
		}
		log.Printf("[DEBUG] 成功插入批次 %d，用户数: %d\n", i/batchSize, len(batch))
	}

	log.Printf("[INFO] 所有接收者插入完成，去重后总数: %d\n", len(uniqueUsers))
	return nil
}

// MyPageNotices 获取我的公告列表
func (s *NoticesServiceImpl) MyPageNotices(ctx *gin.Context, userID uint, keywords string, isRead uint, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
	return s.repo.MyPageNotices(ctx, userID, keywords, isRead, pageNum, pageSize)
}

// MarkAllAsRead 标记全部为已读
func (s *NoticesServiceImpl) MarkAllAsRead(ctx *gin.Context, userID uint) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
