package com.yupi.airouter.service;

import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.entity.ApiKey;
import com.yupi.airouter.model.entity.User;

import java.util.List;

/**
 * API Key 服务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
public interface ApiKeyService extends IService<ApiKey> {

    /**
     * 创建 API Key
     *
     * @param keyName Key名称/备注
     * @param loginUser 当前登录用户
     * @return API Key
     */
    ApiKey createApiKey(String keyName, User loginUser);

    /**
     * 获取用户的 API Key 列表
     *
     * @param userId 用户ID
     * @return API Key 列表
     */
    List<ApiKey> listUserApiKeys(Long userId);

    /**
     * 撤销 API Key
     *
     * @param id API Key ID
     * @param userId 用户ID
     * @return 是否成功
     */
    boolean revokeApiKey(Long id, Long userId);

    /**
     * 根据 Key 值查询 API Key
     *
     * @param keyValue Key值
     * @return API Key
     */
    ApiKey getByKeyValue(String keyValue);

    /**
     * 更新 API Key 的使用统计
     *
     * @param apiKeyId API Key ID
     * @param tokens 本次消耗的Token数
     */
    void updateUsageStats(Long apiKeyId, Integer tokens);
}
