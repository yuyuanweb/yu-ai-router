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

  type BaseResponseApiKeyVO = {
    code?: number
    data?: ApiKeyVO
    message?: string
  }

  type BaseResponseBoolean = {
    code?: number
    data?: boolean
    message?: string
  }

  type BaseResponseListModelVO = {
    code?: number
    data?: ModelVO[]
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

  type BaseResponsePageUserVO = {
    code?: number
    data?: PageUserVO
    message?: string
  }

  type BaseResponseProviderVO = {
    code?: number
    data?: ProviderVO
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

  type BaseResponseUserVO = {
    code?: number
    data?: UserVO
    message?: string
  }

  type BlacklistRequest = {
    ip?: string
    reason?: string
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
  }

  type StreamResponse = {
    id: string
    object: string
    created: number
    model: string
    choices: StreamChoice[]
  }

  type StreamChoice = {
    index: number
    delta: StreamDelta
    finishReason: string | null
  }

  type StreamDelta = {
    role?: string
    content?: string
    reasoningContent?: string
  }

  type checkBlacklistParams = {
    ip: string
  }

  type DeleteRequest = {
    id?: number
  }

  type getModelVOByIdParams = {
    id: number
  }

  type getMyLogsParams = {
    limit?: number
  }

  type getProviderVOByIdParams = {
    id: number
  }

  type getUserByIdParams = {
    id: number
  }

  type getUserVOByIdParams = {
    id: number
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

  type PageUserVO = {
    records?: UserVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
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
    createTime?: string
  }

  // 阶段五新增类型定义

  type QuotaVO = {
    tokenQuota?: number
    usedTokens?: number
    remainingQuota?: number
  }

  type QuotaUpdateRequest = {
    userId?: number
    tokenQuota?: number
  }

  type CostStatsVO = {
    totalCost?: number
    todayCost?: number
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

  type PageRequestLog = {
    records?: RequestLog[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type BaseResponseQuotaVO = {
    code?: number
    data?: QuotaVO
    message?: string
  }

  type BaseResponseCostStatsVO = {
    code?: number
    data?: CostStatsVO
    message?: string
  }

  type BaseResponseUserSummaryStatsVO = {
    code?: number
    data?: UserSummaryStatsVO
    message?: string
  }

  type BaseResponseUserAnalysisVO = {
    code?: number
    data?: UserAnalysisVO
    message?: string
  }

  type BaseResponsePageRequestLog = {
    code?: number
    data?: PageRequestLog
    message?: string
  }

  type BaseResponseRequestLog = {
    code?: number
    data?: RequestLog
    message?: string
  }

  type BaseResponseListMapStringObject = {
    code?: number
    data?: Record<string, any>[]
    message?: string
  }
}
