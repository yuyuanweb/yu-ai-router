package com.yupi.airouter.service;

/**
 * 用户配额服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface QuotaService {

    /**
     * 检查用户配额是否充足
     *
     * @param userId 用户ID
     * @return true-配额充足，false-配额不足
     */
    boolean checkQuota(Long userId);

    /**
     * 扣减用户Token配额
     *
     * @param userId 用户ID
     * @param tokens 消耗的Token数量
     * @return true-扣减成功，false-扣减失败
     */
    boolean deductTokens(Long userId, int tokens);

    /**
     * 获取用户剩余配额
     *
     * @param userId 用户ID
     * @return 剩余配额，-1表示无限制
     */
    long getRemainingQuota(Long userId);
}
