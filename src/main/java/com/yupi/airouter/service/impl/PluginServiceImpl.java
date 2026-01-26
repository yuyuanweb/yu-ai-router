/**
 * 插件服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.service.impl;

import cn.hutool.core.bean.BeanUtil;
import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import com.yupi.airouter.exception.BusinessException;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.mapper.PluginConfigMapper;
import com.yupi.airouter.model.dto.plugin.PluginExecuteRequest;
import com.yupi.airouter.model.dto.plugin.PluginUpdateRequest;
import com.yupi.airouter.model.entity.PluginConfig;
import com.yupi.airouter.model.enums.PluginStatusEnum;
import com.yupi.airouter.model.vo.PluginConfigVO;
import com.yupi.airouter.model.vo.PluginExecuteVO;
import com.yupi.airouter.plugin.Plugin;
import com.yupi.airouter.plugin.PluginContext;
import com.yupi.airouter.plugin.PluginManager;
import com.yupi.airouter.plugin.PluginResult;
import com.yupi.airouter.service.PluginService;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationContext;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Slf4j
@Service
public class PluginServiceImpl extends ServiceImpl<PluginConfigMapper, PluginConfig> implements PluginService {

    @Resource
    private PluginManager pluginManager;

    @Resource
    private ApplicationContext applicationContext;

    @PostConstruct
    public void init() {
        initPlugins();
    }

    @Override
    public List<PluginConfigVO> listAllPlugins() {
        List<PluginConfig> pluginConfigs = list(QueryWrapper.create()
                .orderBy("priority", false));
        return pluginConfigs.stream()
                .map(this::convertToVO)
                .collect(Collectors.toList());
    }

    @Override
    public List<PluginConfigVO> listEnabledPlugins() {
        List<PluginConfig> pluginConfigs = list(QueryWrapper.create()
                .eq("status", PluginStatusEnum.ACTIVE.getValue())
                .orderBy("priority", false));
        return pluginConfigs.stream()
                .map(this::convertToVO)
                .collect(Collectors.toList());
    }

    @Override
    public PluginConfigVO getPluginByKey(String pluginKey) {
        PluginConfig pluginConfig = getOne(QueryWrapper.create()
                .eq("pluginKey", pluginKey));
        if (pluginConfig == null) {
            return null;
        }
        return convertToVO(pluginConfig);
    }

    @Override
    public boolean updatePlugin(PluginUpdateRequest request) {
        if (request.getId() == null) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件 ID 不能为空");
        }

        PluginConfig pluginConfig = getById(request.getId());
        if (pluginConfig == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在");
        }

        // 更新字段
        if (request.getPluginName() != null) {
            pluginConfig.setPluginName(request.getPluginName());
        }
        if (request.getDescription() != null) {
            pluginConfig.setDescription(request.getDescription());
        }
        if (request.getConfig() != null) {
            pluginConfig.setConfig(request.getConfig());
        }
        if (request.getStatus() != null) {
            pluginConfig.setStatus(request.getStatus());
        }
        if (request.getPriority() != null) {
            pluginConfig.setPriority(request.getPriority());
        }

        boolean result = updateById(pluginConfig);
        if (result) {
            // 更新插件管理器中的配置
            pluginManager.updatePluginConfig(pluginConfig);
        }
        return result;
    }

    @Override
    public boolean enablePlugin(String pluginKey) {
        PluginConfig pluginConfig = getOne(QueryWrapper.create()
                .eq("pluginKey", pluginKey));
        if (pluginConfig == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在");
        }

        pluginConfig.setStatus(PluginStatusEnum.ACTIVE.getValue());
        boolean result = updateById(pluginConfig);
        if (result) {
            pluginManager.updatePluginConfig(pluginConfig);
            log.info("启用插件: {}", pluginKey);
        }
        return result;
    }

    @Override
    public boolean disablePlugin(String pluginKey) {
        PluginConfig pluginConfig = getOne(QueryWrapper.create()
                .eq("pluginKey", pluginKey));
        if (pluginConfig == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在");
        }

        pluginConfig.setStatus(PluginStatusEnum.INACTIVE.getValue());
        boolean result = updateById(pluginConfig);
        if (result) {
            pluginManager.updatePluginConfig(pluginConfig);
            log.info("禁用插件: {}", pluginKey);
        }
        return result;
    }

    @Override
    public PluginExecuteVO executePlugin(PluginExecuteRequest request, Long userId) {
        String pluginKey = request.getPluginKey();
        if (pluginKey == null || pluginKey.isEmpty()) {
            throw new BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空");
        }

        // 构建插件上下文
        PluginContext context = PluginContext.builder()
                .userId(userId)
                .input(request.getInput())
                .fileUrl(request.getFileUrl())
                .fileBytes(request.getFileBytes())
                .fileType(request.getFileType())
                .params(request.getParams())
                .build();

        // 执行插件
        PluginResult result = pluginManager.executePlugin(pluginKey, context);

        // 转换结果
        return PluginExecuteVO.builder()
                .success(result.isSuccess())
                .pluginKey(pluginKey)
                .content(result.getContent())
                .errorMessage(result.getErrorMessage())
                .duration(result.getDuration())
                .data(result.getData())
                .build();
    }

    @Override
    public void initPlugins() {
        log.info("初始化插件系统...");

        // 从数据库加载所有插件配置
        List<PluginConfig> pluginConfigs = list();

        // 获取所有实现了 Plugin 接口的 Bean
        Map<String, Plugin> pluginBeans = applicationContext.getBeansOfType(Plugin.class);

        // 注册插件
        for (Plugin plugin : pluginBeans.values()) {
            String pluginKey = plugin.getPluginKey();
            PluginConfig config = pluginConfigs.stream()
                    .filter(c -> c.getPluginKey().equals(pluginKey))
                    .findFirst()
                    .orElse(null);

            if (config != null) {
                pluginManager.registerPlugin(plugin, config);
            } else {
                log.warn("插件 {} 在数据库中没有配置", pluginKey);
            }
        }

        log.info("插件系统初始化完成，共加载 {} 个插件", pluginManager.getPluginCount());
    }

    @Override
    public void reloadPlugin(String pluginKey) {
        // 获取插件配置
        PluginConfig pluginConfig = getOne(QueryWrapper.create()
                .eq("pluginKey", pluginKey));
        if (pluginConfig == null) {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在");
        }

        // 获取插件实例
        Map<String, Plugin> pluginBeans = applicationContext.getBeansOfType(Plugin.class);
        Plugin plugin = pluginBeans.values().stream()
                .filter(p -> p.getPluginKey().equals(pluginKey))
                .findFirst()
                .orElse(null);

        if (plugin != null) {
            pluginManager.registerPlugin(plugin, pluginConfig);
            log.info("重新加载插件: {}", pluginKey);
        } else {
            throw new BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件实例不存在");
        }
    }

    /**
     * 转换为 VO
     */
    private PluginConfigVO convertToVO(PluginConfig pluginConfig) {
        PluginConfigVO vo = new PluginConfigVO();
        BeanUtil.copyProperties(pluginConfig, vo);
        return vo;
    }
}
