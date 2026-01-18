export interface JsonApiResource<T> {
  id: string;
  type: string;
  attributes: T;
}

export interface SingleResponse<T> {
  data: JsonApiResource<T>;
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
