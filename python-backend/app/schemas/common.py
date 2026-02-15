"""Common schemas."""

from __future__ import annotations

from datetime import datetime
from typing import Generic, TypeVar

from pydantic import BaseModel, ConfigDict, Field, field_serializer

from app.core.constants import DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE

T = TypeVar("T")


class CamelBaseModel(BaseModel):
    model_config = ConfigDict(populate_by_name=True, from_attributes=True)


class BaseResponse(CamelBaseModel, Generic[T]):
    code: int
    data: T | None = None
    message: str


class DeleteRequest(CamelBaseModel):
    id: int


class PageRequest(CamelBaseModel):
    page_num: int = Field(default=DEFAULT_PAGE_NUM, alias="pageNum")
    page_size: int = Field(default=DEFAULT_PAGE_SIZE, alias="pageSize")
    sort_field: str | None = Field(default=None, alias="sortField")
    sort_order: str | None = Field(default="descend", alias="sortOrder")


class PageData(CamelBaseModel, Generic[T]):
    records: list[T] = Field(default_factory=list)
    page_number: int = Field(alias="pageNumber")
    page_size: int = Field(alias="pageSize")
    total_page: int = Field(alias="totalPage")
    total_row: int = Field(alias="totalRow")
    optimize_count_query: bool = Field(default=True, alias="optimizeCountQuery")


class LongIdModel(CamelBaseModel):
    id: int | None = None

    @field_serializer("id", when_used="json")
    def serialize_id(self, value: int | None) -> str | None:
        return str(value) if value is not None else None


class TimeModel(CamelBaseModel):
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")
