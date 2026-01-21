/**
 * 黑名单服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service.impl;

import com.yupi.airouter.service.BlacklistService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RSet;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;

/**
 * 基于 Redisson 的黑名单服务实现
 * 使用 Redis Set 存储黑名单 IP
 */
@Slf4j
@Service
public class BlacklistServiceImpl implements BlacklistService {

    @Resource
    private RedissonClient redissonClient;

    /**
     * 黑名单 Redis Key
     */
    private static final String BLACKLIST_KEY = "blacklist:ip";

    @Override
    public boolean isBlocked(String ip) {
        if (ip == null || ip.isBlank()) {
            return false;
        }
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        return blacklist.contains(ip);
    }

    @Override
    public void addToBlacklist(String ip, String reason) {
        if (ip == null || ip.isBlank()) {
            return;
        }
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        blacklist.add(ip);
        log.info("IP added to blacklist: {}, reason: {}", ip, reason);
    }

    @Override
    public void removeFromBlacklist(String ip) {
        if (ip == null || ip.isBlank()) {
            return;
        }
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        blacklist.remove(ip);
        log.info("IP removed from blacklist: {}", ip);
    }

    @Override
    public Set<String> getAllBlacklist() {
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        return new HashSet<>(blacklist.readAll());
    }

    @Override
    public void clearBlacklist() {
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        blacklist.delete();
        log.info("Blacklist cleared");
    }

    @Override
    public long getBlacklistCount() {
        RSet<String> blacklist = redissonClient.getSet(BLACKLIST_KEY);
        return blacklist.size();
    }
}
