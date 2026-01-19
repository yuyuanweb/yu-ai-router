package com.yupi.airouter.model.dto.model;

import com.yupi.airouter.common.PageRequest;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serial;
import java.io.Serializable;

/**
 * 模型查询请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@EqualsAndHashCode(callSuper = true)
public class ModelQueryRequest extends PageRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 提供者id
     */
    private Long providerId;

    /**
     * 模型标识
     */
    private String modelKey;

    /**
     * 模型显示名称
     */
    private String modelName;

    /**
     * 模型类型：chat/embedding/image/audio
     */
    private String modelType;

    /**
     * 状态：active/inactive/deprecated
     */
    private String status;
}
