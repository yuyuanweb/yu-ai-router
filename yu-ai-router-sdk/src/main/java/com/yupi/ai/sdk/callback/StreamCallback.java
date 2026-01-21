/**
 * 流式回调接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.ai.sdk.callback;

import com.yupi.ai.sdk.model.ChatChunk;

/**
 * 流式响应回调接口
 */
public interface StreamCallback {

    /**
     * 接收到消息块时调用
     *
     * @param chunk 消息块
     */
    void onMessage(ChatChunk chunk);

    /**
     * 流式响应完成时调用
     */
    void onComplete();

    /**
     * 发生错误时调用
     *
     * @param error 异常
     */
    void onError(Throwable error);
}
