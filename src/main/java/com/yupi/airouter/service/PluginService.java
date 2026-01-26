/**
 * 插件服务接口
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service;

import com.mybatisflex.core.service.IService;
import com.yupi.airouter.model.dto.plugin.PluginExecuteRequest;
import com.yupi.airouter.model.dto.plugin.PluginUpdateRequest;
import com.yupi.airouter.model.entity.PluginConfig;
import com.yupi.airouter.model.vo.PluginConfigVO;
import com.yupi.airouter.model.vo.PluginExecuteVO;

import java.util.List;

public interface PluginService extends IService<PluginConfig> {

    /**
     * 获取所有插件配置
     *
     * @return 插件配置列表
     */
    List<PluginConfigVO> listAllPlugins();

    /**
     * 获取所有启用的插件
     *
     * @return 启用的插件配置列表
     */
    List<PluginConfigVO> listEnabledPlugins();

    /**
     * 根据插件标识获取插件配置
     *
     * @param pluginKey 插件标识
     * @return 插件配置
     */
    PluginConfigVO getPluginByKey(String pluginKey);

    /**
     * 更新插件配置
     *
     * @param request 更新请求
     * @return 是否成功
     */
    boolean updatePlugin(PluginUpdateRequest request);

    /**
     * 启用插件
     *
     * @param pluginKey 插件标识
     * @return 是否成功
     */
    boolean enablePlugin(String pluginKey);

    /**
     * 禁用插件
     *
     * @param pluginKey 插件标识
     * @return 是否成功
     */
    boolean disablePlugin(String pluginKey);

    /**
     * 执行插件
     *
     * @param request 执行请求
     * @param userId  用户 ID
     * @return 执行结果
     */
    PluginExecuteVO executePlugin(PluginExecuteRequest request, Long userId);

    /**
     * 初始化插件（应用启动时调用）
     */
    void initPlugins();

    /**
     * 重新加载插件
     *
     * @param pluginKey 插件标识
     */
    void reloadPlugin(String pluginKey);
}
