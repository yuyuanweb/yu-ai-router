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
import com.yupi.airouter.model.dto.model.ModelAddRequest;
import com.yupi.airouter.model.dto.model.ModelQueryRequest;
import com.yupi.airouter.model.dto.model.ModelUpdateRequest;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.enums.ModelStatusEnum;
import com.yupi.airouter.model.vo.ModelVO;
import com.yupi.airouter.service.ModelService;
import jakarta.annotation.Resource;
import org.springframework.web.bind.annotation.*;

import java.util.List;

/**
 * 模型控制层
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@RestController
@RequestMapping("/model")
public class ModelController {

    @Resource
    private ModelService modelService;

    /**
     * 创建模型（管理员）
     *
     * @param modelAddRequest 创建请求
     * @return 创建的模型ID
     */
    @PostMapping("/add")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Long> addModel(@RequestBody ModelAddRequest modelAddRequest) {
        ThrowUtils.throwIf(modelAddRequest == null, ErrorCode.PARAMS_ERROR);

        Model model = new Model();
        BeanUtil.copyProperties(modelAddRequest, model);

        // 设置默认值
        if (model.getPriority() == null) {
            model.setPriority(100);
        }
        if (model.getDefaultTimeout() == null) {
            model.setDefaultTimeout(60000);
        }
        if (model.getContextLength() == null) {
            model.setContextLength(4096);
        }
        model.setStatus(ModelStatusEnum.ACTIVE.getValue());

        // 保存到数据库
        boolean result = modelService.save(model);
        ThrowUtils.throwIf(!result, ErrorCode.OPERATION_ERROR);

        return ResultUtils.success(model.getId());
    }

    /**
     * 删除模型（管理员）
     *
     * @param deleteRequest 删除请求
     * @return 是否删除成功
     */
    @PostMapping("/delete")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> deleteModel(@RequestBody DeleteRequest deleteRequest) {
        ThrowUtils.throwIf(deleteRequest == null || deleteRequest.getId() <= 0, ErrorCode.PARAMS_ERROR);

        boolean result = modelService.removeById(deleteRequest.getId());
        return ResultUtils.success(result);
    }

    /**
     * 更新模型（管理员）
     *
     * @param modelUpdateRequest 更新请求
     * @return 是否更新成功
     */
    @PostMapping("/update")
    @AuthCheck(mustRole = UserConstant.ADMIN_ROLE)
    public BaseResponse<Boolean> updateModel(@RequestBody ModelUpdateRequest modelUpdateRequest) {
        ThrowUtils.throwIf(modelUpdateRequest == null || modelUpdateRequest.getId() == null, ErrorCode.PARAMS_ERROR);

        Model model = new Model();
        BeanUtil.copyProperties(modelUpdateRequest, model);

        boolean result = modelService.updateById(model);
        return ResultUtils.success(result);
    }

    /**
     * 根据 ID 获取模型（含提供者信息）
     *
     * @param id 模型ID
     * @return 模型信息
     */
    @GetMapping("/get/vo")
    public BaseResponse<ModelVO> getModelVOById(long id) {
        ThrowUtils.throwIf(id <= 0, ErrorCode.PARAMS_ERROR);

        Model model = modelService.getById(id);
        ThrowUtils.throwIf(model == null, ErrorCode.NOT_FOUND_ERROR);

        return ResultUtils.success(modelService.getModelVO(model));
    }

    /**
     * 分页获取模型列表（含提供者信息）
     *
     * @param modelQueryRequest 查询请求
     * @return 分页结果
     */
    @PostMapping("/list/page/vo")
    public BaseResponse<Page<ModelVO>> listModelVOByPage(@RequestBody ModelQueryRequest modelQueryRequest) {
        long pageNum = modelQueryRequest.getPageNum();
        long pageSize = modelQueryRequest.getPageSize();

        // 查询数据库
        QueryWrapper queryWrapper = modelService.getQueryWrapper(modelQueryRequest);
        Page<Model> modelPage = modelService.page(Page.of(pageNum, pageSize), queryWrapper);

        // 转换为 VO
        Page<ModelVO> modelVOPage = new Page<>(pageNum, pageSize, modelPage.getTotalRow());
        modelVOPage.setRecords(modelService.getModelVOList(modelPage.getRecords()));

        return ResultUtils.success(modelVOPage);
    }

    /**
     * 获取所有模型列表（含提供者信息）
     *
     * @return 模型列表
     */
    @GetMapping("/list/vo")
    public BaseResponse<List<ModelVO>> listModelVO() {
        List<Model> modelList = modelService.list();
        List<ModelVO> modelVOList = modelService.getModelVOList(modelList);
        return ResultUtils.success(modelVOList);
    }

    /**
     * 获取所有启用的模型列表
     *
     * @return 启用的模型列表
     */
    @GetMapping("/list/active")
    public BaseResponse<List<ModelVO>> listActiveModels() {
        List<Model> activeModels = modelService.getActiveModels();
        List<ModelVO> modelVOList = modelService.getModelVOList(activeModels);
        return ResultUtils.success(modelVOList);
    }

    /**
     * 根据提供者ID获取启用的模型列表
     *
     * @param providerId 提供者ID
     * @return 启用的模型列表
     */
    @GetMapping("/list/active/provider/{providerId}")
    public BaseResponse<List<ModelVO>> listActiveModelsByProvider(@PathVariable Long providerId) {
        ThrowUtils.throwIf(providerId == null || providerId <= 0, ErrorCode.PARAMS_ERROR);

        List<Model> models = modelService.getActiveModelsByProviderId(providerId);
        List<ModelVO> modelVOList = modelService.getModelVOList(models);
        return ResultUtils.success(modelVOList);
    }

    /**
     * 根据模型类型获取启用的模型列表
     *
     * @param modelType 模型类型
     * @return 启用的模型列表
     */
    @GetMapping("/list/active/type/{modelType}")
    public BaseResponse<List<ModelVO>> listActiveModelsByType(@PathVariable String modelType) {
        ThrowUtils.throwIf(modelType == null || modelType.isEmpty(), ErrorCode.PARAMS_ERROR);

        List<Model> models = modelService.getActiveModelsByType(modelType);
        List<ModelVO> modelVOList = modelService.getModelVOList(models);
        return ResultUtils.success(modelVOList);
    }
}
