/**
 * 插件管理器
 * 负责插件的加载、管理和执行
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.plugin;

import com.yupi.airouter.model.entity.PluginConfig;
import com.yupi.airouter.model.enums.PluginStatusEnum;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Component
public class PluginManager {

    /**
     * 已注册的插件（pluginKey -> Plugin）
     */
    private final Map<String, Plugin> plugins = new ConcurrentHashMap<>();

    /**
     * 插件配置（pluginKey -> PluginConfig）
     */
    private final Map<String, PluginConfig> pluginConfigs = new ConcurrentHashMap<>();

    @PostConstruct
    public void init() {
        log.info("插件管理器初始化完成");
    }

    @PreDestroy
    public void destroy() {
        log.info("插件管理器销毁，卸载所有插件");
        plugins.values().forEach(Plugin::destroy);
        plugins.clear();
        pluginConfigs.clear();
    }

    /**
     * 注册插件
     *
     * @param plugin 插件实例
     * @param config 插件配置
     */
    public void registerPlugin(Plugin plugin, PluginConfig config) {
        String pluginKey = plugin.getPluginKey();
        
        // 如果已存在，先销毁旧插件
        if (plugins.containsKey(pluginKey)) {
            Plugin oldPlugin = plugins.get(pluginKey);
            oldPlugin.destroy();
            log.info("卸载旧插件: {}", pluginKey);
        }

        // 初始化并注册新插件
        plugin.init(config != null ? config.getConfig() : null);
        plugins.put(pluginKey, plugin);
        if (config != null) {
            pluginConfigs.put(pluginKey, config);
        }
        log.info("注册插件: {} ({})", pluginKey, plugin.getPluginName());
    }

    /**
     * 注销插件
     *
     * @param pluginKey 插件标识
     */
    public void unregisterPlugin(String pluginKey) {
        Plugin plugin = plugins.remove(pluginKey);
        pluginConfigs.remove(pluginKey);
        if (plugin != null) {
            plugin.destroy();
            log.info("注销插件: {}", pluginKey);
        }
    }

    /**
     * 获取插件
     *
     * @param pluginKey 插件标识
     * @return 插件实例
     */
    public Optional<Plugin> getPlugin(String pluginKey) {
        return Optional.ofNullable(plugins.get(pluginKey));
    }

    /**
     * 获取插件配置
     *
     * @param pluginKey 插件标识
     * @return 插件配置
     */
    public Optional<PluginConfig> getPluginConfig(String pluginKey) {
        return Optional.ofNullable(pluginConfigs.get(pluginKey));
    }

    /**
     * 更新插件配置
     *
     * @param config 新配置
     */
    public void updatePluginConfig(PluginConfig config) {
        String pluginKey = config.getPluginKey();
        pluginConfigs.put(pluginKey, config);
        
        // 重新初始化插件
        Plugin plugin = plugins.get(pluginKey);
        if (plugin != null) {
            plugin.init(config.getConfig());
            log.info("更新插件配置: {}", pluginKey);
        }
    }

    /**
     * 检查插件是否启用
     *
     * @param pluginKey 插件标识
     * @return 是否启用
     */
    public boolean isPluginEnabled(String pluginKey) {
        PluginConfig config = pluginConfigs.get(pluginKey);
        if (config == null) {
            return false;
        }
        return PluginStatusEnum.ACTIVE.getValue().equals(config.getStatus());
    }

    /**
     * 执行插件
     *
     * @param pluginKey 插件标识
     * @param context   执行上下文
     * @return 执行结果
     */
    public PluginResult executePlugin(String pluginKey, PluginContext context) {
        // 检查插件是否存在
        Plugin plugin = plugins.get(pluginKey);
        if (plugin == null) {
            return PluginResult.fail("插件不存在: " + pluginKey);
        }

        // 检查插件是否启用
        if (!isPluginEnabled(pluginKey)) {
            return PluginResult.fail("插件未启用: " + pluginKey);
        }

        // 检查插件是否支持当前上下文
        if (!plugin.supports(context)) {
            return PluginResult.fail("插件不支持当前请求");
        }

        // 执行插件
        long startTime = System.currentTimeMillis();
        try {
            PluginResult result = plugin.execute(context);
            result.setDuration(System.currentTimeMillis() - startTime);
            log.info("插件执行完成: {}, 耗时: {}ms", pluginKey, result.getDuration());
            return result;
        } catch (Exception e) {
            log.error("插件执行失败: {}", pluginKey, e);
            PluginResult result = PluginResult.fail("插件执行异常: " + e.getMessage());
            result.setDuration(System.currentTimeMillis() - startTime);
            return result;
        }
    }

    /**
     * 获取所有已注册的插件
     *
     * @return 插件列表
     */
    public List<Plugin> getAllPlugins() {
        return new ArrayList<>(plugins.values());
    }

    /**
     * 获取所有启用的插件
     *
     * @return 启用的插件列表
     */
    public List<Plugin> getEnabledPlugins() {
        List<Plugin> enabledPlugins = new ArrayList<>();
        for (Map.Entry<String, Plugin> entry : plugins.entrySet()) {
            if (isPluginEnabled(entry.getKey())) {
                enabledPlugins.add(entry.getValue());
            }
        }
        return enabledPlugins;
    }

    /**
     * 获取已注册插件数量
     */
    public int getPluginCount() {
        return plugins.size();
    }
}
