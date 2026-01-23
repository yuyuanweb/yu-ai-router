package com.yupi.airouter.model.vo;

import lombok.Data;

import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * 脱敏后的登录用户信息
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class LoginUserVO implements Serializable {

    /**
     * 用户 id
     */
    private Long id;

    /**
     * 账号
     */
    private String userAccount;

    /**
     * 用户昵称
     */
    private String userName;

    /**
     * 用户头像
     */
    private String userAvatar;

    /**
     * 用户简介
     */
    private String userProfile;

    /**
     * 用户角色：user/admin
     */
    private String userRole;

    /**
     * 用户状态：active/disabled
     */
    private String userStatus;

    /**
     * Token配额（-1表示无限制）
     */
    private Long tokenQuota;

    /**
     * 已使用Token数
     */
    private Long usedTokens;

    /**
     * 创建时间
     */
    private LocalDateTime createTime;

    /**
     * 更新时间
     */
    private LocalDateTime updateTime;

    private static final long serialVersionUID = 1L;
}