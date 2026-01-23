package com.yupi.airouter.service;

import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.entity.RechargeRecord;

/**
 * 充值记录服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface RechargeService extends IService<RechargeRecord> {

    /**
     * 创建充值记录
     *
     * @param userId 用户ID
     * @param amount 充值金额
     * @param paymentMethod 支付方式
     * @return 充值记录
     */
    RechargeRecord createRechargeRecord(Long userId, java.math.BigDecimal amount, String paymentMethod);

    /**
     * 更新充值记录状态
     *
     * @param recordId 记录ID
     * @param status 状态
     * @param paymentId 支付ID
     * @return 是否更新成功
     */
    boolean updateRechargeStatus(Long recordId, String status, String paymentId);

    /**
     * 根据支付ID获取充值记录
     *
     * @param paymentId 支付ID
     * @return 充值记录
     */
    RechargeRecord getByPaymentId(String paymentId);

    /**
     * 获取用户的充值记录列表（分页）
     *
     * @param userId 用户ID
     * @param pageNum 页码
     * @param pageSize 每页大小
     * @return 分页结果
     */
    Page<RechargeRecord> listUserRechargeRecords(Long userId, int pageNum, int pageSize);

    /**
     * 完成充值（更新状态并增加余额）
     *
     * @param recordId 记录ID
     * @param paymentId 支付ID
     * @return 是否成功
     */
    boolean completeRecharge(Long recordId, String paymentId);
}
