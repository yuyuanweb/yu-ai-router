/**
 * 用户提供者密钥服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service.impl;

import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.mapper.UserProviderKeyMapper;
import com.yupi.airouter.model.dto.byok.UserProviderKeyAddRequest;
import com.yupi.airouter.model.dto.byok.UserProviderKeyUpdateRequest;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.model.entity.UserProviderKey;
import com.yupi.airouter.model.vo.UserProviderKeyVO;
import com.yupi.airouter.service.ModelProviderService;
import com.yupi.airouter.service.UserProviderKeyService;
import com.yupi.airouter.utils.EncryptionUtils;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

@Slf4j
@Service
public class UserProviderKeyServiceImpl extends ServiceImpl<UserProviderKeyMapper, UserProviderKey> 
        implements UserProviderKeyService {

    @Resource
    private EncryptionUtils encryptionUtils;

    @Resource
    private ModelProviderService modelProviderService;

    @Override
    public boolean addUserProviderKey(UserProviderKeyAddRequest request, Long userId) {
        if (request.getProviderId() == null || StringUtils.isBlank(request.getApiKey())) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }

        // 检查提供者是否存在
        ModelProvider provider = modelProviderService.getById(request.getProviderId());
        if (provider == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "提供者不存在");
        }

        // 检查是否已经配置过该提供者的密钥
        UserProviderKey existing = getOne(QueryWrapper.create()
                .eq("userId", userId)
                .eq("providerId", request.getProviderId()));
        
        if (existing != null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "已配置过该提供者的密钥，请使用更新功能");
        }

        // 加密 API Key
        String encryptedApiKey = encryptionUtils.encrypt(request.getApiKey());

        // 创建记录
        UserProviderKey userProviderKey = UserProviderKey.builder()
                .userId(userId)
                .providerId(request.getProviderId())
                .providerName(provider.getProviderName())
                .apiKey(encryptedApiKey)
                .status("active")
                .build();

        boolean result = save(userProviderKey);
        if (result) {
            log.info("用户 {} 添加了提供者 {} 的密钥", userId, provider.getProviderName());
        }
        return result;
    }

    @Override
    public boolean updateUserProviderKey(UserProviderKeyUpdateRequest request, Long userId) {
        if (request.getId() == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR);
        }

        // 查询原记录
        UserProviderKey userProviderKey = getById(request.getId());
        if (userProviderKey == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR);
        }

        // 验证是否是当前用户的密钥
        if (!userProviderKey.getUserId().equals(userId)) {
            throw new BusinessException(ErrorCode.NO_AUTH_ERROR);
        }

        // 更新字段
        if (StringUtils.isNotBlank(request.getApiKey())) {
            // 加密新的 API Key
            String encryptedApiKey = encryptionUtils.encrypt(request.getApiKey());
            userProviderKey.setApiKey(encryptedApiKey);
        }
        if (StringUtils.isNotBlank(request.getStatus())) {
            userProviderKey.setStatus(request.getStatus());
        }

        return updateById(userProviderKey);
    }

    @Override
    public boolean deleteUserProviderKey(Long id, Long userId) {
        UserProviderKey userProviderKey = getById(id);
        if (userProviderKey == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR);
        }

        // 验证是否是当前用户的密钥
        if (!userProviderKey.getUserId().equals(userId)) {
            throw new BusinessException(ErrorCode.NO_AUTH_ERROR);
        }

        return removeById(id);
    }

    @Override
    public List<UserProviderKeyVO> listUserProviderKeys(Long userId) {
        List<UserProviderKey> list = list(QueryWrapper.create()
                .eq("userId", userId)
                .orderBy("createTime", false));
        
        return list.stream()
                .map(UserProviderKeyVO::objToVo)
                .collect(Collectors.toList());
    }

    @Override
    public String getUserProviderApiKey(Long userId, Long providerId) {
        UserProviderKey userProviderKey = getOne(QueryWrapper.create()
                .eq("userId", userId)
                .eq("providerId", providerId)
                .eq("status", "active"));
        
        if (userProviderKey == null) {
            return null;
        }

        // 解密 API Key
        try {
            return encryptionUtils.decrypt(userProviderKey.getApiKey());
        } catch (Exception e) {
            log.error("解密用户密钥失败", e);
            return null;
        }
    }

    @Override
    public boolean hasUserProviderKey(Long userId, Long providerId) {
        return getUserProviderApiKey(userId, providerId) != null;
    }
}
