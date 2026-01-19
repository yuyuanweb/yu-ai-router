package com.yupi.airouter.task;

import com.yupi.airouter.service.HealthCheckService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

/**
 * 健康检查定时任务
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Component
public class HealthCheckTask {

    @Resource
    private HealthCheckService healthCheckService;

    /**
     * 定时健康检查
     * 每30秒执行一次
     */
    @Scheduled(fixedRate = 30000)
    public void executeHealthCheck() {
        log.debug("执行定时健康检查任务");
        try {
            healthCheckService.checkAllProviders();
        } catch (Exception e) {
            log.error("健康检查任务执行失败", e);
        }
    }
}
