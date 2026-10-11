export interface JsonApiResource<T> {
  id: string;
  type: string;
  attributes: T;
}

export interface SingleResponse<T> {
  data: JsonApiResource<T>;
}

export interface NullableSingleResponse<T> {
  data: JsonApiResource<T> | null;
}

export interface CollectionResponse<T> {
  data: JsonApiResource<T>[];
}

export interface PaginationMeta {
  total: number;
  limit: number;
  offset: number;
}

export interface PaginationLinks {
  self: string;
  first: string;
  prev: string | null;
  next: string | null;
  last: string;
}

export interface PaginatedCollectionResponse<T> extends CollectionResponse<T> {
  meta: PaginationMeta;
  links: PaginationLinks;
}

export interface ErrorObject {
  status?: string;
  title?: string;
  detail?: string;
  code?: string;
}

export interface ErrorResponse {
  errors: ErrorObject[];
}
