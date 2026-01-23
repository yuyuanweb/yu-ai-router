package com.yupi.airouter.service.impl;

import com.yupi.airouter.mapper.UserMapper;
import com.yupi.airouter.model.entity.User;
import com.yupi.airouter.service.QuotaService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

/**
 * 用户配额服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class QuotaServiceImpl implements QuotaService {

    @Resource
    private UserMapper userMapper;

    /**
     * 无限制配额标识
     */
    private static final long UNLIMITED_QUOTA = -1L;

    @Override
    public boolean checkQuota(Long userId) {
        if (userId == null) {
            return true;
        }
        
        User user = userMapper.selectOneById(userId);
        if (user == null) {
            return false;
        }
        
        Long tokenQuota = user.getTokenQuota();
        // 无限制配额
        if (tokenQuota == null || tokenQuota == UNLIMITED_QUOTA) {
            return true;
        }
        
        Long usedTokens = user.getUsedTokens();
        if (usedTokens == null) {
            usedTokens = 0L;
        }
        
        // 检查是否还有剩余配额
        return usedTokens < tokenQuota;
    }

    @Override
    public boolean deductTokens(Long userId, int tokens) {
        if (userId == null || tokens <= 0) {
            return true;
        }
        
        try {
            // 先查询当前用户的已使用Token数
            User user = userMapper.selectOneById(userId);
            if (user == null) {
                return false;
            }
            
            Long currentUsedTokens = user.getUsedTokens();
            if (currentUsedTokens == null) {
                currentUsedTokens = 0L;
            }
            
            // 更新已使用Token数
            User updateUser = new User();
            updateUser.setId(userId);
            updateUser.setUsedTokens(currentUsedTokens + tokens);
            int updated = userMapper.update(updateUser);
            
            if (updated > 0) {
                log.debug("用户 {} 扣减Token {} 成功", userId, tokens);
            }
            return updated > 0;
        } catch (Exception e) {
            log.error("用户 {} 扣减Token失败", userId, e);
            return false;
        }
    }

    @Override
    public long getRemainingQuota(Long userId) {
        if (userId == null) {
            return UNLIMITED_QUOTA;
        }
        
        User user = userMapper.selectOneById(userId);
        if (user == null) {
            return 0;
        }
        
        Long tokenQuota = user.getTokenQuota();
        // 无限制配额
        if (tokenQuota == null || tokenQuota == UNLIMITED_QUOTA) {
            return UNLIMITED_QUOTA;
        }
        
        Long usedTokens = user.getUsedTokens();
        if (usedTokens == null) {
            usedTokens = 0L;
        }
        
        return Math.max(0, tokenQuota - usedTokens);
    }
}
