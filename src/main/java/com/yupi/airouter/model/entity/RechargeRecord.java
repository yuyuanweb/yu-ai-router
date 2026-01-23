package com.yupi.airouter.model.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import com.mybatisflex.core.keygen.KeyGenerators;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDateTime;

/**
 * 充值记录 实体类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("recharge_record")
public class RechargeRecord implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    @Id(keyType = KeyType.Generator, value = KeyGenerators.snowFlakeId)
    private Long id;

    /**
     * 用户id
     */
    @Column("userId")
    private Long userId;

    /**
     * 充值金额（元）
     */
    @Column("amount")
    private BigDecimal amount;

    /**
     * 支付方式：stripe/alipay/wechat
     */
    @Column("paymentMethod")
    private String paymentMethod;

    /**
     * 第三方支付ID
     */
    @Column("paymentId")
    private String paymentId;

    /**
     * 状态：pending/success/failed/refunded
     */
    @Column("status")
    private String status;

    /**
     * 充值说明
     */
    @Column("description")
    private String description;

    /**
     * 创建时间
     */
    @Column("createTime")
    private LocalDateTime createTime;

    /**
     * 更新时间
     */
    @Column("updateTime")
    private LocalDateTime updateTime;

}
