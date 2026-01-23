package com.yupi.airouter.service;

import com.yupi.airouter.model.entity.User;

import java.math.BigDecimal;

/**
 * 余额管理服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface BalanceService {

    /**
     * 检查用户余额是否充足
     *
     * @param userId 用户ID
     * @param amount 需要的金额
     * @return 是否充足
     */
    boolean checkBalance(Long userId, BigDecimal amount);

    /**
     * 扣减用户余额
     *
     * @param userId 用户ID
     * @param amount 扣减金额
     * @param requestLogId 关联的请求日志ID
     * @param description 扣减说明
     * @return 是否扣减成功
     */
    boolean deductBalance(Long userId, BigDecimal amount, Long requestLogId, String description);

    /**
     * 增加用户余额（充值）
     *
     * @param userId 用户ID
     * @param amount 充值金额
     * @param description 充值说明
     * @return 是否充值成功
     */
    boolean addBalance(Long userId, BigDecimal amount, String description);

    /**
     * 获取用户余额
     *
     * @param userId 用户ID
     * @return 用户余额
     */
    BigDecimal getUserBalance(Long userId);

    /**
     * 更新用户余额
     *
     * @param user 用户
     * @return 是否更新成功
     */
    boolean updateBalance(User user);
}
