/**
 * 用户提供者密钥（BYOK）实体
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Date;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("user_provider_key")
public class UserProviderKey implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Auto)
    private Long id;

    /**
     * 用户 ID
     */
    @Column("userId")
    private Long userId;

    /**
     * 提供者 ID
     */
    @Column("providerId")
    private Long providerId;

    /**
     * 提供者名称（冗余字段，便于查询）
     */
    @Column("providerName")
    private String providerName;

    /**
     * API Key（加密存储）
     */
    @Column("apiKey")
    private String apiKey;

    /**
     * 状态：active/inactive
     */
    @Column("status")
    private String status;

    /**
     * 创建时间
     */
    @Column("createTime")
    private Date createTime;

    /**
     * 更新时间
     */
    @Column("updateTime")
    private Date updateTime;

    /**
     * 是否删除
     */
    @Column(value = "isDelete", isLogicDelete = true)
    private Integer isDelete;
}
