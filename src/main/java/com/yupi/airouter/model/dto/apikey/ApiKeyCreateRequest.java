package com.yupi.airouter.model.dto.apikey;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 创建 API Key 请求
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
public class ApiKeyCreateRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * Key名称/备注
     */
    private String keyName;
}
