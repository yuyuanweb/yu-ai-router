"""Plugin manager."""

from __future__ import annotations

import logging
import time

from app.models.plugin_config import PluginConfig
from app.plugin.plugin import Plugin
from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_result import PluginResult

logger = logging.getLogger("app")


class PluginManager:
    def __init__(self) -> None:
        self._plugins: dict[str, Plugin] = {}
        self._plugin_configs: dict[str, PluginConfig] = {}

    def register_plugin(self, plugin: Plugin, config: PluginConfig | None) -> None:
        plugin_key = plugin.get_plugin_key()
        old_plugin = self._plugins.get(plugin_key)
        if old_plugin is not None:
            old_plugin.destroy()
            logger.info("卸载旧插件: %s", plugin_key)
        plugin.init(config.config if config is not None else None)
        self._plugins[plugin_key] = plugin
        if config is not None:
            self._plugin_configs[plugin_key] = config
        logger.info("注册插件: %s (%s)", plugin_key, plugin.get_plugin_name())

    def unregister_plugin(self, plugin_key: str) -> None:
        plugin = self._plugins.pop(plugin_key, None)
        self._plugin_configs.pop(plugin_key, None)
        if plugin is not None:
            plugin.destroy()
            logger.info("注销插件: %s", plugin_key)

    def get_plugin(self, plugin_key: str) -> Plugin | None:
        return self._plugins.get(plugin_key)

    def get_plugin_config(self, plugin_key: str) -> PluginConfig | None:
        return self._plugin_configs.get(plugin_key)

    def update_plugin_config(self, config: PluginConfig) -> None:
        plugin_key = config.plugin_key
        self._plugin_configs[plugin_key] = config
        plugin = self._plugins.get(plugin_key)
        if plugin is not None:
            plugin.init(config.config)
            logger.info("更新插件配置: %s", plugin_key)

    def is_plugin_enabled(self, plugin_key: str) -> bool:
        config = self._plugin_configs.get(plugin_key)
        if config is None:
            return False
        return config.status == "active"

    async def execute_plugin(self, plugin_key: str, context: PluginContext) -> PluginResult:
        plugin = self._plugins.get(plugin_key)
        if plugin is None:
            return PluginResult.fail(f"插件不存在: {plugin_key}")
        if not self.is_plugin_enabled(plugin_key):
            return PluginResult.fail(f"插件未启用: {plugin_key}")
        if not plugin.supports(context):
            return PluginResult.fail("插件不支持当前请求")
        start = time.perf_counter()
        try:
            result = await plugin.execute(context)
            result.duration = int((time.perf_counter() - start) * 1000)
            logger.info("插件执行完成: %s, 耗时: %sms", plugin_key, result.duration)
            return result
        except Exception as exc:
            logger.error("插件执行失败: %s", plugin_key, exc_info=True)
            result = PluginResult.fail(f"插件执行异常: {exc}")
            result.duration = int((time.perf_counter() - start) * 1000)
            return result

    def get_all_plugins(self) -> list[Plugin]:
        return list(self._plugins.values())

    def get_enabled_plugins(self) -> list[Plugin]:
        return [plugin for key, plugin in self._plugins.items() if self.is_plugin_enabled(key)]

    def get_plugin_count(self) -> int:
        return len(self._plugins)


plugin_manager = PluginManager()
