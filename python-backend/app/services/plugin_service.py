"""Plugin service."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import ErrorCode
from app.exceptions.business_exception import BusinessException
from app.models.plugin_config import PluginConfig
from app.plugin.impl.image_recognition_plugin import ImageRecognitionPlugin
from app.plugin.impl.pdf_parser_plugin import PdfParserPlugin
from app.plugin.impl.web_search_plugin import WebSearchPlugin
from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_manager import plugin_manager
from app.schemas.plugin import PluginConfigVO, PluginExecuteRequest, PluginExecuteVO, PluginUpdateRequest


class PluginService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def list_all_plugins(self) -> list[PluginConfigVO]:
        stmt = (
            select(PluginConfig)
            .where(PluginConfig.is_delete == 0)
            .order_by(PluginConfig.priority.desc())
        )
        rows = list((await self.db.scalars(stmt)).all())
        return [PluginConfigVO.model_validate(item) for item in rows]

    async def list_enabled_plugins(self) -> list[PluginConfigVO]:
        stmt = (
            select(PluginConfig)
            .where(PluginConfig.is_delete == 0, PluginConfig.status == "active")
            .order_by(PluginConfig.priority.desc())
        )
        rows = list((await self.db.scalars(stmt)).all())
        return [PluginConfigVO.model_validate(item) for item in rows]

    async def get_plugin_by_key(self, plugin_key: str) -> PluginConfigVO | None:
        stmt = select(PluginConfig).where(
            PluginConfig.plugin_key == plugin_key,
            PluginConfig.is_delete == 0,
        )
        entity = await self.db.scalar(stmt)
        if entity is None:
            return None
        return PluginConfigVO.model_validate(entity)

    async def update_plugin(self, request: PluginUpdateRequest) -> bool:
        if request.id is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "插件 ID 不能为空")
        entity = await self.db.scalar(
            select(PluginConfig).where(
                PluginConfig.id == request.id,
                PluginConfig.is_delete == 0,
            )
        )
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在")
        if request.plugin_name is not None:
            entity.plugin_name = request.plugin_name
        if request.description is not None:
            entity.description = request.description
        if request.config is not None:
            entity.config = request.config
        if request.status is not None:
            entity.status = request.status
        if request.priority is not None:
            entity.priority = request.priority
        await self.db.commit()
        plugin_manager.update_plugin_config(entity)
        return True

    async def enable_plugin(self, plugin_key: str) -> bool:
        entity = await self._get_entity_by_key(plugin_key)
        entity.status = "active"
        await self.db.commit()
        plugin_manager.update_plugin_config(entity)
        return True

    async def disable_plugin(self, plugin_key: str) -> bool:
        entity = await self._get_entity_by_key(plugin_key)
        entity.status = "inactive"
        await self.db.commit()
        plugin_manager.update_plugin_config(entity)
        return True

    async def execute_plugin(self, request: PluginExecuteRequest, user_id: int | None) -> PluginExecuteVO:
        if not request.plugin_key:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空")
        context = PluginContext(
            user_id=user_id,
            input=request.input,
            file_url=request.file_url,
            file_bytes=request.file_bytes,
            file_type=request.file_type,
            params=request.params or {},
        )
        result = await plugin_manager.execute_plugin(request.plugin_key, context)
        return PluginExecuteVO(
            success=result.success,
            pluginKey=request.plugin_key,
            content=result.content,
            errorMessage=result.error_message,
            duration=result.duration,
            data=result.data or None,
        )

    async def init_plugins(self) -> None:
        plugin_configs = list(
            (
                await self.db.scalars(
                    select(PluginConfig).where(PluginConfig.is_delete == 0)
                )
            ).all()
        )
        builtins = [
            WebSearchPlugin(),
            PdfParserPlugin(),
            ImageRecognitionPlugin(),
        ]
        for plugin in builtins:
            plugin_key = plugin.get_plugin_key()
            config = next((item for item in plugin_configs if item.plugin_key == plugin_key), None)
            if config is not None:
                plugin_manager.register_plugin(plugin, config)

    async def reload_plugin(self, plugin_key: str) -> None:
        config = await self._get_entity_by_key(plugin_key)
        plugin_map = {
            "web_search": WebSearchPlugin(),
            "pdf_parser": PdfParserPlugin(),
            "image_recognition": ImageRecognitionPlugin(),
        }
        plugin = plugin_map.get(plugin_key)
        if plugin is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件实例不存在")
        plugin_manager.register_plugin(plugin, config)

    async def _get_entity_by_key(self, plugin_key: str) -> PluginConfig:
        entity = await self.db.scalar(
            select(PluginConfig).where(
                PluginConfig.plugin_key == plugin_key,
                PluginConfig.is_delete == 0,
            )
        )
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在")
        return entity
