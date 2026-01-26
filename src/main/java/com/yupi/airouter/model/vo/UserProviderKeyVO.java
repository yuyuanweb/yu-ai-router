/**
 * 用户提供者密钥视图
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.model.vo;

import com.yupi.airouter.model.entity.UserProviderKey;
import lombok.Data;
import org.springframework.beans.BeanUtils;

import java.io.Serial;
import java.io.Serializable;
import java.util.Date;

@Data
public class UserProviderKeyVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * id
     */
    private Long id;

    /**
     * 提供者 ID
     */
    private Long providerId;

    /**
     * 提供者名称
     */
    private String providerName;

    /**
     * API Key（脱敏显示：只显示前8位和后4位）
     */
    private String apiKey;

    /**
     * 状态：active/inactive
     */
    private String status;

    /**
     * 创建时间
     */
    private Date createTime;

    /**
     * 更新时间
     */
    private Date updateTime;

    /**
     * 对象转包装类
     */
    public static UserProviderKeyVO objToVo(UserProviderKey userProviderKey) {
        if (userProviderKey == null) {
            return null;
        }
        UserProviderKeyVO vo = new UserProviderKeyVO();
        BeanUtils.copyProperties(userProviderKey, vo);
        
        // API Key 脱敏
        String apiKey = userProviderKey.getApiKey();
        if (apiKey != null && apiKey.length() > 12) {
            vo.setApiKey(apiKey.substring(0, 8) + "****" + apiKey.substring(apiKey.length() - 4));
        } else {
            vo.setApiKey("****");
        }
        
        return vo;
    }
}
