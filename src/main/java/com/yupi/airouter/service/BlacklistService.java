/**
 * 黑名单服务接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service;

import java.util.Set;

/**
 * 黑名单服务
 * 提供 IP 黑名单管理能力
 */
public interface BlacklistService {

    /**
     * 检查 IP 是否在黑名单中
     *
     * @param ip IP 地址
     * @return true 表示在黑名单中
     */
    boolean isBlocked(String ip);

    /**
     * 将 IP 加入黑名单
     *
     * @param ip     IP 地址
     * @param reason 封禁原因
     */
    void addToBlacklist(String ip, String reason);

    /**
     * 将 IP 从黑名单移除
     *
     * @param ip IP 地址
     */
    void removeFromBlacklist(String ip);

    /**
     * 获取所有黑名单 IP
     *
     * @return 黑名单 IP 集合
     */
    Set<String> getAllBlacklist();

    /**
     * 清空黑名单
     */
    void clearBlacklist();

    /**
     * 获取黑名单数量
     *
     * @return 黑名单数量
     */
    long getBlacklistCount();
}
