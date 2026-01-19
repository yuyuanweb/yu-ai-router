package com.yupi.airouter.model.dto.chat;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;

/**
 * 聊天消息
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ChatMessage implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    /**
     * 角色：system/user/assistant
     */
    private String role;

    /**
     * 消息内容
     */
    private String content;
}
