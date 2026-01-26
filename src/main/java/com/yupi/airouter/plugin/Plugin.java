/**
 * 插件接口
 * 所有插件都需要实现此接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin;

public interface Plugin {

    /**
     * 获取插件标识
     */
    String getPluginKey();

    /**
     * 获取插件名称
     */
    String getPluginName();

    /**
     * 获取插件描述
     */
    String getDescription();

    /**
     * 执行插件
     *
     * @param context 插件上下文
     * @return 执行结果
     */
    PluginResult execute(PluginContext context);

    /**
     * 检查插件是否支持给定的上下文
     *
     * @param context 插件上下文
     * @return 是否支持
     */
    default boolean supports(PluginContext context) {
        return true;
    }

    /**
     * 初始化插件（在插件加载时调用）
     *
     * @param config 插件配置（JSON 格式）
     */
    default void init(String config) {
        // 默认空实现
    }

    /**
     * 销毁插件（在插件卸载时调用）
     */
    default void destroy() {
        // 默认空实现
    }
}
