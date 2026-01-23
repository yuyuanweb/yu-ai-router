package com.yupi.airouter.service.impl;

import com.mybatisflex.core.query.QueryWrapper;
import com.yupi.airouter.mapper.ModelMapper;
import com.yupi.airouter.mapper.RequestLogMapper;
import com.yupi.airouter.model.entity.Model;
import com.yupi.airouter.model.entity.RequestLog;
import com.yupi.airouter.service.BillingService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

/**
 * 费用计算服务实现
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
@Slf4j
@Service
public class BillingServiceImpl implements BillingService {

    @Resource
    private ModelMapper modelMapper;

    @Resource
    private RequestLogMapper requestLogMapper;

    /**
     * 每千Token的价格基数
     */
    private static final BigDecimal TOKENS_PER_UNIT = new BigDecimal("1000");

    @Override
    public BigDecimal calculateCost(Model model, int promptTokens, int completionTokens) {
        if (model == null) {
            return BigDecimal.ZERO;
        }

        BigDecimal inputPrice = model.getInputPrice();
        BigDecimal outputPrice = model.getOutputPrice();

        if (inputPrice == null) {
            inputPrice = BigDecimal.ZERO;
        }
        if (outputPrice == null) {
            outputPrice = BigDecimal.ZERO;
        }

        // 费用 = (输入Token数 * 输入价格 + 输出Token数 * 输出价格) / 1000
        BigDecimal inputCost = inputPrice.multiply(new BigDecimal(promptTokens))
                .divide(TOKENS_PER_UNIT, 6, RoundingMode.HALF_UP);
        BigDecimal outputCost = outputPrice.multiply(new BigDecimal(completionTokens))
                .divide(TOKENS_PER_UNIT, 6, RoundingMode.HALF_UP);

        return inputCost.add(outputCost);
    }

    @Override
    public BigDecimal calculateCost(Long modelId, int promptTokens, int completionTokens) {
        if (modelId == null) {
            return BigDecimal.ZERO;
        }

        Model model = modelMapper.selectOneById(modelId);
        return calculateCost(model, promptTokens, completionTokens);
    }

    @Override
    public BigDecimal getUserTotalCost(Long userId) {
        if (userId == null) {
            return BigDecimal.ZERO;
        }

        List<RequestLog> logs = requestLogMapper.selectListByQuery(
                QueryWrapper.create()
                        .where("userId = " + userId)
                        .and("status = 'success'")
        );

        BigDecimal totalCost = BigDecimal.ZERO;
        for (RequestLog log : logs) {
            if (log.getCost() != null) {
                totalCost = totalCost.add(log.getCost());
            }
        }

        return totalCost;
    }

    @Override
    public BigDecimal getUserTodayCost(Long userId) {
        if (userId == null) {
            return BigDecimal.ZERO;
        }

        LocalDateTime todayStart = LocalDateTime.of(LocalDate.now(), LocalTime.MIN);
        LocalDateTime todayEnd = LocalDateTime.of(LocalDate.now(), LocalTime.MAX);

        List<RequestLog> logs = requestLogMapper.selectListByQuery(
                QueryWrapper.create()
                        .where("userId = " + userId)
                        .and("status = 'success'")
                        .and("createTime >= '" + todayStart + "'")
                        .and("createTime <= '" + todayEnd + "'")
        );

        BigDecimal totalCost = BigDecimal.ZERO;
        for (RequestLog log : logs) {
            if (log.getCost() != null) {
                totalCost = totalCost.add(log.getCost());
            }
        }

        return totalCost;
    }
}
