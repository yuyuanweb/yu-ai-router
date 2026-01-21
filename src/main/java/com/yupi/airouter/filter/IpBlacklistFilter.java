/**
 * IP 黑名单过滤器
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.filter;

import cn.hutool.extra.servlet.JakartaServletUtil;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.service.BlacklistService;
import jakarta.annotation.Resource;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

/**
 * IP 黑名单过滤器
 * 在请求处理的最前面拦截黑名单 IP
 */
@Slf4j
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class IpBlacklistFilter extends OncePerRequestFilter {

    @Resource
    private BlacklistService blacklistService;

    @Resource
    private ObjectMapper objectMapper;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        // 获取客户端 IP
        String clientIp = JakartaServletUtil.getClientIP(request);

        // 检查是否在黑名单中
        if (blacklistService.isBlocked(clientIp)) {
            log.warn("Blocked request from blacklisted IP: {}", clientIp);
            writeErrorResponse(response, ErrorCode.FORBIDDEN_ERROR, "您的 IP 已被封禁");
            return;
        }

        filterChain.doFilter(request, response);
    }

    /**
     * 写入错误响应
     */
    private void writeErrorResponse(HttpServletResponse response, ErrorCode errorCode, String message) throws IOException {
        response.setStatus(HttpServletResponse.SC_FORBIDDEN);
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        response.setCharacterEncoding(StandardCharsets.UTF_8.name());

        BaseResponse<Object> errorResponse = new BaseResponse<>(errorCode.getCode(), null, message);
        response.getWriter().write(objectMapper.writeValueAsString(errorResponse));
    }
}
