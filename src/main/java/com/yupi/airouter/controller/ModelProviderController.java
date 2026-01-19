package com.yupi.airouter.controller;

import cn.hutool.core.bean.BeanUtil;
import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.annotation.AuthCheck;
import com.yupi.airouter.common.BaseResponse;
import com.yupi.airouter.common.DeleteRequest;
import com.yupi.airouter.common.ResultUtils;
import com.yupi.airouter.constant.UserConstant;
import com.yupi.airouter.exception.ErrorCode;
import com.yupi.airouter.exception.ThrowUtils;
import com.yupi.airouter.model.dto.provider.ProviderAddRequest;
import com.yupi.airouter.model.dto.provider.ProviderQueryRequest;
import com.yupi.airouter.model.dto.provider.ProviderUpdateRequest;
import com.yupi.airouter.model.entity.ModelProvider;
import com.yupi.airouter.model.enums.ProviderStatusEnum;
import com.yupi.airouter.model.vo.ProviderVO;
import com.yupi.airouter.service.ModelProviderService;
import jakarta.annotation.Resource;
import org.springframework.web.bind.annotation.*;

import java.util.List;

/**
 * 模型提供者控制层
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/provider")
public class ModelProviderController {

    @Resource
    private ModelProviderService modelProviderService;

    /**
     * 创建模型提供者（管理员）
     *
     * @param providerAddRequest 创建请求
     * @return 创建的提供者ID
     */
    @PostMapping("/add")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Long> addProvider(@RequestBody ProviderAddRequest providerAddRequest) {
        ThrowUtils.throwIf(providerAddRequest == null, ErrorCode.PARAMS_ERROR);
        
        ModelProvider modelProvider = new ModelProvider();
        BeanUtil.copyProperties(providerAddRequest, modelProvider);
        
        // 设置默认值
        if (modelProvider.getPriority() == null) {
            modelProvider.setPriority(100);
        }
        modelProvider.setStatus(ProviderStatusEnum.ACTIVE.getValue());
        
        // 保存到数据库
        boolean result = modelProviderService.save(modelProvider);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);
        
        return ResultUtils.success(modelProvider.getId());
    }

    /**
     * 删除模型提供者（管理员）
     *
     * @param deleteRequest 删除请求
     * @return 是否删除成功
     */
    @PostMapping("/delete")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> deleteProvider(@RequestBody DeleteRequest deleteRequest) {
        ThrowUtils.throwIf(deleteRequest == null || deleteRequest.getId() <= 0, ErrorCode.PARAMS_ERROR);
        
        boolean result = modelProviderService.removeById(deleteRequest.getId());
        return ResultUtils.success(result);
    }

    /**
     * 更新模型提供者（管理员）
     *
     * @param providerUpdateRequest 更新请求
     * @return 是否更新成功
     */
    @PostMapping("/update")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> updateProvider(@RequestBody ProviderUpdateRequest providerUpdateRequest) {
        ThrowUtils.throwIf(providerUpdateRequest == null || providerUpdateRequest.getId() == null, ErrorCode.PARAMS_ERROR);
        
        ModelProvider modelProvider = new ModelProvider();
        BeanUtil.copyProperties(providerUpdateRequest, modelProvider);
        
        boolean result = modelProviderService.updateById(modelProvider);
        return ResultUtils.success(result);
    }

    /**
     * 根据 ID 获取模型提供者（脱敏）
     *
     * @param id 提供者ID
     * @return 提供者信息
     */
    @GetMapping("/get/vo")
    public BaseResponse<ProviderVO> getProviderVOById(long id) {
        ThrowUtils.throwIf(id <= 0, ErrorCode.PARAMS_ERROR);
        
        ModelProvider modelProvider = modelProviderService.getById(id);
        ThrowUtils.throwIf(modelProvider == null, ErrorCode.NOT_FOUND_ERROR);
        
        return ResultUtils.success(modelProviderService.getProviderVO(modelProvider));
    }

    /**
     * 分页获取模型提供者列表（脱敏）
     *
     * @param providerQueryRequest 查询请求
     * @return 分页结果
     */
    @PostMapping("/list/page/vo")
    public BaseResponse<Page<ProviderVO>> listProviderVOByPage(@RequestBody ProviderQueryRequest providerQueryRequest) {
        long pageNum = providerQueryRequest.getPageNum();
        long pageSize = providerQueryRequest.getPageSize();
        
        // 查询数据库
        QueryWrapper queryWrapper = modelProviderService.getQueryWrapper(providerQueryRequest);
        Page<ModelProvider> providerPage = modelProviderService.page(Page.of(pageNum, pageSize), queryWrapper);
        
        // 转换为 VO
        Page<ProviderVO> providerVOPage = new Page<>(pageNum, pageSize, providerPage.getTotalRow());
        providerVOPage.setRecords(modelProviderService.getProviderVOList(providerPage.getRecords()));
        
        return ResultUtils.success(providerVOPage);
    }

    /**
     * 获取所有模型提供者列表（脱敏）
     *
     * @return 提供者列表
     */
    @GetMapping("/list/vo")
    public BaseResponse<List<ProviderVO>> listProviderVO() {
        List<ModelProvider> providerList = modelProviderService.list();
        List<ProviderVO> providerVOList = modelProviderService.getProviderVOList(providerList);
        return ResultUtils.success(providerVOList);
    }

    /**
     * 获取健康的提供者列表
     *
     * @return 健康的提供者列表
     */
    @GetMapping("/list/healthy")
    public BaseResponse<List<ProviderVO>> listHealthyProviders() {
        List<ModelProvider> healthyProviders = modelProviderService.getHealthyProviders();
        List<ProviderVO> providerVOList = modelProviderService.getProviderVOList(healthyProviders);
        return ResultUtils.success(providerVOList);
    }
}
