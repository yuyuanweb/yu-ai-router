package com.yupi.airouter.service.impl;

import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.mapper.RechargeRecordMapper;
import com.yupi.airouter.model.entity.RechargeRecord;
import com.yupi.airouter.service.BalanceService;
import com.yupi.airouter.service.RechargeService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 充值记录服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Service
@Slf4j
public class RechargeServiceImpl extends ServiceImpl<RechargeRecordMapper, RechargeRecord> implements RechargeService {

    @Resource
    private BalanceService balanceService;

    @Override
    public RechargeRecord createRechargeRecord(Long userId, BigDecimal amount, String paymentMethod) {
        RechargeRecord record = RechargeRecord.builder()
                .userId(userId)
                .amount(amount)
                .paymentMethod(paymentMethod)
                .status("pending")
                .description("账户充值")
                .createTime(LocalDateTime.now())
                .updateTime(LocalDateTime.now())
                .build();

        save(record);
        log.info("创建充值记录：用户 {}, 金额 ¥{}", userId, amount);
        return record;
    }

    @Override
    public boolean updateRechargeStatus(Long recordId, String status, String paymentId) {
        RechargeRecord record = getById(recordId);
        if (record == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "充值记录不存在");
        }

        record.setStatus(status);
        record.setPaymentId(paymentId);
        record.setUpdateTime(LocalDateTime.now());

        return updateById(record);
    }

    @Override
    public RechargeRecord getByPaymentId(String paymentId) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .where("paymentId = ?", paymentId);

        return getOne(queryWrapper);
    }

    @Override
    public Page<RechargeRecord> listUserRechargeRecords(Long userId, int pageNum, int pageSize) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .where("userId = ?", userId)
                .orderBy("createTime", false);

        return page(Page.of(pageNum, pageSize), queryWrapper);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public boolean completeRecharge(Long recordId, String paymentId) {
        RechargeRecord record = getById(recordId);
        if (record == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "充值记录不存在");
        }

        // 防止重复处理
        if ("success".equals(record.getStatus())) {
            log.warn("充值记录 {} 已经处理过，跳过", recordId);
            return true;
        }

        // 更新充值记录状态
        record.setStatus("success");
        record.setPaymentId(paymentId);
        record.setUpdateTime(LocalDateTime.now());
        updateById(record);

        // 增加用户余额
        balanceService.addBalance(record.getUserId(), record.getAmount(), "充值：" + paymentId);

        log.info("充值完成：用户 {}, 金额 ¥{}, 支付ID {}", record.getUserId(), record.getAmount(), paymentId);
        return true;
    }
}
