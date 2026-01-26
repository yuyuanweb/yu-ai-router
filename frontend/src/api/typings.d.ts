declare namespace API {
  type ApiKeyCreateRequest = {
    keyName?: string
  }

  type ApiKeyVO = {
    id?: number
    keyValue?: string
    keyName?: string
    status?: string
    totalTokens?: number
    lastUsedTime?: string
    createTime?: string
  }

  type BalanceVO = {
    balance?: number
    totalSpending?: number
    totalRecharge?: number
  }

  type BaseResponseApiKeyVO = {
    code?: number
    data?: ApiKeyVO
    message?: string
  }

  type BaseResponseBalanceVO = {
    code?: number
    data?: BalanceVO
    message?: string
  }

  type BaseResponseBoolean = {
    code?: number
    data?: boolean
    message?: string
  }

  type BaseResponseCostStatsVO = {
    code?: number
    data?: CostStatsVO
    message?: string
  }

  type BaseResponseCreateRechargeResponse = {
    code?: number
    data?: CreateRechargeResponse
    message?: string
  }

  type BaseResponseListMapStringObject = {
    code?: number
    data?: Record<string, any>[]
    message?: string
  }

  type BaseResponseListModelVO = {
    code?: number
    data?: ModelVO[]
    message?: string
  }

  type BaseResponseListPluginConfigVO = {
    code?: number
    data?: PluginConfigVO[]
    message?: string
  }

  type BaseResponseListProviderVO = {
    code?: number
    data?: ProviderVO[]
    message?: string
  }

  type BaseResponseListRequestLog = {
    code?: number
    data?: RequestLog[]
    message?: string
  }

  type BaseResponseLoginUserVO = {
    code?: number
    data?: LoginUserVO
    message?: string
  }

  type BaseResponseLong = {
    code?: number
    data?: number
    message?: string
  }

  type BaseResponseModelVO = {
    code?: number
    data?: ModelVO
    message?: string
  }

  type BaseResponsePageApiKeyVO = {
    code?: number
    data?: PageApiKeyVO
    message?: string
  }

  type BaseResponsePageBillingRecord = {
    code?: number
    data?: PageBillingRecord
    message?: string
  }

  type BaseResponsePageImageGenerationRecord = {
    code?: number
    data?: PageImageGenerationRecord
    message?: string
  }

  type BaseResponsePageModelVO = {
    code?: number
    data?: PageModelVO
    message?: string
  }

  type BaseResponsePageProviderVO = {
    code?: number
    data?: PageProviderVO
    message?: string
  }

  type BaseResponsePageRechargeRecord = {
    code?: number
    data?: PageRechargeRecord
    message?: string
  }

  type BaseResponsePageRequestLog = {
    code?: number
    data?: PageRequestLog
    message?: string
  }

  type BaseResponsePageUserVO = {
    code?: number
    data?: PageUserVO
    message?: string
  }

  type BaseResponsePluginConfigVO = {
    code?: number
    data?: PluginConfigVO
    message?: string
  }

  type BaseResponsePluginExecuteVO = {
    code?: number
    data?: PluginExecuteVO
    message?: string
  }

  type BaseResponseProviderVO = {
    code?: number
    data?: ProviderVO
    message?: string
  }

  type BaseResponseQuotaVO = {
    code?: number
    data?: QuotaVO
    message?: string
  }

  type BaseResponseRequestLog = {
    code?: number
    data?: RequestLog
    message?: string
  }

  type BaseResponseSetString = {
    code?: number
    data?: string[]
    message?: string
  }

  type BaseResponseString = {
    code?: number
    data?: string
    message?: string
  }

  type BaseResponseTokenStatsVO = {
    code?: number
    data?: TokenStatsVO
    message?: string
  }

  type BaseResponseUser = {
    code?: number
    data?: User
    message?: string
  }

  type BaseResponseUserAnalysisVO = {
    code?: number
    data?: UserAnalysisVO
    message?: string
  }

  type BaseResponseUserSummaryStatsVO = {
    code?: number
    data?: UserSummaryStatsVO
    message?: string
  }

  type BaseResponseUserVO = {
    code?: number
    data?: UserVO
    message?: string
  }

  type BillingRecord = {
    id?: number
    userId?: number
    requestLogId?: number
    amount?: number
    balanceBefore?: number
    balanceAfter?: number
    description?: string
    billingType?: string
    createTime?: string
  }

  type BlacklistRequest = {
    ip?: string
    reason?: string
  }

  type chatCompletionsWithFileParams = {
    model?: string
    messages: string
    stream?: boolean
    routing_strategy?: string
    plugin_key?: string
    enable_reasoning?: boolean
  }

  type ChatMessage = {
    role?: string
    content?: string
  }

  type ChatRequest = {
    model?: string
    messages?: ChatMessage[]
    stream?: boolean
    temperature?: number
    max_tokens?: number
    enable_reasoning?: boolean
    routing_strategy?: string
    plugin_key?: string
    file_url?: string
    file_type?: string
  }

  type checkBlacklistParams = {
    ip: string
  }

  type CostStatsVO = {
    totalCost?: number
    todayCost?: number
  }

  type CreateRechargeRequest = {
    amount?: number
  }

  type CreateRechargeResponse = {
    checkoutUrl?: string
    sessionId?: string
  }

  type DeleteRequest = {
    id?: number
  }

  type disablePluginParams = {
    pluginKey: string
  }

  type disableUserParams = {
    userId: number
  }

  type enablePluginParams = {
    pluginKey: string
  }

  type enableUserParams = {
    userId: number
  }

  type getHistoryDetailParams = {
    id: number
  }

  type getModelVOByIdParams = {
    id: number
  }

  type getMyBillingRecordsParams = {
    pageNum?: number
    pageSize?: number
  }

  type getMyDailyStatsParams = {
    startDate?: string
    endDate?: string
  }

  type getMyLogsParams = {
    limit?: number
  }

  type getMyRechargeRecordsParams = {
    pageNum?: number
    pageSize?: number
  }

  type getMyRecordsParams = {
    pageNum?: number
    pageSize?: number
  }

  type getPluginParams = {
    pluginKey: string
  }

  type getProviderVOByIdParams = {
    id: number
  }

  type getUserAnalysisParams = {
    userId: number
  }

  type getUserByIdParams = {
    id: number
  }

  type getUserVOByIdParams = {
    id: number
  }

  type ImageData = {
    url?: string
    b64Json?: string
    revisedPrompt?: string
  }

  type ImageGenerationRecord = {
    id?: number
    userId?: number
    apiKeyId?: number
    modelId?: number
    modelKey?: string
    prompt?: string
    revisedPrompt?: string
    imageUrl?: string
    imageData?: string
    size?: string
    quality?: string
    status?: string
    cost?: number
    duration?: number
    errorMessage?: string
    clientIp?: string
    createTime?: string
  }

  type ImageGenerationRequest = {
    prompt?: string
    model?: string
    size?: string
    quality?: string
    user?: string
    n?: number
    response_format?: string
  }

  type ImageGenerationResponse = {
    created?: number
    data?: ImageData[]
  }

  type listActiveModelsByProviderParams = {
    providerId: number
  }

  type listActiveModelsByTypeParams = {
    modelType: string
  }

  type listMyApiKeysParams = {
    pageNum?: number
    pageSize?: number
  }

  type LoginUserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    userStatus?: string
    tokenQuota?: number
    usedTokens?: number
    balance?: number
    createTime?: string
    updateTime?: string
  }

  type ModelAddRequest = {
    providerId?: number
    modelKey?: string
    modelName?: string
    modelType?: string
    description?: string
    contextLength?: number
    inputPrice?: number
    outputPrice?: number
    priority?: number
    defaultTimeout?: number
    capabilities?: string
  }

  type ModelQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    providerId?: number
    modelKey?: string
    modelName?: string
    modelType?: string
    status?: string
  }

  type ModelUpdateRequest = {
    id?: number
    modelName?: string
    description?: string
    contextLength?: number
    inputPrice?: number
    outputPrice?: number
    status?: string
    priority?: number
    defaultTimeout?: number
    capabilities?: string
  }

  type ModelVO = {
    id?: number
    providerId?: number
    providerName?: string
    providerDisplayName?: string
    modelKey?: string
    modelName?: string
    modelType?: string
    description?: string
    contextLength?: number
    inputPrice?: number
    outputPrice?: number
    status?: string
    healthStatus?: string
    avgLatency?: number
    successRate?: number
    priority?: number
    defaultTimeout?: number
    supportReasoning?: number
    capabilities?: string
    createTime?: string
    updateTime?: string
  }

  type PageApiKeyVO = {
    records?: ApiKeyVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageBillingRecord = {
    records?: BillingRecord[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageImageGenerationRecord = {
    records?: ImageGenerationRecord[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageModelVO = {
    records?: ModelVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageProviderVO = {
    records?: ProviderVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageRechargeRecord = {
    records?: RechargeRecord[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageRequestLog = {
    records?: RequestLog[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageUserVO = {
    records?: UserVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PluginConfigVO = {
    id?: number
    pluginKey?: string
    pluginName?: string
    pluginType?: string
    description?: string
    config?: string
    status?: string
    priority?: number
    createTime?: string
    updateTime?: string
  }

  type PluginExecuteRequest = {
    pluginKey?: string
    input?: string
    fileUrl?: string
    fileType?: string
    params?: Record<string, any>
    fileBytes?: string[]
  }

  type PluginExecuteVO = {
    success?: boolean
    pluginKey?: string
    content?: string
    errorMessage?: string
    duration?: number
    data?: Record<string, any>
  }

  type PluginUpdateRequest = {
    id?: number
    pluginName?: string
    description?: string
    config?: string
    status?: string
    priority?: number
  }

  type ProviderAddRequest = {
    providerName?: string
    displayName?: string
    baseUrl?: string
    apiKey?: string
    priority?: number
    config?: string
  }

  type ProviderQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    providerName?: string
    displayName?: string
    status?: string
    healthStatus?: string
  }

  type ProviderUpdateRequest = {
    id?: number
    displayName?: string
    baseUrl?: string
    apiKey?: string
    status?: string
    priority?: number
    config?: string
  }

  type ProviderVO = {
    id?: number
    providerName?: string
    displayName?: string
    baseUrl?: string
    status?: string
    healthStatus?: string
    avgLatency?: number
    successRate?: number
    priority?: number
    config?: string
    createTime?: string
    updateTime?: string
  }

  type QuotaUpdateRequest = {
    userId?: number
    tokenQuota?: number
  }

  type QuotaVO = {
    tokenQuota?: number
    usedTokens?: number
    remainingQuota?: number
  }

  type RechargeRecord = {
    id?: number
    userId?: number
    amount?: number
    paymentMethod?: string
    paymentId?: string
    status?: string
    description?: string
    createTime?: string
    updateTime?: string
  }

  type reloadPluginParams = {
    pluginKey: string
  }

  type RequestLog = {
    id?: number
    traceId?: string
    userId?: number
    apiKeyId?: number
    modelId?: number
    requestModel?: string
    modelName?: string
    requestType?: string
    source?: string
    promptTokens?: number
    completionTokens?: number
    totalTokens?: number
    cost?: number
    duration?: number
    status?: string
    errorMessage?: string
    errorCode?: string
    routingStrategy?: string
    isFallback?: number
    clientIp?: string
    userAgent?: string
    createTime?: string
    updateTime?: string
  }

  type RequestLogQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    userId?: number
    requestModel?: string
    requestType?: string
    source?: string
    status?: string
    startDate?: string
    endDate?: string
  }

  type resetUserQuotaParams = {
    userId: number
  }

  type stripeSuccessParams = {
    session_id: string
  }

  type TokenStatsVO = {
    totalTokens?: number
  }

  type User = {
    id?: number
    userAccount?: string
    userPassword?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    userStatus?: string
    tokenQuota?: number
    usedTokens?: number
    balance?: number
    editTime?: string
    createTime?: string
    updateTime?: string
    isDelete?: number
  }

  type UserAddRequest = {
    userName?: string
    userAccount?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
  }

  type UserAnalysisVO = {
    userId?: number
    userAccount?: string
    userName?: string
    userStatus?: string
    userRole?: string
    tokenQuota?: number
    usedTokens?: number
    remainingQuota?: number
    totalRequests?: number
    successRequests?: number
    totalTokens?: number
    totalCost?: number
    todayCost?: number
  }

  type UserLoginRequest = {
    userAccount?: string
    userPassword?: string
  }

  type UserQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    id?: number
    userName?: string
    userAccount?: string
    userProfile?: string
    userRole?: string
  }

  type UserRegisterRequest = {
    userAccount?: string
    userPassword?: string
    checkPassword?: string
  }

  type UserSummaryStatsVO = {
    totalTokens?: number
    tokenQuota?: number
    usedTokens?: number
    remainingQuota?: number
    totalCost?: number
    todayCost?: number
    totalRequests?: number
    successRequests?: number
  }

  type UserUpdateRequest = {
    id?: number
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    tokenQuota?: number
  }

  type UserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    userStatus?: string
    tokenQuota?: number
    usedTokens?: number
    balance?: number
    createTime?: string
  }
}
